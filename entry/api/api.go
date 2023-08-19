package entry

import (
	"github.com/gofiber/fiber/v2"
	"github.com/rapid-downloader/rapid/api"
	"github.com/rapid-downloader/rapid/entry"
	response "github.com/rapid-downloader/rapid/helper"
)

type entryService struct {
}

func (s *entryService) Close() error {
	return nil
}

func (s *entryService) browserRequest(ctx *fiber.Ctx) error {
	var request browserRequest

	if err := ctx.BodyParser(&request); err != nil {
		return response.Error(ctx, err.Error())
	}

	_, err := entry.Fetch(request.Url, request.toOptions()...)
	if err != nil {
		return response.Error(ctx, err.Error())
	}

	// TODO: call the app to spawn if not openned yet

	return response.Created(ctx)
}

func (s *entryService) cliRequest(ctx *fiber.Ctx) error {
	var request cliRequest

	if err := ctx.BodyParser(&request); err != nil {
		return response.Error(ctx, err.Error())
	}

	entry, err := entry.Fetch(request.Url)
	if err != nil {
		return response.Error(ctx, err.Error())
	}

	return response.Success(ctx, entry)
}

func (s *entryService) hello(ctx *fiber.Ctx) error {
	return ctx.SendString("hello world")
}

func (s *entryService) Router() []api.Route {
	return []api.Route{
		{
			Path:    "/browser",
			Method:  "POST",
			Handler: s.browserRequest,
		},
		{
			Path:    "/cli",
			Method:  "POST",
			Handler: s.cliRequest,
		},
		{
			Path:    "/",
			Method:  "GET",
			Handler: s.hello,
		},
	}
}

func init() {
	api.RegisterService(&entryService{})
}
