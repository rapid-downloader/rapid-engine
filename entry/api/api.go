package entry

import (
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
	s.ws = api.NewChannel()
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
	client := ctx.Params("client")
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

	// s.set(entry.ID(), entry)

	// TODO: call the app to spawn if not openned yet
	// TODO; perform logic to get user auth if user, for example, choose gdrive provider (for future)

	s.ws.Post(client, entry)

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

	// s.ws.Post(ClientCLI, entry)

	return response.Success(ctx, entry)
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
