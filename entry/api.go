package entry

import (
	"github.com/gofiber/fiber/v2"
	response "github.com/rapid-downloader/rapid/helper"
)

func BrowserRequest(ctx *fiber.Ctx) error {
	var request browserRequest

	if err := ctx.BodyParser(&request); err != nil {
		return response.Error(ctx, err.Error())
	}

	_, err := Fetch(request.Url, request.toOptions()...)
	if err != nil {
		return response.Error(ctx, err.Error())
	}

	// TODO: call the app to spawn if not openned yet

	return response.Created(ctx)
}
