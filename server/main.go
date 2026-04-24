package main

import (
    "log"
    "net"

    pb "grpc-chat/proto"

    "google.golang.org/grpc"
)

func main() {
    lis, err := net.Listen("tcp", "0.0.0.0:8080")
    if err != nil {
        log.Fatalf("Ошибка запуска сервера: %v", err)
    }

    grpcServer := grpc.NewServer()
    chatServer := NewChatServer()
    pb.RegisterChatServiceServer(grpcServer, chatServer)

    log.Println("gRPC чат сервер запущен на :8080")
    if err := grpcServer.Serve(lis); err != nil {
        log.Fatalf("Ошибка: %v", err)
    }
}