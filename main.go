package main

import (
	"log"

	"github.com/goccy/go-json"
	"github.com/gofiber/contrib/websocket"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/rapid-downloader/rapid/api"
	_ "github.com/rapid-downloader/rapid/downloader/api"
	_ "github.com/rapid-downloader/rapid/entry/api"
)

func main() {
	app := fiber.New(fiber.Config{
		JSONEncoder: json.Marshal,
		JSONDecoder: json.Unmarshal,
	})

	app.Use(logger.New())
	app.Use("/ws", func(c *fiber.Ctx) error {
		if websocket.IsWebSocketUpgrade(c) {
			c.Locals("allowed", true)
			return c.Next()
		}

		return fiber.ErrUpgradeRequired
	})

	api := api.Create(app)

	api.Run()
	api.Shutdown()

	log.Fatal(app.Listen(":3333"))
}
