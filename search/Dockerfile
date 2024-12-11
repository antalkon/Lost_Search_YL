FROM golang:1.23 AS builder

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN go build -o search_service ./cmd/app/main.go

EXPOSE 9090 8080


CMD ["./search_service"]
