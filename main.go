package main

import (
	"log"

	"github.com/goccy/go-json"
	"github.com/gofiber/contrib/websocket"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/rapid-downloader/rapid/api"
	"github.com/rapid-downloader/rapid/db"
	_ "github.com/rapid-downloader/rapid/downloader/api"
	_ "github.com/rapid-downloader/rapid/entry/api"
	_ "github.com/rapid-downloader/rapid/log/api"
)

func main() {
	db.Open()
	defer db.Close()

	app := fiber.New(fiber.Config{
		JSONEncoder: json.Marshal,
		JSONDecoder: json.Unmarshal,
	})

	app.Use(cors.New())
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

	if err := app.Listen(":9999"); err != nil {
		log.Fatal(err)
	}
}
