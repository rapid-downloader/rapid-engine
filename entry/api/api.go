package entry

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/rapid-downloader/rapid/api"
	"github.com/rapid-downloader/rapid/downloader"
	"github.com/rapid-downloader/rapid/entry"
	response "github.com/rapid-downloader/rapid/helper"
)

const (
	ClientCLI = "cli"
	ClientGUI = "gui"
)

type entryService struct {
	ws      api.Channel
	entries map[string]entry.Entry
}

func NewService() api.Service {
	return &entryService{}
}

func (s *entryService) Init() error {
	s.entries = make(map[string]entry.Entry)

	return nil
}

func (s *entryService) Close() error {
	return nil
}

// TODO: use better approach than this
func (s *entryService) set(id string, entry entry.Entry) {
	s.entries[id] = entry
}

func (s *entryService) get(id string) entry.Entry {
	return s.entries[id]
}

func (s *entryService) delete(id string) {
	delete(s.entries, id)
}

func (s *entryService) createRequestWithOpen(ctx *fiber.Ctx) error {
	// client := ctx.Params("client")
	provider := ctx.Params("provider", downloader.Default)

	var req request

	if err := ctx.BodyParser(&req); err != nil {
		return response.Error(ctx, err.Error())
	}

	options := req.toOptions()
	options = append(options, entry.UseDownloadProvider(provider))

	entry, err := entry.Fetch(req.Url, options...)
	if err != nil {
		return response.Error(ctx, err.Error())
	}

	// TODO: find a way to resume and cancel for certain entry based on their pointer
	// TODO: call the app to spawn if not openned yet
	// TODO; perform logic to get user auth if user, for example, choose gdrive provider (for future)

	go s.download(entry)

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

	go s.download(entry)

	return response.Success(ctx, entry)
}

func (s *entryService) download(entry entry.Entry) error {
	dl := downloader.New(entry.DownloadProvider())
	channel := api.CreateChannel(ClientCLI)

	if watcher, ok := dl.(downloader.Watcher); ok {
		watcher.Watch(func(data ...interface{}) {
			channel.Post(data[0])
		})
	}

	if err := dl.Download(entry); err != nil {
		log.Printf("Error downloading %s: %v", entry.Name(), err)
		return err
	}

	channel.Post(map[string]bool{
		"done": true,
	})

	return nil
}

func (s *entryService) Router() []api.Route {
	return []api.Route{
		{
			Path:    "/download/:client/:provider",
			Method:  "POST",
			Handler: s.createRequestWithOpen,
		},
		{
			Path:    "/download/cli/",
			Method:  "POST",
			Handler: s.cliRequest,
		},
	}
}

func init() {
	api.RegisterService(NewService())
}
