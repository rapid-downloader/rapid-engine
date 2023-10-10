package api

import (
	"github.com/gofiber/fiber/v2"
	"github.com/rapid-downloader/rapid/api"
	"github.com/rapid-downloader/rapid/entry"
	response "github.com/rapid-downloader/rapid/helper"
)

const (
	ClientCLI = "cli"
	ClientGUI = "gui"
)

type entryService struct {
	channel api.Channel
}

func newService() api.Service {
	return &entryService{
		channel: api.CreateChannel("memstore"),
	}
}

func (s *entryService) Close() error {
	return s.channel.Close()
}

func (s *entryService) fetch(ctx *fiber.Ctx) error {
	var req request

	if err := ctx.BodyParser(&req); err != nil {
		return response.Error(ctx, err.Error())
	}

	client := ctx.Params("client")
	if client != "browser" {
		//TODO: get cookies etc with headless browser
	}

	entry, err := entry.Fetch(req.Url, req.toOptions()...)
	if err != nil {
		return response.Error(ctx, err.Error())
	}

	s.channel.Publish(entry)

	return response.Success(ctx, entry)
}

func (s *entryService) Router() []api.Route {
	return []api.Route{
		{
			Path:    "/:client/fetch",
			Method:  "POST",
			Handler: s.fetch,
		},
	}
}

func init() {
	api.RegisterService(newService)
}
