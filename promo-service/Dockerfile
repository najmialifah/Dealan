FROM golang:latest

WORKDIR /app

COPY go.mod .
COPY go.sum .

RUN go mod tidy

COPY . .

RUN go build -o main ./cmd/main.go

EXPOSE 8085

CMD ["./main"]