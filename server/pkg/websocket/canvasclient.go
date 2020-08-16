package websocket

import (
	"fmt"

	"github.com/gorilla/websocket"
	"github.com/tokonoma-art/tokonoma/pkg/canvas"
)

// CanvasClient …
type CanvasClient struct {
	Conn *websocket.Conn
}

// Start …
func (client *CanvasClient) Start() {
	for {
		_, p, err := client.Conn.ReadMessage()
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Println(string(p))
	}
}

// CanvasClientPool …
type CanvasClientPool struct {
	Register   chan *CanvasClient
	Unregister chan *CanvasClient
	Clients    map[*CanvasClient]bool
	Broadcast  chan canvas.Canvas
}

// NewCanvasClientPool …
func NewCanvasClientPool() *CanvasClientPool {
	return &CanvasClientPool{
		Register:   make(chan *CanvasClient),
		Unregister: make(chan *CanvasClient),
		Clients:    make(map[*CanvasClient]bool),
		Broadcast:  make(chan canvas.Canvas),
	}
}

// Start …
func (pool *CanvasClientPool) Start() {
	for {
		select {
		case client := <-pool.Register:
			pool.Clients[client] = true
			fmt.Println("Size of Connection Pool: ", len(pool.Clients))
		case client := <-pool.Unregister:
			delete(pool.Clients, client)
			fmt.Println("Size of Connection Pool: ", len(pool.Clients))
		case message := <-pool.Broadcast:
			for client := range pool.Clients {
				if err := client.Conn.WriteJSON(message); err != nil {
					fmt.Println(err)
					return
				}
			}
		}
	}
}
