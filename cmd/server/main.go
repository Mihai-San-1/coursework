package main

import (
    "log"
    "net/http"
    "web-chat/internal/handler"
)

func main() {
    fs := http.FileServer(http.Dir("./public"))
    http.Handle("/", fs)
    http.HandleFunc("/ws", handler.HandleConnections)

    log.Println("✅ Сервер на порту 8080")
    log.Fatal(http.ListenAndServe("0.0.0.0:8080", nil))
}