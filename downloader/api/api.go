package api

import (
	"log"

	"github.com/goccy/go-json"
	"github.com/gofiber/contrib/websocket"
	"github.com/rapid-downloader/rapid/api"
	"github.com/rapid-downloader/rapid/downloader"
	"github.com/rapid-downloader/rapid/entry"
)

type downloaderWebSocket struct {
	channel api.Channel
}

func NewWebsocket() api.WebSocket {
	return &downloaderWebSocket{}
}

func (s *downloaderWebSocket) Init() error {
	s.channel = api.NewChannel()

	return nil
}

// TODO: test this
func (s *downloaderWebSocket) progressBar(c *websocket.Conn) {
	client := c.Params("client")

	for {
		select {
		case data := <-s.channel.Signal(client):
			entry := data.(entry.Entry)

			dl := downloader.New(entry.DownloadProvider())

			if watcher, ok := dl.(downloader.Watcher); ok {
				watcher.Watch(func(data ...interface{}) {
					payload, err := json.Marshal(data[0])
					if err != nil {
						log.Println("Error marshalling data:", err)
						return
					}

					if err := c.WriteMessage(websocket.TextMessage, payload); err != nil {
						log.Println("Error sending payload:", err)
						return
					}
				})
			}

			go func() {
				if err := dl.Download(entry); err != nil {
					log.Println("Error downloading entry:", err)
					return
				}
			}()

		default:
			_, message, err := c.ReadMessage()
			if err != nil {
				log.Println("Error reading message:", err)
			}

			log.Println(string(message))
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
