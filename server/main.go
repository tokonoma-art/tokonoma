package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"path/filepath"
	"strings"

	"github.com/tokonoma-art/tokonoma/pkg/canvas"
	"github.com/tokonoma-art/tokonoma/pkg/websocket"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

const artworksDirectory = "artworks"

var (
	managedCanvases = make(map[string]*websocket.ManagedCanvas) // to ease future support for multiple canvases
	appPath         string
	storagePath     string
	artworksPath    string
)

func init() {
	flag.StringVar(&appPath, "app", "..", "path to the app folder")
	flag.StringVar(&storagePath, "storage", "../storage", "path to the storage folder")
	flag.Parse()
	artworksPath = filepath.Join(storagePath, artworksDirectory)
}

// ArtbundleSetting represents the artbundle path setting
type ArtbundleSetting struct {
	Artbundle string `json:"artbundle"`
}

// RotationSetting represents a rotation setting
type RotationSetting struct {
	Rotation canvas.Rotation `json:"rotation"`
}

// ExtractManagedCanvas is an Echo middleware that puts the managedCanvas in the context
func ExtractManagedCanvas(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		key := c.Param("key")
		if managedCanvas, ok := managedCanvases[key]; ok {
			c.Set("managedCanvas", managedCanvas)
			return next(c)
		}
		return c.String(http.StatusNotFound, fmt.Sprintf("Cannot find canvas '%s'.", key))
	}
}

func main() {

	// Create default canvas manually
	managedCanvases["default"] = websocket.NewManagedCanvas("default")

	// Configure the HTTP server…

	e := echo.New()
	e.Use(middleware.Logger())
	// e.Use(middleware.Recover())

	// Expose home page
	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello from Tokonoma!")
	})

	// Expose artworks
	e.Static("/artworks", artworksPath)

	// Expose client
	e.Pre(middleware.RemoveTrailingSlash())
	e.File("/canvas", filepath.Join(appPath, "client/dist/index.html"))
	e.Static("/canvas", filepath.Join(appPath, "client/dist"))

	// Expose WS canvas API
	// at the moment: "default" canvas is hard-coded
	defaultManagedCanvas := managedCanvases["default"]
	e.GET("/ws", func(c echo.Context) (err error) {
		ws, err := websocket.Upgrade(c)
		if err == nil {
			client := &websocket.CanvasClient{Conn: ws}
			defaultManagedCanvas.ClientPool.Register <- client
			defer func() {
				defaultManagedCanvas.ClientPool.Unregister <- client
				ws.Close()
			}()
			defaultManagedCanvas.UpdateClient(client)
			client.Start()
		}
		return
	})

	// Expose HTTP controller API
	api := e.Group("/api/v1")

	// Lists all available artworks
	api.GET("/artworks", func(c echo.Context) (err error) {
		files, err := ioutil.ReadDir(artworksPath)
		if err != nil {
			return
		}
		artworks := []string{}
		for _, f := range files {
			if name := f.Name(); strings.HasSuffix(name, ".artbundle") {
				artworks = append(artworks, name)
			}
		}
		return c.JSON(http.StatusOK, artworks)
	})

	// Lists all canvases
	api.GET("/canvases", func(c echo.Context) error {
		keys := []string{}
		for key := range managedCanvases {
			keys = append(keys, key)
		}
		return c.JSON(http.StatusOK, keys)
	})

	// API about one specific canvas
	canvasAPI := api.Group("/canvases/:key", ExtractManagedCanvas)

	// Returns canvas' current settings
	canvasAPI.GET("", func(c echo.Context) error {
		managedCanvas := c.Get("managedCanvas").(*websocket.ManagedCanvas)
		return c.JSON(http.StatusOK, managedCanvas.Canvas)
	})

	// Sets the canvas' artbundle
	canvasAPI.POST("/artbundle", func(c echo.Context) (err error) {
		managedCanvas := c.Get("managedCanvas").(*websocket.ManagedCanvas)
		setting := new(ArtbundleSetting)
		if err = c.Bind(setting); err != nil {
			return
		}
		managedCanvas.Canvas.Artbundle = setting.Artbundle
		managedCanvas.BroadcastUpdate()
		return c.JSON(http.StatusOK, setting)
	})

	// Sets the canvas' rotation
	canvasAPI.POST("/rotation", func(c echo.Context) (err error) {
		managedCanvas := c.Get("managedCanvas").(*websocket.ManagedCanvas)
		setting := new(RotationSetting)
		if err = c.Bind(setting); err != nil {
			return
		}
		managedCanvas.Canvas.Rotation = setting.Rotation
		managedCanvas.BroadcastUpdate()
		return c.JSON(http.StatusOK, setting)
	})

	// Start the server
	e.Logger.Fatal(e.Start(":1323"))
}
