# Pakai environment bahasa Go
FROM golang:1.21-alpine

# Bikin folder di dalam kotaknya
WORKDIR /app

# Copy semua kodingan kamu ke dalam kotak
COPY . .

# Download library yang dibutuhin
RUN go mod tidy

# Bangun aplikasinya
RUN go build -o main cmd/main.go

# Nyalain aplikasinya
CMD ["./main"]