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
	memstore entry.Store
}

func newService(memstore entry.Store) api.Service {
	return &entryService{
		memstore: memstore,
	}
}

func (s *entryService) Init() error {
	return nil
}

func (s *entryService) fetch(ctx *fiber.Ctx) error {
	var req request

	if err := ctx.BodyParser(&req); err != nil {
		return response.Error(ctx, err.Error())
	}

	entry, err := entry.Fetch(req.Url, req.toOptions()...)
	if err != nil {
		return response.Error(ctx, err.Error())
	}

	s.memstore.Set(entry.ID(), entry)

	return response.Success(ctx, entry)
}

func (s *entryService) Router() []api.Route {
	return []api.Route{
		{
			Path:    "/fetch",
			Method:  "POST",
			Handler: s.fetch,
		},
	}
}

func init() {
	api.RegisterService(newService)
}
