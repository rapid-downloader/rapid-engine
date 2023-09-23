package api

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/rapid-downloader/rapid/api"
	"github.com/rapid-downloader/rapid/entry"
	response "github.com/rapid-downloader/rapid/helper"
	"github.com/rapid-downloader/rapid/logger"
	"github.com/rapid-downloader/rapid/setting"
)

const (
	ClientCLI = "cli"
	ClientGUI = "gui"
)

type entryService struct {
	entries *entry.Listing
}

func newService(entries *entry.Listing) api.Service {
	return &entryService{
		entries: entries,
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

	logger := logger.New(logger.FS, setting.Get())

	options := req.toOptions()
	options = append(options, entry.UseLogger(logger))

	entry, err := entry.Fetch(req.Url, options...)
	if err != nil {
		return response.Error(ctx, err.Error())
	}

	s.entries.List.Insert(entry)

	return response.Success(ctx, entry)
}

func (s *entryService) queue(ctx *fiber.Ctx) error {
	var req queueRequest

	if err := ctx.BodyParser(&req); err != nil {
		return response.Error(ctx, err.Error())
	}

	entries := make([]entry.Entry, len(req.Requests))
	for i, request := range req.Requests {
		entry, err := entry.Fetch(request.Url, request.toOptions()...)
		if err != nil {
			return response.Error(ctx, fmt.Sprint("Error fetching url:", entry))
		}

		entries[i] = entry
		s.entries.Queue.Push(entry)
	}

	return response.Success(ctx, entries)
}

func (s *entryService) Router() []api.Route {
	return []api.Route{
		{
			Path:    "/fetch",
			Method:  "POST",
			Handler: s.fetch,
		},
		{
			Path:    "/queue",
			Method:  "POST",
			Handler: s.queue,
		},
	}
}

func init() {
	api.RegisterService(newService)
}
