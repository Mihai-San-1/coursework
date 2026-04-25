# Этап 1: сборка бинарника
FROM golang:1.24 AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o chat-server ./cmd/server

# Этап 2: минимальный образ
FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder /app/chat-server .
COPY --from=builder /app/public ./public
EXPOSE 8080
CMD ["./chat-server"]