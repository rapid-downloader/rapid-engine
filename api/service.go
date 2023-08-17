package api

import (
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

func Services() []Service {
	return services
}
