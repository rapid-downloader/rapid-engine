package api

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/rapid-downloader/rapid/api"
	"github.com/rapid-downloader/rapid/db"
	"github.com/rapid-downloader/rapid/entry"
	response "github.com/rapid-downloader/rapid/helper"
	"github.com/rapid-downloader/rapid/log"
	"github.com/rapid-downloader/rapid/setting"
	"github.com/rapid-downloader/rapid/utils"
)

const (
	ClientCLI = "cli"
	ClientGUI = "gui"
)

type entryService struct {
	app     *fiber.App
	channel api.Channel
	store   Store
}

func newService(app *fiber.App) api.Service {
	return &entryService{
		app:     app,
		channel: api.CreateChannel("memstore"),
		store:   NewStore("download", db.DB()),
	}
}

func (s *entryService) fetch(ctx *fiber.Ctx) error {
	var req request

	if err := ctx.BodyParser(&req); err != nil {
		return response.BadRequest(ctx, err)
	}

	entry, err := entry.Fetch(req.Url, req.toOptions()...)
	if err != nil {
		return response.BadRequest(ctx, err)
	}

	s.channel.Publish(entry)

	toDownload := Download{
		ID:               entry.ID(),
		Name:             entry.Name(),
		Location:         entry.Location(),
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
		return response.InternalServerError(ctx, err)
	}

	return response.Ok(ctx, toDownload)
}

func (s *entryService) getEntry(ctx *fiber.Ctx) error {
	id := ctx.Params("id")
	res := s.store.Get(id)
	if res == nil {
		return response.Success(ctx, fiber.StatusNoContent)
	}

	return response.Ok(ctx, res)
}

func (s *entryService) getAllEntry(ctx *fiber.Ctx) error {
	setting := setting.Get()
	page := utils.Parse(ctx.Query("page")).Int(1)

	res := s.store.GetAll(page, setting.DisplayedEntriesCount)
	if res == nil {
		return response.Success(ctx, fiber.StatusNoContent)
	}

	return response.Ok(ctx, res)
}

func (s *entryService) updateEntry(ctx *fiber.Ctx) error {
	id := ctx.Params("id")

	var payload UpdateDownload
	if err := ctx.BodyParser(&payload); err != nil {
		return response.BadRequest(ctx, err)
	}

	if err := s.store.Update(id, payload); err != nil {
		return response.BadRequest(ctx, err)
	}

	return response.Ok(ctx)
}

func (s *entryService) updateAllEntry(ctx *fiber.Ctx) error {
	var payload BatchUpdateDownload
	if err := ctx.BodyParser(&payload); err != nil {
		return response.BadRequest(ctx, err)
	}

	if err := s.store.BatchUpdate(payload.IDs, payload.Payload); err != nil {
		return response.BadRequest(ctx, err)
	}

	return response.Ok(ctx)
}

func (s *entryService) deleteEntry(ctx *fiber.Ctx) error {
	id := ctx.Params("id")
	fromDisk := ctx.QueryBool("fromDisk", false)

	entry := s.store.Get(id)
	if entry == nil {
		return response.Success(ctx, fiber.StatusNoContent)
	}

	if err := s.store.Delete(id); err != nil {
		return response.Success(ctx, fiber.StatusNoContent)
	}

	if !fromDisk {
		return response.Ok(ctx)
	}

	for i := 0; i < entry.ChunkLen; i++ {
		dirpath := strings.Replace(entry.Location, filepath.Base(entry.Location), "", 1)
		path := fmt.Sprintf("%s/%s-%d", dirpath, entry.ID, i)

		if err := os.Remove(path); err != nil {
			log.Println(err)
		}
	}

	if err := os.Remove(entry.Location); err != nil {
		log.Println(err)
	}

	return response.Ok(ctx)
}

func (s *entryService) CreateRoutes() {
	s.app.Add("POST", "/fetch", s.fetch)

	s.app.Add("GET", "/entries/:id", s.getEntry)
	s.app.Add("GET", "/entries", s.getAllEntry)
	s.app.Add("PUT", "/entries/:id", s.updateEntry)
	s.app.Add("PUT", "/entries", s.updateAllEntry)
	s.app.Add("DELETE", "/entries/:id", s.deleteEntry)
}

func init() {
	api.RegisterService(newService)
}
