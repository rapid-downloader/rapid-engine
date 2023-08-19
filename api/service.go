package api

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/gofiber/contrib/websocket"
	"github.com/gofiber/fiber/v2"
)

type (
	ServiceInitter interface {
		Init() error
	}

	ServiceCloser interface {
		Close() error
	}

	Service interface {
		Router() []Route
	}

	Route struct {
		Method  string
		Path    string
		Handler func(ctx *fiber.Ctx) error
	}
)

var services = make([]Service, 0)

func RegisterService(s Service) {
	services = append(services, s)
}

func Create(app *fiber.App) {
	for _, service := range services {
		if init, ok := service.(ServiceInitter); ok {
			if err := init.Init(); err != nil {
				log.Fatal("Error initiating service:", err)
			}
		}

		for _, route := range service.Router() {
			app.Add(route.Method, route.Path, route.Handler)
		}

		if ch, ok := service.(Realtime); ok {
			for _, channel := range ch.Channels() {
				app.Add(channel.Method, channel.Path, websocket.New(channel.Handler))
			}
		}
	}
}

func Shutdown(app *fiber.App) {
	signals := []os.Signal{syscall.SIGINT, syscall.SIGKILL, syscall.SIGTERM, syscall.SIGSTOP}
	ch := make(chan os.Signal, 1)

	signal.Notify(ch, signals...)

	go func() {
		<-ch
		log.Println("Shutting down...")

		for _, service := range services {
			if closer, ok := service.(ServiceCloser); ok {
				if err := closer.Close(); err != nil {
					log.Fatal("Error closing service:", err)
				}
			}
		}

		app.Shutdown()
	}()
}
