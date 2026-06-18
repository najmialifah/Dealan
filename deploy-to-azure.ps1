param (
    [Parameter(Mandatory=$false)]
    [string]$ResourceGroupName = "dealan-rg",

    [Parameter(Mandatory=$false)]
    [string]$AcrName = "dealanregistry",

    [Parameter(Mandatory=$false)]
    [string]$AksClusterName = "dealan-aks-cluster",

    [Parameter(Mandatory=$false)]
    [string]$Location = "eastasia",

    [Parameter(Mandatory=$false)]
    [string]$NodeVmSize = "standard_b2s_v2"
)

# Hentikan eksekusi jika terjadi error
$ErrorActionPreference = "Stop"

$services = @(
    "auth-service",
    "chat-service",
    "driver-service",
    "location-service",
    "map-route-service",
    "matching-service",
    "notification-service",
    "order-service",
    "payment-service",
    "pricing-service",
    "promo-service",
    "punishment-service",
    "rating-review-service",
    "shipment-service",
    "user-service"
)

Write-Host "==================================================" -ForegroundColor Cyan
Write-Host "DEALAN AUTOMATED DEPLOYMENT TO AZURE" -ForegroundColor Cyan
Write-Host "==================================================" -ForegroundColor Cyan

# 1. Pengecekan Login Azure CLI
Write-Host "1. Memeriksa login Azure CLI..." -ForegroundColor Green
$azAccount = az account show --query "name" -o tsv 2>$null
if (-not $azAccount) {
    Write-Error "Anda belum login ke Azure CLI. Harap jalankan 'az login' terlebih dahulu di terminal Anda."
    exit 1
}
Write-Host "Terhubung sebagai akun: $azAccount" -ForegroundColor Yellow

# 1.5. Registrasi Resource Provider (Wajib untuk subskripsi baru / Azure for Students)
Write-Host "`n1.5. Mendaftarkan Resource Provider ke subskripsi..." -ForegroundColor Green
Write-Host "Mendaftarkan Microsoft.ContainerRegistry..." -ForegroundColor Blue
az provider register --namespace Microsoft.ContainerRegistry --wait
if ($LASTEXITCODE -ne 0) {
    Write-Error "Gagal meregistrasi Microsoft.ContainerRegistry provider."
    exit 1
}
Write-Host "Mendaftarkan Microsoft.ContainerService..." -ForegroundColor Blue
az provider register --namespace Microsoft.ContainerService --wait
if ($LASTEXITCODE -ne 0) {
    Write-Error "Gagal meregistrasi Microsoft.ContainerService provider."
    exit 1
}
Write-Host "Pendaftaran provider selesai." -ForegroundColor Yellow

# 2. Pembuatan / Pemeriksaan Resource Group
Write-Host "`n2. Memeriksa Resource Group..." -ForegroundColor Green
$rgExists = az group exists --name $ResourceGroupName
if ($rgExists -eq "false") {
    Write-Host "Membuat Resource Group '$ResourceGroupName' di lokasi '$Location'..." -ForegroundColor Blue
    az group create --name $ResourceGroupName --location $Location > $null
    if ($LASTEXITCODE -ne 0) {
        Write-Error "Gagal membuat Resource Group."
        exit 1
    }
    Write-Host "Berhasil membuat Resource Group." -ForegroundColor Yellow
} else {
    Write-Host "Resource Group '$ResourceGroupName' sudah tersedia." -ForegroundColor Yellow
}

# 3. Pembuatan / Pemeriksaan ACR
Write-Host "`n3. Memeriksa Azure Container Registry..." -ForegroundColor Green
$acrLoginServer = az acr list --query "[?name=='$AcrName'].loginServer" -o tsv
if ($acrLoginServer) { $acrLoginServer = $acrLoginServer.Trim() }

