package websocket

import (
	"fmt"

	"github.com/gorilla/websocket"
)

// CurrentArtworkCanvasMessage …
type CurrentArtworkCanvasMessage struct {
	ArtworkKey string `json:"artworkKey"`
}

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

// CanvasPool …
type CanvasPool struct {
	Register   chan *CanvasClient
	Unregister chan *CanvasClient
	Clients    map[*CanvasClient]bool
	Broadcast  chan CurrentArtworkCanvasMessage
}

// NewCanvasPool …
func NewCanvasPool() *CanvasPool {
	return &CanvasPool{
		Register:   make(chan *CanvasClient),
		Unregister: make(chan *CanvasClient),
		Clients:    make(map[*CanvasClient]bool),
		Broadcast:  make(chan CurrentArtworkCanvasMessage),
	}
}

// Start …
func (pool *CanvasPool) Start() {
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
