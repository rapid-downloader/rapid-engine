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
	Initter interface {
		Init() error
	}

	Closer interface {
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

	Socket struct {
		Path    string
		Method  string
		Handler func(c *websocket.Conn)
	}

	serviceRunner struct {
		services []Service
		app      *fiber.App
	}

	ServiceFactory func() Service
)

var services = make([]ServiceFactory, 0)

func RegisterService(s ServiceFactory) {
	services = append(services, s)
}

func Create(app *fiber.App) serviceRunner {
	svcs := make([]Service, 0)

	for _, service := range services {
		svcs = append(svcs, service())
	}

	return serviceRunner{
		app:      app,
		services: svcs,
	}
}

func (s *serviceRunner) Run() {
	for _, service := range s.services {
		if init, ok := service.(Initter); ok {
			if err := init.Init(); err != nil {
				log.Fatal(err)
			}
		}

		// create the service
		for _, route := range service.Router() {
			s.app.Add(route.Method, route.Path, route.Handler)
		}
	}
}

func (s *serviceRunner) Shutdown() {
	signals := []os.Signal{syscall.SIGINT, syscall.SIGKILL, syscall.SIGTERM, syscall.SIGSTOP, os.Interrupt}
	ch := make(chan os.Signal, 1)

	signal.Notify(ch, signals...)

	go func() {
		<-ch
		log.Println("Shutting down...")

		defer s.app.Shutdown()

		for _, service := range s.services {
			if closer, ok := service.(Closer); ok {
				if err := closer.Close(); err != nil {
					log.Fatal(err)
				}
			}
		}

		for _, channel := range channels {
			channel.Close()
		}
	}()
}
