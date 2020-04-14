package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"path/filepath"

	"github.com/gorilla/websocket"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

var (
	current  = "default"
	upgrader = websocket.Upgrader{}
	lastWS   *websocket.Conn
)

var basePath string

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

	// Expose API
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
		if lastWS != nil {
			lastWS.WriteJSON(ca)
		}
		return c.JSON(http.StatusOK, ca)
	})

	e.GET("/ws", wsAPI)

	// Start the server
	e.Logger.Fatal(e.Start(":1323"))
}

func wsAPI(c echo.Context) (err error) {
	upgrader.CheckOrigin = func(r *http.Request) bool { return true }
	ws, err := upgrader.Upgrade(c.Response(), c.Request(), nil)
	if err != nil {
		return
	}
	defer ws.Close()

	log.Println("Client Connected")
	lastWS = ws
	reader(ws)

	return nil
}

func reader(conn *websocket.Conn) {
	for {
		// read in a message
		_, p, err := conn.ReadMessage()
		if err != nil {
			log.Println(err)
			return
		}
		// print out that message for clarity
		fmt.Println(string(p))

	}
}