if (-not $acrLoginServer) {
    Write-Host "Membuat ACR '$AcrName' (ini memakan waktu sekitar 1-2 menit)..." -ForegroundColor Blue
    az acr create --resource-group $ResourceGroupName --name $AcrName --sku Basic > $null
    if ($LASTEXITCODE -ne 0) {
        Write-Error "Gagal membuat Azure Container Registry."
        exit 1
    }
    $acrLoginServer = az acr list --query "[?name=='$AcrName'].loginServer" -o tsv
    if ($acrLoginServer) { $acrLoginServer = $acrLoginServer.Trim() }
    Write-Host "Berhasil membuat ACR: $acrLoginServer" -ForegroundColor Yellow
} else {
    Write-Host "ACR '$AcrName' sudah tersedia: $acrLoginServer" -ForegroundColor Yellow
}

# 4. Pembuatan / Pemeriksaan AKS
Write-Host "`n4. Memeriksa AKS Cluster..." -ForegroundColor Green
$aksExists = az aks list --query "[?name=='$AksClusterName'].id" -o tsv

if (-not $aksExists) {
    Write-Host "Membuat AKS Cluster '$AksClusterName' (VM Size: $NodeVmSize) & menghubungkannya ke ACR '$AcrName'..." -ForegroundColor Blue
    Write-Host "(Ini memakan waktu sekitar 5-10 menit. Harap bersabar...)" -ForegroundColor Blue
    az aks create `
        --resource-group $ResourceGroupName `
        --name $AksClusterName `
        --node-count 2 `
        --node-vm-size $NodeVmSize `
        --generate-ssh-keys `
        --attach-acr $AcrName > $null
    if ($LASTEXITCODE -ne 0) {
        Write-Error "Gagal membuat AKS Cluster."
        exit 1
    }
    Write-Host "Berhasil membuat cluster AKS." -ForegroundColor Yellow
} else {
    Write-Host "AKS Cluster '$AksClusterName' sudah tersedia." -ForegroundColor Yellow
}

# 5. Login ke Registry ACR
Write-Host "`n5. Melakukan login ke registry ACR..." -ForegroundColor Green
az acr login --name $AcrName

# 6. Tagging & Pushing Local Images ke ACR
Write-Host "`n6. Mengunggah Docker Image ke ACR..." -ForegroundColor Green
foreach ($service in $services) {
    $localImage = "dealan-$($service):latest"
    $remoteImage = "$acrLoginServer/dealan-$($service):latest"
    
    # Memeriksa apakah image lokal tersedia
    $imageCheck = docker images -q $localImage
    if (-not $imageCheck) {
        Write-Warning "Docker image lokal '$localImage' tidak ditemukan. Harap pastikan image sudah di-build."
        continue
    }

    Write-Host "Tagging $localImage -> $remoteImage..." -ForegroundColor Blue
    docker tag $localImage $remoteImage

    Write-Host "Pushing $remoteImage ke ACR..." -ForegroundColor Blue
    docker push $remoteImage
}

# 7. Menghubungkan kubectl ke AKS
Write-Host "`n7. Menghubungkan kubectl ke cluster AKS..." -ForegroundColor Green
az aks get-credentials --resource-group $ResourceGroupName --name $AksClusterName --overwrite-existing

# 8. Menerapkan Manifest Kubernetes dengan Referensi Image ACR
Write-Host "`n8. Memproses dan menerapkan manifest K8s ke AKS..." -ForegroundColor Green

# Membuat folder temporer untuk modifikasi manifest agar tidak merusak manifest asli
$tmpDir = Join-Path $PSScriptRoot ".azure_deploy_tmp"
if (Test-Path $tmpDir) {
    Remove-Item -Path $tmpDir -Recurse -Force
}
New-Item -ItemType Directory -Path $tmpDir | Out-Null

Write-Host "Menerapkan ConfigMap dan Infrastruktur K8s (PostgreSQL, Redis, Kafka)..." -ForegroundColor Blue
kubectl apply -f (Join-Path $PSScriptRoot "dealan-config-k8s.yaml")
kubectl apply -f (Join-Path $PSScriptRoot "infrastructure-k8s.yaml")

Write-Host "Tunggu sekitar 60 detik agar infrastruktur stabil..." -ForegroundColor Yellow
Start-Sleep -Seconds 60


