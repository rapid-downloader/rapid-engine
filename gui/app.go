package main

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/rapid-downloader/rapid/client"
	"github.com/rapid-downloader/rapid/client/websocket"
	"github.com/rapid-downloader/rapid/env"
	"github.com/rapid-downloader/rapid/log"
	"github.com/wailsapp/wails/v2/pkg/runtime"
)

// App struct
type App struct {
	ctx context.Context
	ws  websocket.Websocket
}

// NewApp creates a new App application struct
func NewApp() *App {
	return &App{}
}

// startup is called when the app starts. The context is saved
// so we can call the runtime methods
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
	host := env.Get("API_HOST").String("localhost")
	port := env.Get("API_PORT").String(":9999")
	id := "gui"

	wsUrl := fmt.Sprintf("ws://%s%s/ws/%s", host, port, id)

	a.ws = websocket.Connect(ctx, wsUrl)
	go a.ws.Listen(func(msg []byte) {
		var progress client.Progress
		if err := json.Unmarshal(msg, &progress); err != nil {
			log.Println("error unmarshalling progress bar:", err)
			return
		}

		runtime.EventsEmit(ctx, "progress", progress)
	})
}

func (a *App) shutdown(ctx context.Context) {
	a.ws.Close()
}
