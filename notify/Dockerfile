FROM golang:1.22 AS builder
WORKDIR /app

COPY go.mod go.sum ./
RUN go mod tidy

COPY . .
RUN go build -o notify-service ./cmd/app/main.go

FROM gcr.io/distroless/base-debian11
WORKDIR /app

COPY --from=builder /app/notify-service /app/notify-service

COPY .env .env
RUN ls -la .env || echo "File .env not found"

EXPOSE 8080
CMD ["/app/notify-service"]