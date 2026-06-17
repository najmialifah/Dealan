param (
    [Parameter(Mandatory=$true)]
    [string]$ServiceName,
    
    [Parameter(Mandatory=$false)]
    [string]$ServicePort
)

# Set default ports based on the service if no port is provided
if (-not $ServicePort) {
    switch ($ServiceName) {
        "auth-service" { $ServicePort = "3001" }
        "user-service" { $ServicePort = "3002" }
        "driver-service" { $ServicePort = "3003" }
        "order-service" { $ServicePort = "3004" }
        "matching-service" { $ServicePort = "3005" }
        "pricing-service" { $ServicePort = "3006" }
        "chat-service" { $ServicePort = "8087" }
        "map-route-service" { $ServicePort = "8088" }
        "notification-service" { $ServicePort = "8084" }
        "payment-service" { $ServicePort = "8093" }
        "promo-service" { $ServicePort = "3012" }
        "punishment-service" { $ServicePort = "8086" }
        "rating-review-service" { $ServicePort = "8085" }
        "shipment-service" { $ServicePort = "8094" }
        default { $ServicePort = "8080" }
    }
}

Write-Host "=========================================" -ForegroundColor Cyan
Write-Host "Memulai $ServiceName di port $ServicePort..." -ForegroundColor Green
Write-Host "Pastikan infrastruktur (PostgreSQL, Kafka) sudah jalan!" -ForegroundColor Yellow
Write-Host "=========================================" -ForegroundColor Cyan

# Set Environment Variables
$env:DB_HOST="127.0.0.1"
$env:DB_USER="postgres"
$env:DB_PASSWORD="password"
$env:DB_NAME="dealan"
$env:DB_PORT="5433"
$env:PORT=$ServicePort
$env:KAFKA_BROKER="127.0.0.1:9092"
$env:KAFKA_BROKERS="127.0.0.1:9092"
$env:DB_URL="host=127.0.0.1 user=postgres password=password dbname=dealan port=5433 sslmode=disable"

# Pindah ke direktori service
if (Test-Path $ServiceName) {
    Set-Location $ServiceName
    
    # Jalankan service
    Write-Host "Menjalankan: go run cmd/main.go"
    go run cmd/main.go
    
    # Kembali ke root setelah selesai
    Set-Location ..
} else {
    Write-Host "Error: Folder service '$ServiceName' tidak ditemukan!" -ForegroundColor Red
}
