package handler

import (
	"log"
	"net/http"
	"sync"

	"web-chat/internal/message"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

var (
	clients   = make(map[*websocket.Conn]bool)
	broadcast = make(chan message.Message)
	mutex     = sync.Mutex{}
)

func init() {
	go message.HandleMessages(broadcast, clients, &mutex)
}

func HandleConnections(w http.ResponseWriter, r *http.Request) {
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("Ошибка апгрейда: %v", err)
		return
	}
	defer ws.Close()

	mutex.Lock()
	clients[ws] = true
	mutex.Unlock()

	log.Printf("✅ Клиент подключен. Всего: %d", len(clients))

	for {
		var msg message.Message
		err := ws.ReadJSON(&msg)
		if err != nil {
			log.Printf("Клиент отключился: %v", err)
			mutex.Lock()
			delete(clients, ws)
			mutex.Unlock()
			break
		}
		broadcast <- msg
	}
}