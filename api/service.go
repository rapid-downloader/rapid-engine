package api

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/gofiber/contrib/websocket"
	"github.com/gofiber/fiber/v2"
	"github.com/rapid-downloader/rapid/entry"
	"github.com/rapid-downloader/rapid/logger"
	"github.com/rapid-downloader/rapid/setting"
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
		lists    *entry.Listing
		services []Service
		app      *fiber.App
	}

	ServiceFactory func(entries *entry.Listing) Service
)

var services = make([]ServiceFactory, 0)

func RegisterService(s ServiceFactory) {
	services = append(services, s)
}

func Create(app *fiber.App) serviceRunner {
	setting := setting.Get()
	logger := logger.New(logger.FS, setting)

	svcs := make([]Service, 0)
	lists := entry.NewListing(setting, logger)

	for _, service := range services {
		svcs = append(svcs, service(lists))
	}

	return serviceRunner{
		lists:    lists,
		app:      app,
		services: svcs,
	}
}

func (s *serviceRunner) Run() {
	if listInitter, ok := s.lists.List.(entry.ListInitter); ok {
		if err := listInitter.Init(); err != nil {
			log.Fatal(err)
		}
	}

	if listInitter, ok := s.lists.Queue.(entry.ListInitter); ok {
		if err := listInitter.Init(); err != nil {
			log.Fatal(err)
		}
	}

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

		if listCloser, ok := s.lists.List.(entry.ListCloser); ok {
			if err := listCloser.Close(); err != nil {
				log.Fatal(err)
			}
		}

		if listCloser, ok := s.lists.Queue.(entry.ListCloser); ok {
			if err := listCloser.Close(); err != nil {
				log.Fatal(err)
			}
		}
	}()
}
