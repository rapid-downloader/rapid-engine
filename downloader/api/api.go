package api

import (
	"log"

	"github.com/goccy/go-json"
	"github.com/gofiber/contrib/websocket"
	"github.com/rapid-downloader/rapid/api"
)

type downloaderService struct {
	channel api.Channel
}

func NewService() api.Realtime {
	return &downloaderService{}
}

func (s *downloaderService) Init() error {
	s.channel = api.NewChannel()

	return nil
}

// TODO: test this
func (s *downloaderService) progressBar(c *websocket.Conn) {
	client := c.Params("client")

	for data := range s.channel.Signal(client) {
		payload, err := json.Marshal(data)
		if err != nil {
			log.Println("Error marshalling data:", err)
			return
		}

		if err := c.WriteMessage(websocket.TextMessage, payload); err != nil {
			log.Println("Error sending payload:", err)
			return
		}
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
	api.RegisterRealtime(NewService())
}
