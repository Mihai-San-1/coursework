package message

import (
	"log"
	"sync"

	"github.com/gorilla/websocket"
)

type Message struct {
	User string `json:"user"`
	Text string `json:"text"`
}

func HandleMessages(broadcast chan Message, clients map[*websocket.Conn]bool, mutex *sync.Mutex) {
	for {
		msg := <-broadcast
		mutex.Lock()
		for client := range clients {
			err := client.WriteJSON(msg)
			if err != nil {
				log.Printf("Ошибка отправки: %v", err)
				client.Close()
				delete(clients, client)
			}
		}
		mutex.Unlock()
	}
}