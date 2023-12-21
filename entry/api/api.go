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

	toDownload := Download{
		ID:               entry.ID(),
		Name:             entry.Name(),
		URL:              entry.URL(),
		Size:             entry.Size(),
		Type:             entry.Type(),
		ChunkLen:         entry.ChunkLen(),
		Provider:         entry.Downloader(),
		Resumable:        entry.Resumable(),
		Progress:         0,
		DownloadedChunks: make([]int64, entry.ChunkLen()),
		TimeLeft:         0,
		Speed:            0,
		Status:           "Queued",
		Date:             time.Now(),
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

func (s *entryService) updateEntry(ctx *fiber.Ctx) error {
	id := ctx.Params("id")

	var payload UpdateDownload
	if err := ctx.BodyParser(&payload); err != nil {
		return response.Error(ctx, err.Error())
	}

	if err := s.store.Update(id, payload); err != nil {
		return response.Error(ctx, err.Error())
	}

	return response.Success(ctx)
}

func (s *entryService) updateAllEntry(ctx *fiber.Ctx) error {
	var payload BatchUpdateDownload
	if err := ctx.BodyParser(&payload); err != nil {
		return response.Error(ctx, err.Error())
	}

	if err := s.store.BatchUpdate(payload.IDs, payload.Payload); err != nil {
		return response.Error(ctx, err.Error())
	}

	return response.Success(ctx)
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
			Path:    "/entries/:id",
			Method:  "GET",
			Handler: s.getEntry,
		},
		{
			Path:    "/entries",
			Method:  "GET",
			Handler: s.getAllEntry,
		},
		{
			Path:    "/entries/:id",
			Method:  "PUT",
			Handler: s.updateEntry,
		},
		{
			Path:    "/entries",
			Method:  "PUT",
			Handler: s.updateAllEntry,
		},
		{
			Path:    "/delete/entries/:id",
			Method:  "DELETE",
			Handler: s.deleteEntry,
		},
	}
}

func init() {
	api.RegisterService(newService)
}
