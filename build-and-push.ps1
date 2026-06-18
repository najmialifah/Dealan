param (
    [string]$AcrName = "dealanregistry"
)

$ErrorActionPreference = "Stop"

$services = @(
    "auth-service", "chat-service", "driver-service", "location-service",
    "map-route-service", "matching-service", "notification-service",
    "order-service", "payment-service", "pricing-service", "promo-service",
    "punishment-service", "rating-review-service", "shipment-service", "user-service"
)

$acrLoginServer = "dealanregistry.azurecr.io"
Write-Host "ACR Login Server: $acrLoginServer" -ForegroundColor Cyan
Write-Host "`nMemulai proses Build dan Push 15 Services..." -ForegroundColor Green

foreach ($service in $services) {
    $remoteImage = "$acrLoginServer/dealan-$($service):latest"
    $servicePath = Join-Path $PSScriptRoot $service
    
    if (Test-Path "$servicePath/Dockerfile") {
        Write-Host "========================================" -ForegroundColor Cyan
        Write-Host "Building image: $remoteImage" -ForegroundColor Yellow
        Write-Host "Path: $servicePath" -ForegroundColor Yellow
        Write-Host "========================================" -ForegroundColor Cyan
        
        # Build image using local Docker daemon
        docker build -t $remoteImage $servicePath
        if ($LASTEXITCODE -ne 0) {
            Write-Error "Build gagal untuk $service"
            exit 1
        }
        
        Write-Host "Pushing image: $remoteImage..." -ForegroundColor Magenta
        docker push $remoteImage
        if ($LASTEXITCODE -ne 0) {
            Write-Error "Push gagal untuk $service"
            exit 1
        }
    } else {
        Write-Warning "Dockerfile tidak ditemukan di $servicePath"
    }
}

Write-Host "`nSemua image berhasil di-build dan di-push ke ACR!" -ForegroundColor Green
