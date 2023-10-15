package api

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/rapid-downloader/rapid/api"
	"github.com/rapid-downloader/rapid/db"
	"github.com/rapid-downloader/rapid/entry"
	response "github.com/rapid-downloader/rapid/helper"
)

const (
	ClientCLI = "cli"
	ClientGUI = "gui"
)

type entryService struct {
	channel api.Channel
	store   Store
}

func newService() api.Service {
	return &entryService{
		channel: api.CreateChannel("memstore"),
		store:   NewStore("download", db.DB()),
	}
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

	s.channel.Publish(entry)

	chunkProgress := make([]int, entry.ChunkLen())
	for i := range chunkProgress {
		chunkProgress[i] = 0
	}

	toDownload := Download{
		ID:            entry.ID(),
		Name:          entry.Name(),
		URL:           entry.URL(),
		Size:          entry.Size(),
		Type:          entry.Type(),
		ChunkLen:      entry.ChunkLen(),
		Provider:      entry.Downloader(),
		Resumable:     entry.Resumable(),
		Progress:      0,
		ChunkProgress: chunkProgress,
		TimeLeft:      time.Time{},
		Speed:         0,
		Status:        "Queued",
		Date:          time.Now(),
	}

	if err := s.store.Create(entry.ID(), toDownload); err != nil {
		return response.Error(ctx, err.Error(), fiber.StatusInternalServerError)
	}

	return response.Success(ctx, toDownload)
}

func (s *entryService) getEntry(ctx *fiber.Ctx) error {
	id := ctx.Params("id")
	res := s.store.Get(id)
	if res == nil {
		return response.NotFound(ctx)
	}

	return response.Success(ctx, res)
}

func (s *entryService) getAllEntry(ctx *fiber.Ctx) error {
	res := s.store.GetAll()
	if res == nil {
		return response.NotFound(ctx)
	}

	return response.Success(ctx, res)
}

func (s *entryService) deleteEntry(ctx *fiber.Ctx) error {
	id := ctx.Params("id")

	if err := s.store.Delete(id); err != nil {
		return response.NotFound(ctx)
	}

	return response.Success(ctx)
}

func (s *entryService) Routes() []api.Route {
	return []api.Route{
		{
			Path:    "/fetch",
			Method:  "POST",
			Handler: s.fetch,
		},
		{
			Path:    "/entry/:id",
			Method:  "GET",
			Handler: s.getEntry,
		},
		{
			Path:    "/entries",
			Method:  "GET",
			Handler: s.getAllEntry,
		},
		{
			Path:    "/delete/entry/:id",
			Method:  "DELETE",
			Handler: s.deleteEntry,
		},
	}
}

func init() {
	api.RegisterService(newService)
}
