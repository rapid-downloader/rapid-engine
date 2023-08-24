package entry

import (
	"github.com/gofiber/fiber/v2"
	"github.com/rapid-downloader/rapid/api"
	"github.com/rapid-downloader/rapid/downloader"
	"github.com/rapid-downloader/rapid/entry"
	response "github.com/rapid-downloader/rapid/helper"
)

const (
	ClientCLI = "cli"
	ClientGUI = "gui"
)

type entryService struct {
	entries *entry.Listing
}

func NewService(entries *entry.Listing) api.Service {
	return &entryService{
		entries: entries,
	}
}

func (s *entryService) Init() error {
	return nil
}

func (s *entryService) fetch(ctx *fiber.Ctx) error {
	provider := ctx.Params("provider", downloader.Default)

	var req request

	if err := ctx.BodyParser(&req); err != nil {
		return response.Error(ctx, err.Error())
	}

	options := req.toOptions()
	options = append(options, entry.UseDownloadProvider(provider))

	entry, err := entry.Fetch(req.Url, options...)
	if err != nil {
		return response.Error(ctx, err.Error())
	}

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
	api.RegisterService(NewService)
}
