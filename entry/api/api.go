package entry

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/rapid-downloader/rapid/api"
	"github.com/rapid-downloader/rapid/downloader"
	"github.com/rapid-downloader/rapid/entry"
	response "github.com/rapid-downloader/rapid/helper"
)

type entryService struct {
	channel api.Channel
	entries map[string]entry.Entry
}

type Channel string

const (
	GUIChannel Channel = "gui"
	CLIChannel Channel = "cli"
)

func NewService() api.Service {
	return &entryService{}
}

func (s *entryService) Init() error {
	s.channel = api.NewChannel()
	s.entries = make(map[string]entry.Entry)

	return nil
}

func (s *entryService) Close() error {
	return nil
}

func (s *entryService) set(id string, entry entry.Entry) {
	s.entries[id] = entry
}

func (s *entryService) get(id string) entry.Entry {
	return s.entries[id]
}

func (s *entryService) delete(id string) {
	delete(s.entries, id)
}

func (s *entryService) createRequest(ctx *fiber.Ctx) error {
	var request browserRequest

	if err := ctx.BodyParser(&request); err != nil {
		return response.Error(ctx, err.Error())
	}

	entry, err := entry.Fetch(request.Url, request.toOptions()...)
	if err != nil {
		return response.Error(ctx, err.Error())
	}

	s.set(entry.ID(), entry)

	// TODO: call the app to spawn if not openned yet
	// TODO; perform logic to get user auth if user, for example, choose gdrive provider (for future)
	provider := ctx.Params("provider", downloader.Default)
	dl := downloader.New(provider)
	if watcher, ok := dl.(downloader.Watcher); ok {
		watcher.Watch(func(data ...interface{}) {
			s.channel.Post(string(GUIChannel), data)
		})
	}

	go func() {
		defer s.delete(entry.ID())

		if err := dl.Download(entry); err != nil {
			log.Printf("Error downloading %s:%s", entry.Name(), err.Error())
		}
	}()

	return response.Created(ctx)
}

// TODO: resume, stop, restart

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

func (s *entryService) Router() []api.Route {
	return []api.Route{
		{
			Path:    "/browser/:provider",
			Method:  "POST",
			Handler: s.createRequest,
		},
		{
			Path:    "/cli",
			Method:  "POST",
			Handler: s.cliRequest,
		},
	}
}

func init() {
	api.RegisterService(NewService())
}
