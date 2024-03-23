package api

import (
	"fmt"
	logger "log"
	"time"

	rapidClient "github.com/rapid-downloader/rapid/client"
	"github.com/rapid-downloader/rapid/log"

	"github.com/goccy/go-json"

	"github.com/gofiber/contrib/websocket"
	"github.com/gofiber/fiber/v2"
	"github.com/rapid-downloader/rapid/api"
	"github.com/rapid-downloader/rapid/db"
	"github.com/rapid-downloader/rapid/downloader"
	"github.com/rapid-downloader/rapid/entry"
	entryApi "github.com/rapid-downloader/rapid/entry/api"
	response "github.com/rapid-downloader/rapid/helper"
	"github.com/rapid-downloader/rapid/setting"
)

type downloaderService struct {
	app      *fiber.App
	memstore entry.Store
	channel  api.Channel
	store    entryApi.Store
}

func newService(app *fiber.App) api.Service {
	return &downloaderService{
		app:      app,
		memstore: entry.Memstore(),
		channel:  api.CreateChannel("memstore"),
		store:    entryApi.NewStore("download", db.DB()),
	}
}

func (s *downloaderService) Init() error {
	go s.channel.Subscribe(func(data interface{}) {
		entry, ok := data.(entry.Entry)
		if !ok {
			return
		}

		if err := s.memstore.Set(entry.ID(), entry); err != nil {
			log.Println("error inserting into memstore:", err.Error())
			return
		}
	})

	return nil
}

// TODO: call the app to spawn if not openned yet
// TODO; perform logic to get user auth if user, for example, choose gdrive provider (for future)

func (s *downloaderService) download(ctx *fiber.Ctx) error {
	client := ctx.Params("client")

	id := ctx.Params("id")
	entry := s.memstore.Get(id)
	if entry == nil {
		return response.NotFound(ctx)
	}

	status := "Downloading"
	err := s.store.Update(entry.ID(), entryApi.UpdateDownload{
		Status: &status,
	})

	if err != nil {
		return response.InternalServerError(ctx, err)
	}

	go s.doDownload(entry, client)

	return response.Ok(ctx)
}

func (s *downloaderService) doDownload(entry entry.Entry, client string) {
	channel := api.CreateChannel(client)
	setting := setting.Get()

	dl := downloader.New(entry.Downloader(),
		downloader.UseSetting(setting),
	)

	if watcher, ok := dl.(downloader.Watcher); ok {
		watcher.Watch(func(data ...interface{}) {
			channel.Publish(data[0])
		})
	}

	if err := dl.Download(entry); err != nil {
		log.Printf("error downloading %s: %s", entry.Name(), err.Error())
		return
	}

	channel.Publish(rapidClient.Progress{
		ID:     entry.ID(),
		Lenght: entry.ChunkLen(),
		Done:   true,
	})
}

func (s *downloaderService) resume(ctx *fiber.Ctx) error {
	client := ctx.Params("client")

	id := ctx.Params("id")

	entry := s.memstore.Get(id)
	if entry == nil {
		return response.Success(ctx, fiber.StatusNoContent)
	}

	go s.doResume(entry, client)

	return response.Ok(ctx)
}

func (s *downloaderService) doResume(entry entry.Entry, client string) {
	channel := api.CreateChannel(client)
	setting := setting.Get()

	dl := downloader.New(entry.Downloader(),
		downloader.UseSetting(setting),
	)

	if watcher, ok := dl.(downloader.Watcher); ok {
		watcher.Watch(func(data ...interface{}) {
			channel.Publish(data[0])
		})
	}

	if err := dl.Resume(entry); err != nil {
		log.Printf("error downloading %s: %s", entry.Name(), err.Error())
		return
	}

	channel.Publish(rapidClient.Progress{
		ID:     entry.ID(),
		Lenght: entry.ChunkLen(),
		Done:   true,
	})
}

func (s *downloaderService) restart(ctx *fiber.Ctx) error {
	client := ctx.Params("client")

	id := ctx.Params("id")
	entry := s.memstore.Get(id)
	if entry == nil {
		return response.NotFound(ctx)
	}

	go s.doRestart(entry, client)

	return response.Ok(ctx)
}

func (s *downloaderService) doRestart(entry entry.Entry, client string) {
	channel := api.CreateChannel(client)
	setting := setting.Get()

	dl := downloader.New(entry.Downloader(),
		downloader.UseSetting(setting),
	)

	if watcher, ok := dl.(downloader.Watcher); ok {
		watcher.Watch(func(data ...interface{}) {
			channel.Publish(data[0])
		})
	}

	if err := dl.Restart(entry); err != nil {
		log.Printf("error restarting %s: %s", entry.Name(), err.Error())
		return
	}

	channel.Publish(rapidClient.Progress{
		ID:     entry.ID(),
		Lenght: entry.ChunkLen(),
		Done:   true,
	})
}

func (s *downloaderService) pause(ctx *fiber.Ctx) error {
	id := ctx.Params("id")
	entry := s.memstore.Get(id)
	if entry == nil {
		return response.Success(ctx, fiber.StatusNoContent)
	}

	return s.doStop(entry, ctx)
}

func (s *downloaderService) stop(ctx *fiber.Ctx) error {
	id := ctx.Params("id")
	entry := s.memstore.Get(id)
	if entry == nil {
		return response.NotFound(ctx)
	}

	s.memstore.Delete(entry.ID())

	return s.doStop(entry, ctx)
}

func (s *downloaderService) doStop(entry entry.Entry, ctx *fiber.Ctx) error {
	setting := setting.Get()

	dl := downloader.New(entry.Downloader(),
		downloader.UseSetting(setting),
	)

	if err := dl.Stop(entry); err != nil {
		return response.InternalServerError(ctx, fmt.Errorf("error stopping download: %s", err.Error()))
	}

	return response.Ok(ctx)
}

func (s *downloaderService) progressBar(c *websocket.Conn) {
	channel := api.CreateChannel(c.Params("client"))
	done := make(chan bool)

	ping := time.NewTicker(time.Second)

	go func() {
		for {
			t, _, err := c.ReadMessage()
			if err != nil {
				logger.Println("error reading message:", err)
				return
			}

			if t == websocket.CloseMessage {
				done <- true
			}
		}
	}()

	for {
		select {
		case <-done:
			return
		case data, ok := <-channel.Subscribe():
			if !ok {
				return
			}

			payload, err := json.Marshal(data)
			if err != nil {
				logger.Println("error marshalling data:", err)
				break
			}

			c.SetWriteDeadline(time.Now().Add(10 * time.Second))
			if err := c.WriteMessage(websocket.TextMessage, payload); err != nil {
				logger.Println("error sending progress data:", err)
				return
			}
		case <-ping.C:
			if err := c.WriteMessage(websocket.PingMessage, nil); err != nil {
				logger.Println("error ping:", err)
				return
			}
		}
	}
}

func (s *downloaderService) CreateRoutes() {
	s.app.Add("GET", "/:client/download/:id", s.download)
	s.app.Add("PUT", "/:client/restart/:id", s.restart)
	s.app.Add("PUT", "/:client/resume/:id", s.resume)
	s.app.Add("PUT", "/pause/:id", s.pause)
	s.app.Add("PUT", "/stop/:id", s.stop)
	s.app.Add("GET", "/ws/:client", websocket.New(s.progressBar))
}

func init() {
	api.RegisterService(newService)
}
