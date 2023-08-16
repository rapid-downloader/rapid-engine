package main

import (
	"github.com/gofiber/fiber/v2"
	"github.com/rapid-downloader/rapid/entry"
)

func Router(app *fiber.App) {
	app.Post("/browser", entry.BrowserRequest)
}
