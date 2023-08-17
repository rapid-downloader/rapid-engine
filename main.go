package main

import (
	"github.com/goccy/go-json"
	"github.com/gofiber/contrib/websocket"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/rapid-downloader/rapid/api"
	_ "github.com/rapid-downloader/rapid/downloader"
	_ "github.com/rapid-downloader/rapid/entry"
)

func main() {
	app := fiber.New(fiber.Config{
		JSONEncoder: json.Marshal,
		JSONDecoder: json.Unmarshal,
	})

	app.Use(logger.New())
	app.Use("/ws", func(c *fiber.Ctx) error {
		// IsWebSocketUpgrade returns true if the client
		// requested upgrade to the WebSocket protocol.
		if websocket.IsWebSocketUpgrade(c) {
			c.Locals("allowed", true)
			return c.Next()
		}

		return fiber.ErrUpgradeRequired
	})

	for _, service := range api.Services() {
		if init, ok := service.(api.ServiceInitter); ok {
			if err := init.Init(); err != nil {
				log.Debug(err)
			}
		}

		for _, route := range service.Router() {
			app.Add(route.Method, route.Path, route.Handler)
		}

		if ch, ok := service.(api.Realtime); ok {
			for _, channel := range ch.Channels() {
				app.Add(channel.Method, channel.Path, websocket.New(channel.Handler))
			}
		}
	}

	app.Listen(":3000")
	// TODO: gracefull shutdown
}
