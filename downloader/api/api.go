package api

import (
	"log"

	"github.com/goccy/go-json"

	"github.com/gofiber/contrib/websocket"
	"github.com/gofiber/fiber/v2"
	"github.com/rapid-downloader/rapid/api"
	"github.com/rapid-downloader/rapid/entry"
	response "github.com/rapid-downloader/rapid/helper"
)

type downloaderService struct {
	entries *entry.Listing
}

func NewWebsocket(entries *entry.Listing) api.Service {
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

// TODO: implement this
func (s *downloaderService) download(ctx *fiber.Ctx) error {

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
			Path:    "/:client/download/:id",
			Method:  "GET",
			Handler: s.download,
		},
		{
			Path:    "/:client/resume/:id",
			Method:  "UPDATE",
			Handler: s.resume,
		},
		{
			Path:    "/:client/pause/:id",
			Method:  "UPDATE",
			Handler: s.pause,
		},
		{
			Path:    "/:client/stop/:id",
			Method:  "UPDATE",
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
	api.RegisterService(NewWebsocket)
}
