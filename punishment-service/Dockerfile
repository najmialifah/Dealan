FROM golang:alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o main ./cmd/main.go

FROM alpine:3.18
WORKDIR /app
COPY --from=builder /app/main .
# Port untuk Punishment Service
EXPOSE 8086
CMD ["./main"]