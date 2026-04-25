package main

import (
	"log"
	"net/http"

	"web-chat/internal/handler"
)

func main() {
	http.Handle("/", http.FileServer(http.Dir("./public")))
	http.HandleFunc("/ws", handler.HandleConnections)

	log.Println("✅ Сервер запущен на http://localhost:8080")
	log.Println("📱 Для доступа с телефона: http://<ВАШ_IP>:8080")
	
	err := http.ListenAndServe("0.0.0.0:8080", nil)
	if err != nil {
		log.Fatal("❌ Ошибка:", err)
	}
}