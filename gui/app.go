package main

import (
	"context"

	"github.com/rapid-downloader/rapid/client"
	"github.com/rapid-downloader/rapid/log"
	"github.com/wailsapp/wails/v2/pkg/runtime"
)

// App struct
type App struct {
	ctx   context.Context
	rapid client.Rapid
}

// NewApp creates a new App application struct
func NewApp() *App {
	return &App{}
}

// startup is called when the app starts. The context is saved
// so we can call the runtime methods
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
	rapid, err := client.New(ctx, "gui")
	if err != nil {
		log.Fatal(err)
	}

	go rapid.Listen(func(progress client.Progress, err error) {
		if err != nil {
			log.Println(err)
			return
		}

		runtime.EventsEmit(ctx, "progress", progress)
	})

	a.rapid = rapid
}

func (a *App) shutdown(ctx context.Context) {
	if closer, ok := a.rapid.(client.RapidCloser); ok {
		closer.Close()
	}
}

func (a *App) Fetch(request client.Request) (client.Download, error) {
	res, err := a.rapid.Fetch(request)
	if err != nil {
		// TODO: add notification
		log.Println(err)
		return client.Download{}, nil
	}

	return *res, nil
}

func (a *App) Download(id string) error {
	return a.rapid.Download(id)
}
