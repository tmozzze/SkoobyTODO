# Build
FROM golang:1.24.4-alpine AS builder

WORKDIR /app

COPY go.mod ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o skooby-todo ./cmd/app/main.go

# Run
FROM alpine:latest

WORKDIR /root/

COPY --from=builder /app/skooby-todo .

COPY .env .

EXPOSE 8080

CMD ["./skooby-todo"]