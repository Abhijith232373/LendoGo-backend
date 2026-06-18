package websockets

import (
	"log"
	"sync"
	"github.com/gofiber/websocket/v2"
)

var Clients = make(map[*websocket.Conn]bool)
var broadcast = make(chan interface{}, 100)
var Mutex = sync.Mutex{}

func StartHub() {
	for {
		msg := <-broadcast
		Mutex.Lock()
		for client := range Clients {
			err := client.WriteJSON(msg)
			if err != nil {
				log.Printf("WebSocket Error: %v", err)
				client.Close()
				delete(Clients, client)
			}
		}
		Mutex.Unlock()
	}
}

func BroadcastMessage(eventType string, data interface{}) {
	broadcast <- map[string]interface{}{
		"event": eventType,
		"data":  data,
	}
}