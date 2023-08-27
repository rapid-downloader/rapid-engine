package api

import (
	"fmt"
	"log"

	"github.com/goccy/go-json"

	"github.com/gofiber/contrib/websocket"
	"github.com/gofiber/fiber/v2"
	"github.com/rapid-downloader/rapid/api"
	"github.com/rapid-downloader/rapid/downloader"
	"github.com/rapid-downloader/rapid/entry"
	response "github.com/rapid-downloader/rapid/helper"
)

type downloaderService struct {
	entries *entry.Listing
}

func newService(entries *entry.Listing) api.Service {
	return &downloaderService{
		entries: entries,
	}
}

func (s *downloaderService) Init() error {
	return nil
}

// TODO: find a way to resume and cancel for certain entry based on their pointer
// TODO: call the app to spawn if not openned yet
// TODO; perform logic to get user auth if user, for example, choose gdrive provider (for future)

func (s *downloaderService) download(ctx *fiber.Ctx) error {
	client := ctx.Params("client")

	id := ctx.Params("id")
	entry, ok := s.entries.List.Find(id)
	if !ok {
		return response.NotFound(ctx)
	}

	go s.doDownload(entry, client)

	return response.Success(ctx, nil)
}

func (s *downloaderService) doDownload(entry entry.Entry, client string) {
	channel := api.CreateChannel(client)

	dl := downloader.New(entry.Downloader())
	if watcher, ok := dl.(downloader.Watcher); ok {
		watcher.Watch(func(data ...interface{}) {
			channel.Publish(data[0])
		})
	}

	if err := dl.Download(entry); err != nil {
		log.Printf("Error downloading %s: %s", entry.Name(), err.Error())
		return
	}

	s.entries.List.Remove(entry.ID())

	channel.Publish(downloader.Progressbar{
		ID:   entry.ID(),
		Done: true,
	})
}

func (s *downloaderService) resume(ctx *fiber.Ctx) error {
	client := ctx.Params("client")

	id := ctx.Params("id")
	entry, ok := s.entries.List.Find(id)
	if !ok {
		return response.NotFound(ctx)
	}

	go s.doResume(entry, client)

	return response.Success(ctx, nil)
}

func (s *downloaderService) doResume(entry entry.Entry, client string) {
	channel := api.CreateChannel(client)

	dl := downloader.New(entry.Downloader())
	if watcher, ok := dl.(downloader.Watcher); ok {
		watcher.Watch(func(data ...interface{}) {
			channel.Publish(data[0])
		})
	}

	if err := dl.Resume(entry); err != nil {
		log.Printf("Error downloading %s: %s", entry.Name(), err.Error())
		return
	}

	s.entries.List.Remove(entry.ID())
	channel.Publish(downloader.Progressbar{
		ID:   entry.ID(),
		Done: true,
	})
}

func (s *downloaderService) restart(ctx *fiber.Ctx) error {
	client := ctx.Params("client")

	id := ctx.Params("id")
	entry, ok := s.entries.List.Find(id)
	if !ok {
		return response.NotFound(ctx)
	}

	go s.doRestart(entry, client)

	return response.Success(ctx, nil)
}

func (s *downloaderService) doRestart(entry entry.Entry, client string) {
	channel := api.CreateChannel(client)

	dl := downloader.New(entry.Downloader())
	if watcher, ok := dl.(downloader.Watcher); ok {
		watcher.Watch(func(data ...interface{}) {
			channel.Publish(data[0])
		})
	}

	defer s.entries.List.Remove(entry.ID())

	if err := dl.Restart(entry); err != nil {
		log.Printf("Error downloading %s: %s", entry.Name(), err.Error())
		return
	}

	channel.Publish(downloader.Progressbar{
		ID:   entry.ID(),
		Done: true,
	})
}

func (s *downloaderService) stop(ctx *fiber.Ctx) error {
	id := ctx.Params("id")
	entry, ok := s.entries.List.Find(id)
	if !ok {
		return response.NotFound(ctx)
	}

	defer s.entries.List.Remove(entry.ID())

	dl := downloader.New(entry.Downloader())

	if err := dl.Stop(entry); err != nil {
		return response.Error(ctx, fmt.Sprint("Error stopping download:", err.Error()))
	}

	return response.Success(ctx, nil)
}

func (s *downloaderService) progressBar(c *websocket.Conn) {
	channel := api.CreateChannel(c.Params("client"))

	defer c.Close()

	for data := range channel.Subscribe() {
		progressBar := data.(downloader.Progressbar)

		payload, err := json.Marshal(progressBar)
		if err != nil {
			log.Println("Error marshalling data:", err)
			break
		}

		if err := c.WriteMessage(websocket.TextMessage, payload); err != nil {
			log.Println("Error sending progress data:", err)
			break
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
			Path:    "/stop/:id",
			Method:  "PUT",
			Handler: s.stop,
		},
	}
}

func (s *downloaderService) Sockets() []api.Socket {
	return []api.Socket{
		{
			Path:    "/ws/:client",
			Method:  "GET",
			Handler: s.progressBar,
		},
	}
}

func init() {
	api.RegisterService(newService)
}
