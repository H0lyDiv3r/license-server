package main

import (
	"client/internals/auth"
	"client/internals/license"
	"client/internals/state"
	"client/internals/store"
	"context"
	"embed"
	"log"

	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
)

//go:embed all:frontend/dist
var assets embed.FS

func main() {
	// Create an instance of the app structure
	s := &state.AppState{}
	app := NewApp()
	license := license.NewLicense()
	auth := auth.NewAuth()
	_, err := store.NewStore()

	if err != nil {
		log.Fatal("failed to initialize db")
	}

	// Create application with options
	err = wails.Run(&options.App{
		Title:  "client",
		Width:  1024,
		Height: 768,
		AssetServer: &assetserver.Options{
			Assets: assets,
		},
		BackgroundColour: &options.RGBA{R: 255, G: 255, B: 255, A: 1},
		OnStartup: func(ctx context.Context) {
			app.startup(ctx, s)
			license.Startup(ctx, s)
			auth.Startup(ctx, s)
		},
		Bind: []interface{}{
			app,
			license,
			auth,
		},
	})

	if err != nil {
		println("Error:", err.Error())
	}
}
