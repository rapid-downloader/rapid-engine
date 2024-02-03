package api

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/rapid-downloader/rapid/api"
	response "github.com/rapid-downloader/rapid/helper"
	"github.com/rapid-downloader/rapid/setting"
)

type logService struct {
	app *fiber.App
}

func newService(app *fiber.App) api.Service {
	return &logService{
		app: app,
	}
}

func today() string {
	now := time.Now()
	d := now.Day()
	m := now.Month()
	y := now.Year()

	return fmt.Sprintf("%d-%d-%d", d, int(m), y)
}

func (s *logService) logs(ctx *fiber.Ctx) error {
	date := ctx.Params("date", today())
	setting := setting.Get()

	path := filepath.Join(setting.DataLocation, "logs", fmt.Sprintf("%s.txt", date))
	file, err := os.Open(path)
	if err != nil {
		return response.NotFound(ctx)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)

	logs := make([]string, 0)
	for scanner.Scan() {
		logs = append(logs, scanner.Text())
	}

	return response.Ok(ctx, logs)
}

func (s *logService) CreateRoutes() {
	s.app.Add("GET", "/logs/:date", s.logs)
}

func init() {
	api.RegisterService(newService)
}