foreach ($service in $services) {
    # Memetakan nama file manifest yang sesuai
    $k8sFileName = "$service-k8s.yaml"
    if ($service -eq "auth-service") { $k8sFileName = "auth-k8s.yaml" }
    elseif ($service -eq "driver-service") { $k8sFileName = "driver-k8s.yaml" }
    elseif ($service -eq "notification-service") { $k8sFileName = "notification-k8s.yaml" }
    elseif ($service -eq "payment-service") { $k8sFileName = "payment-k8s.yaml" }
    elseif ($service -eq "pricing-service") { $k8sFileName = "pricing-k8s.yaml" }
    elseif ($service -eq "punishment-service") { $k8sFileName = "punishment-k8s.yaml" }
    elseif ($service -eq "rating-review-service") { $k8sFileName = "rating-review-k8s.yaml" }
    elseif ($service -eq "shipment-service") { $k8sFileName = "shipment-k8s.yaml" }
    elseif ($service -eq "user-service") { $k8sFileName = "user-k8s.yaml" }

    $k8sFile = Join-Path $PSScriptRoot "$service/$k8sFileName"

    if (Test-Path $k8sFile) {
        Write-Host "Mengonfigurasi manifest $service -> ACR..." -ForegroundColor Blue
        $content = Get-Content $k8sFile -Raw
        
        # Ganti image format kelompokdealan/<service-name>:latest dengan server ACR kita
        $targetImage = "$acrLoginServer/dealan-$($service):latest"
        $newContent = $content -replace "kelompokdealan/[a-zA-Z0-9_-]+:latest", $targetImage
        
        $outFile = Join-Path $tmpDir "$service-k8s.yaml"
        Set-Content -Path $outFile -Value $newContent
        
        Write-Host "Menerapkan manifest $service-k8s.yaml ke AKS..." -ForegroundColor Blue
        kubectl apply -f $outFile
    } else {
        Write-Warning "Manifest untuk $service tidak ditemukan pada path: $k8sFile"
    }
}

# 9. Setup Nginx Ingress Controller
Write-Host "`n9. Mengonfigurasi Ingress Controller..." -ForegroundColor Green
$helmExists = Get-Command helm -ErrorAction SilentlyContinue
if ($helmExists) {
    Write-Host "Menambahkan helm repository ingress-nginx..." -ForegroundColor Blue
    helm repo add ingress-nginx https://kubernetes.github.io/ingress-nginx 2>$null | Out-Null
    helm repo update 2>$null | Out-Null
    
    Write-Host "Menginstal Nginx Ingress Controller di namespace ingress-basic..." -ForegroundColor Blue
    helm install ingress-nginx ingress-nginx/ingress-nginx `
      --create-namespace `
      --namespace ingress-basic `
      --set controller.service.externalTrafficPolicy="Local" 2>$null | Out-Null
      
    Write-Host "Ingress Controller siap/diperbarui." -ForegroundColor Yellow
} else {
    Write-Warning "Helm tidak ditemukan di host. Lewati instalasi otomatis Ingress Controller. Silakan instal Ingress manual pada cluster AKS Anda."
}

# Terapkan ingress.yaml
$ingressFile = Join-Path $PSScriptRoot "ingress.yaml"
if (Test-Path $ingressFile) {
    Write-Host "Menerapkan ingress.yaml ke AKS..." -ForegroundColor Blue
    kubectl apply -f $ingressFile
} else {
    Write-Warning "File ingress.yaml tidak ditemukan di root directory!"
}

# Membersihkan file temporer
Remove-Item -Path $tmpDir -Recurse -Force

Write-Host "`n==================================================" -ForegroundColor Cyan
Write-Host "PROSES DEPLOYMENT KE AZURE BERHASIL DIRENCANAKAN!" -ForegroundColor Green
Write-Host "Gunakan 'kubectl get pods' untuk melihat pod Anda." -ForegroundColor Green
Write-Host "==================================================" -ForegroundColor Cyan
