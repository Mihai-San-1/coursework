package handler

import (
    "net/http"
    "sync"
    "web-chat/internal/message"
    "github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{CheckOrigin: func(r *http.Request) bool { return true }}
var clients = make(map[*websocket.Conn]bool)
var broadcast = make(chan message.Message)
var mutex = sync.Mutex{}

func init() {
    go message.HandleMessages(broadcast, clients, &mutex)
}

func HandleConnections(w http.ResponseWriter, r *http.Request) {
    // ... код из твоего main.go
}