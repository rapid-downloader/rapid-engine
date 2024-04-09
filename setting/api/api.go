package api

import (
	"os"
	"path/filepath"

	"github.com/BurntSushi/toml"
	"github.com/gofiber/fiber/v2"
	"github.com/rapid-downloader/rapid/api"
	response "github.com/rapid-downloader/rapid/helper"
	"github.com/rapid-downloader/rapid/setting"
)

type settingService struct {
	app *fiber.App
}

func newService(app *fiber.App) api.Service {
	return &settingService{
		app: app,
	}
}

func (s *settingService) getSetting(ctx *fiber.Ctx) error {
	setting := setting.Get()
	return response.Ok(ctx, setting)
}

func (s *settingService) updateSetting(ctx *fiber.Ctx) error {
	ss := setting.Get()

	var stg setting.Setting
	if err := ctx.BodyParser(&stg); err != nil {
		return response.BadRequest(ctx)
	}

	// TODO: validation

	location := filepath.Join(ss.DataLocation, "setting.toml")
	f, err := os.OpenFile(location, os.O_RDWR, os.ModePerm)
	if err != nil {
		return response.NotFound(ctx)
	}

	defer f.Close()
	if err := toml.NewEncoder(f).Encode(stg); err != nil {
		return response.InternalServerError(ctx, err)
	}

	return response.Ok(ctx, stg)
}

func (s *settingService) CreateRoutes() {
	s.app.Add("GET", "/settings", s.getSetting)
	s.app.Add("PUT", "/settings", s.updateSetting)
}

func init() {
	api.RegisterService(newService)
}
