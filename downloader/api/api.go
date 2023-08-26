package api

import (
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

func (s *downloaderService) downloadQueue(ctx *fiber.Ctx) error {
	channel := api.CreateChannel("queue")

	if s.entries.Queue.IsEmpty() {
		return response.Error(ctx, "Queue is empty")
	}

	go func() {
		for !s.entries.Queue.IsEmpty() {
			entry := s.entries.Queue.Pop()

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
		}
	}()

	return response.Success(ctx, nil)
}

func (s *downloaderService) download(ctx *fiber.Ctx) error {
	client := ctx.Params("client")
	channel := api.CreateChannel(client)

	id := ctx.Params("id")
	entry, ok := s.entries.List.Find(id)
	if !ok {
		return response.NotFound(ctx)
	}

	dl := downloader.New(entry.Downloader())
	if watcher, ok := dl.(downloader.Watcher); ok {
		watcher.Watch(func(data ...interface{}) {
			channel.Publish(data[0])
		})
	}

	go func() {
		if err := dl.Download(entry); err != nil {
			log.Printf("Error downloading %s: %s", entry.Name(), err.Error())
			return
		}

		s.entries.List.Remove(entry.ID())
		channel.Publish(downloader.Progressbar{
			ID:   id,
			Done: true,
		})
	}()

	return response.Success(ctx, nil)
}

// TODO: implement this
func (s *downloaderService) resume(ctx *fiber.Ctx) error {

	return response.Success(ctx, nil)
}

// TODO: implement this
func (s *downloaderService) pause(ctx *fiber.Ctx) error {

	return response.Success(ctx, nil)
}

// TODO: implement this
func (s *downloaderService) stop(ctx *fiber.Ctx) error {

	return response.Success(ctx, nil)
}

func (s *downloaderService) progressBar(c *websocket.Conn) {
	channel := api.CreateChannel(c.Params("client"))

	for data := range channel.Subscribe() {
		payload, err := json.Marshal(data)
		if err != nil {
			log.Println("Error marshalling data:", err)
			return
		}

		if err := c.WriteMessage(websocket.TextMessage, payload); err != nil {
			log.Println("Error sending progress data:", err)
			return
		}
	}
}

func (s *downloaderService) Router() []api.Route {
	return []api.Route{
		{
			Path:    "/:client/download/queue",
			Method:  "GET",
			Handler: s.downloadQueue,
		},
		{
			Path:    "/:client/download/:id",
			Method:  "GET",
			Handler: s.download,
		},
		{
			Path:    "/:client/resume/:id",
			Method:  "PUT",
			Handler: s.resume,
		},
		{
			Path:    "/:client/pause/:id",
			Method:  "PUT",
			Handler: s.pause,
		},
		{
			Path:    "/:client/stop/:id",
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
