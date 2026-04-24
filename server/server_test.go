package main

import (
    "testing"
    pb "grpc-chat/proto"
)

func TestNewChatServer(t *testing.T) {
    server := NewChatServer()
    if server.streams == nil {
        t.Error("Expected streams map to be initialized")
    }
    if len(server.messages) != 0 {
        t.Error("Expected empty messages slice")
    }
}