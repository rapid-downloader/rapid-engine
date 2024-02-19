package api

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/gofiber/contrib/websocket"
	"github.com/gofiber/fiber/v2"
	"github.com/rapid-downloader/rapid/env"
)

type (
	App interface {
		Run()
		Shutdown()
	}

	Service interface {
		CreateRoutes()
	}

	ServiceInitter interface {
		Init() error
	}

	ServiceCloser interface {
		Close() error
	}

	ServiceFactory func(app *fiber.App) Service

	Socket struct {
		Path    string
		Method  string
		Handler func(c *websocket.Conn)
	}

	service struct {
		services []Service
		app      *fiber.App
	}
)

var services = make([]ServiceFactory, 0)

func RegisterService(s ServiceFactory) {
	services = append(services, s)
}

func Create(app *fiber.App) App {
	svcs := make([]Service, 0)

	for _, service := range services {
		svcs = append(svcs, service(app))
	}

	return &service{
		app:      app,
		services: svcs,
	}
}

func (s *service) Run() {
	for _, service := range s.services {
		if init, ok := service.(ServiceInitter); ok {
			if err := init.Init(); err != nil {
				log.Fatal(err)
			}
		}

		// create the service
		service.CreateRoutes()
	}
}

func (s *service) Shutdown() {
	signals := []os.Signal{syscall.SIGINT, syscall.SIGKILL, syscall.SIGTERM, syscall.SIGSTOP, os.Interrupt}
	ch := make(chan os.Signal, 1)

	signal.Notify(ch, signals...)

	go func() {
		<-ch
		log.Println("Shutting down...")

		defer s.app.Shutdown()

		for _, service := range s.services {
			if closer, ok := service.(ServiceCloser); ok {
				if err := closer.Close(); err != nil {
					log.Fatal(err)
				}
			}
		}

		for _, channel := range channels {
			channel.Close()
		}
	}()

	port := env.Get("API_PORT").String(":9999")

	if err := s.app.Listen(port); err != nil {
		log.Fatal(err)
	}
}
