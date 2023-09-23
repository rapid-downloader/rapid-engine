package api

import (
	"fmt"
	"log"
	"time"

	"github.com/goccy/go-json"

	"github.com/gofiber/contrib/websocket"
	"github.com/gofiber/fiber/v2"
	"github.com/rapid-downloader/rapid/api"
	"github.com/rapid-downloader/rapid/downloader"
	"github.com/rapid-downloader/rapid/entry"
	response "github.com/rapid-downloader/rapid/helper"
	"github.com/rapid-downloader/rapid/setting"
)

type downloaderService struct {
	memstore entry.Store
}

func newService(memstore entry.Store) api.Service {
	return &downloaderService{
		memstore: memstore,
	}
}

func (s *downloaderService) Init() error {
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

	go s.doDownload(entry, client)

	return response.Success(ctx, nil)
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
		log.Printf("Error downloading %s: %s", entry.Name(), err.Error())
		return
	}

	s.memstore.Delete(entry.ID())

	channel.Publish(downloader.Progressbar{
		ID:   entry.ID(),
		Done: true,
	})
}

func (s *downloaderService) resume(ctx *fiber.Ctx) error {
	client := ctx.Params("client")

	id := ctx.Params("id")
	entry := s.memstore.Get(id)
	if entry == nil {
		return response.NotFound(ctx)
	}

	go s.doResume(entry, client)

	return response.Success(ctx, nil)
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
		log.Printf("Error downloading %s: %s", entry.Name(), err.Error())
		return
	}

	s.memstore.Delete(entry.ID())
	channel.Publish(downloader.Progressbar{
		ID:   entry.ID(),
		Done: true,
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

	return response.Success(ctx, nil)
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

	defer s.memstore.Delete(entry.ID())

	if err := dl.Restart(entry); err != nil {
		log.Printf("Error downloading %s: %s", entry.Name(), err.Error())
		return
	}

	channel.Publish(downloader.Progressbar{
		ID:   entry.ID(),
		Done: true,
	})
}

func (s *downloaderService) pause(ctx *fiber.Ctx) error {
	id := ctx.Params("id")
	entry := s.memstore.Get(id)
	if entry == nil {
		return response.NotFound(ctx)
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
		return response.Error(ctx, fmt.Sprint("Error stopping download:", err.Error()))
	}

	return response.Success(ctx, nil)
}

func (s *downloaderService) progressBar(c *websocket.Conn) {
	channel := api.CreateChannel(c.Params("client"))
	done := make(chan bool)

	ping := time.NewTicker(time.Second)

	go func() {
		for {
			t, _, err := c.ReadMessage()
			if err != nil {
				log.Println("Error reading message:", err)
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

			progressBar := data.(downloader.Progressbar)

			payload, err := json.Marshal(progressBar)
			if err != nil {
				log.Println("Error marshalling data:", err)
				break
			}

			c.SetWriteDeadline(time.Now().Add(10 * time.Second))
			if err := c.WriteMessage(websocket.TextMessage, payload); err != nil {
				log.Println("Error sending progress data:", err)
				return
			}
		case <-ping.C:
			if err := c.WriteMessage(websocket.PingMessage, nil); err != nil {
				log.Println("Error ping:", err)
				return
			}
		}
	}
}

func (s *downloaderService) Router() []api.Route {
	return []api.Route{
		{
			Path:    "/:client/download/:id",
			Method:  "GET",
			Handler: s.download,
		},
		{
			Path:    "/:client/restart/:id",
			Method:  "PUT",
			Handler: s.restart,
		},
		{
			Path:    "/:client/resume/:id",
			Method:  "PUT",
			Handler: s.resume,
		},
		{
			Path:    "/pause/:id",
			Method:  "PUT",
			Handler: s.pause,
		},
		{
			Path:    "/stop/:id",
			Method:  "PUT",
			Handler: s.stop,
		},
		{
			Path:    "/ws/:client",
			Method:  "GET",
			Handler: websocket.New(s.progressBar),
		},
	}
}

func init() {
	api.RegisterService(newService)
}
