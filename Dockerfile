FROM golang:1.24-alpine AS builder

WORKDIR /app

# Установка protoc
RUN apk add --no-cache protobuf protobuf-dev

# Копирование исходников
COPY go.mod go.sum ./
RUN go mod download

COPY . .

# Генерация proto
RUN go install google.golang.org/protobuf/cmd/protoc-gen-go@latest && \
    go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest && \
    export PATH=$PATH:$(go env GOPATH)/bin && \
    protoc --go_out=. --go_opt=paths=source_relative \
           --go-grpc_out=. --go-grpc_opt=paths=source_relative \
           proto/chat.proto

# Сборка сервера
RUN go build -o chat-server server/*.go

# Финальный образ
FROM alpine:latest

RUN apk --no-cache add ca-certificates

WORKDIR /root/

COPY --from=builder /app/chat-server .

EXPOSE 50051

CMD ["./chat-server"]