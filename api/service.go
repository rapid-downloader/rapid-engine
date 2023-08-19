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
var websockets = make([]WebSocket, 0)

func RegisterService(s Service) {
	services = append(services, s)
}

func RegisterWebSocket(r WebSocket) {
	websockets = append(websockets, r)
}

func createService(app *fiber.App) error {
	for _, service := range services {
		if init, ok := service.(ServiceInitter); ok {
			if err := init.Init(); err != nil {
				return err
			}
		}

		for _, route := range service.Router() {
			app.Add(route.Method, route.Path, route.Handler)
		}
	}

	return nil
}

func closeService(app *fiber.App) error {
	for _, service := range services {
		if closer, ok := service.(ServiceCloser); ok {
			if err := closer.Close(); err != nil {
				return err
			}
		}
	}

	return nil
}

func createRealtime(app *fiber.App) error {
	for _, service := range websockets {
		if init, ok := service.(ServiceInitter); ok {
			if err := init.Init(); err != nil {
				return err
			}
		}

		for _, channel := range service.Sockets() {
			app.Add(channel.Method, channel.Path, websocket.New(channel.Handler))
		}
	}

	return nil
}

func closeRealtime(app *fiber.App) error {
	for _, service := range websockets {
		if closer, ok := service.(ServiceCloser); ok {
			if err := closer.Close(); err != nil {
				return err
			}
		}
	}

	return nil
}

func Create(app *fiber.App) {
	if err := createService(app); err != nil {
		log.Fatal("Error creating service:", err)
	}

	if err := createRealtime(app); err != nil {
		log.Fatal("Error creating web socket service:", err)
	}
}

func Shutdown(app *fiber.App) {
	signals := []os.Signal{syscall.SIGINT, syscall.SIGKILL, syscall.SIGTERM, syscall.SIGSTOP}
	ch := make(chan os.Signal, 1)

	signal.Notify(ch, signals...)

	go func() {
		<-ch
		log.Println("Shutting down...")

		if err := closeService(app); err != nil {
			log.Fatal("Error creating service:", err)
		}

		if err := closeRealtime(app); err != nil {
			log.Fatal("Error creating web socket service:", err)
		}

		app.Shutdown()
	}()
}
