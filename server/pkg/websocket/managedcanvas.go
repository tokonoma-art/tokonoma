package websocket

import (
	"fmt"

	"github.com/tokonoma-art/tokonoma/pkg/canvas"
)

// ManagedCanvas …
type ManagedCanvas struct {
	Canvas     canvas.Canvas
	ClientPool *CanvasClientPool
}

// NewManagedCanvas …
func NewManagedCanvas(key string) *ManagedCanvas {
	pool := NewCanvasClientPool()
	go pool.Start()
	return &ManagedCanvas{
		Canvas: canvas.Canvas{
			Key:       key,
			Rotation:  canvas.Rotation0,
			Artbundle: "one.artbundle",
		},
		ClientPool: pool,
	}
}

// BroadcastUpdate broadcast current settings to all clients
func (managedCanvas *ManagedCanvas) BroadcastUpdate() {
	managedCanvas.ClientPool.Broadcast <- managedCanvas.Canvas
}

// UpdateClient sends the current settings to the given client
func (managedCanvas *ManagedCanvas) UpdateClient(client *CanvasClient) {
	if err := client.Conn.WriteJSON(managedCanvas.Canvas); err != nil {
		fmt.Println(err)
		return
	}
}
