package main

import (
	"embed"
	"io/fs"
	"net/http"

	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
)

//go:embed all:frontend/dist
var assets embed.FS

//go:embed phase0/*.html
var probeAssets embed.FS

func main() {
	// Create an instance of the app structure
	app := NewApp()
	probeFiles, err := fs.Sub(probeAssets, "phase0")
	if err != nil {
		panic(err)
	}
	probeHandler := http.StripPrefix("/phase0/", http.FileServer(http.FS(probeFiles)))

	// Create application with options
	err = wails.Run(&options.App{
		Title:  "OI-Desk",
		Width:  1024,
		Height: 768,
		AssetServer: &assetserver.Options{
			Assets:  assets,
			Handler: probeHandler,
		},
		BackgroundColour: &options.RGBA{R: 27, G: 38, B: 54, A: 1},
		OnStartup:        app.startup,
		Bind: []interface{}{
			app,
		},
	})

	if err != nil {
		println("Error:", err.Error())
	}
}
