package api

import (
	"log"

	"github.com/goccy/go-json"

	"github.com/gofiber/contrib/websocket"
	"github.com/rapid-downloader/rapid/api"
)

type downloaderWebSocket struct{}

func NewWebsocket() api.WebSocket {
	return &downloaderWebSocket{}
}

func (s *downloaderWebSocket) Init() error {

	return nil
}

func (s *downloaderWebSocket) progressBar(c *websocket.Conn) {
	channel := api.CreateChannel(c.Params("client"))

	for data := range channel.Signal() {
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

func (s *downloaderWebSocket) Sockets() []api.Socket {
	return []api.Socket{
		{
			Path:    "/ws/:client",
			Method:  "GET",
			Handler: s.progressBar,
		},
	}
}

func init() {
	api.RegisterWebSocket(NewWebsocket())
}
