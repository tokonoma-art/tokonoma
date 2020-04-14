package main

import (
	"flag"
	"net/http"
	"path/filepath"

	"github.com/tokonoma-art/tokonoma/pkg/websocket"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

var (
	current  = "default"
	basePath string
)

func init() {
	flag.StringVar(&basePath, "app-path", "..", "path to the app folder")
	flag.Parse()
}

// CurrentArtwork represents an artwork key
type CurrentArtwork struct {
	Artwork string `json:"artwork"`
}

func main() {

	e := echo.New()

	e.Use(middleware.Logger())
	// e.Use(middleware.Recover())

	// Expose home page
	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello from Tokonoma!")
	})

	// Expose client
	e.Pre(middleware.RemoveTrailingSlash())
	e.File("/canvas", filepath.Join(basePath, "client/dist/index.html"))
	e.Static("/canvas", filepath.Join(basePath, "client/dist"))

	// Expose artworks
	e.Static("/artworks", "../storage/artworks")

	// Expose WS canvas API
	pool := websocket.NewCanvasPool()
	go pool.Start()
	e.GET("/ws", func(c echo.Context) (err error) {
		ws, err := websocket.Upgrade(c)
		if err == nil {
			client := &websocket.CanvasClient{Conn: ws}
			pool.Register <- client
			defer func() {
				pool.Unregister <- client
				ws.Close()
			}()
			client.Conn.WriteJSON(websocket.CurrentArtworkCanvasMessage{ArtworkKey: current})
			client.Start()
		}
		return
	})

	// Expose HTTP controller API
	api := e.Group("/api/v1")

	api.GET("/canvases/default", func(c echo.Context) error {
		return c.String(http.StatusOK, current)
	})

	api.POST("/canvases/default/current", func(c echo.Context) (err error) {
		ca := new(CurrentArtwork)
		if err = c.Bind(ca); err != nil {
			return
		}
		current = ca.Artwork
		pool.Broadcast <- websocket.CurrentArtworkCanvasMessage{ArtworkKey: current}
		return c.JSON(http.StatusOK, ca)
	})

	// Start the server
	e.Logger.Fatal(e.Start(":1323"))
}
