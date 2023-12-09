package api

import (
	"fmt"
	logger "log"
	"time"

	"github.com/rapid-downloader/rapid/log"

	"github.com/goccy/go-json"

	"github.com/gofiber/contrib/websocket"
	"github.com/gofiber/fiber/v2"
	"github.com/rapid-downloader/rapid/api"
	"github.com/rapid-downloader/rapid/db"
	"github.com/rapid-downloader/rapid/downloader"
	"github.com/rapid-downloader/rapid/entry"
	entryApi "github.com/rapid-downloader/rapid/entry/api"
	response "github.com/rapid-downloader/rapid/helper"
	"github.com/rapid-downloader/rapid/setting"
)

type downloaderService struct {
	memstore entry.Store
	channel  api.Channel
	store    entryApi.Store
}

func newService() api.Service {
	return &downloaderService{
		memstore: entry.Memstore(),
		channel:  api.CreateChannel("memstore"),
		store:    entryApi.NewStore("download", db.DB()),
	}
}

func (s *downloaderService) Init() error {
	go s.channel.Subscribe(func(data interface{}) {
		entry, ok := data.(entry.Entry)
		if !ok {
			return
		}

		if err := s.memstore.Set(entry.ID(), entry); err != nil {
			log.Println("error inserting into memstore:", err.Error())
			return
		}
	})

	return nil
}

// TODO: call the app to spawn if not openned yet
// TODO; perform logic to get user auth if user, for example, choose gdrive provider (for future)

func (s *downloaderService) download(ctx *fiber.Ctx) error {
	client := ctx.Params("client")

	id := ctx.Params("id")
	entry := s.memstore.Get(id)
	if entry == nil {
		return response.NotFound(ctx)
	}

	status := "Downloading"
	err := s.store.Update(entry.ID(), entryApi.UpdateDownload{
		Status: &status,
	})

	if err != nil {
		return response.Error(ctx, err.Error(), fiber.StatusBadGateway)
	}

	go s.doDownload(entry, client)

	return response.Success(ctx, nil)
}

func (s *downloaderService) doDownload(entry entry.Entry, client string) {
	channel := api.CreateChannel(client)
	setting := setting.Get()

	dl := downloader.New(entry.Downloader(),
		downloader.UseSetting(setting),
	)

	if watcher, ok := dl.(downloader.Watcher); ok {
		watcher.Watch(func(data ...interface{}) {
			channel.Publish(map[string]interface{}{
				"id":         data[0],
				"index":      data[1],
				"downloaded": data[2],
				"size":       data[3],
				"progress":   data[4],
				"done":       false,
			})
		})
	}

	if err := dl.Download(entry); err != nil {
		log.Printf("error downloading %s: %s", entry.Name(), err.Error())
		return
	}

	s.memstore.Delete(entry.ID())

	channel.Publish(map[string]interface{}{
		"id":         entry.ID(),
		"index":      0,
		"downloaded": entry.Size(),
		"size":       entry.Size(),
		"progress":   100,
		"done":       true,
	})

	status := "Completed"
	err := s.store.Update(entry.ID(), entryApi.UpdateDownload{
		Status: &status,
	})
	if err != nil {
		log.Println(err.Error())
	}
}

func (s *downloaderService) resume(ctx *fiber.Ctx) error {
	client := ctx.Params("client")

	id := ctx.Params("id")
	entry := s.memstore.Get(id)
	if entry == nil {
		return response.NotFound(ctx)
	}

	status := "Downloading"
	err := s.store.Update(entry.ID(), entryApi.UpdateDownload{
		Status: &status,
	})

	if err != nil {
		return response.Error(ctx, err.Error(), fiber.StatusBadGateway)
	}

	go s.doResume(entry, client)

	return response.Success(ctx, nil)
}

func (s *downloaderService) doResume(entry entry.Entry, client string) {
	channel := api.CreateChannel(client)
	setting := setting.Get()

	dl := downloader.New(entry.Downloader(),
		downloader.UseSetting(setting),
	)

	if watcher, ok := dl.(downloader.Watcher); ok {
		watcher.Watch(func(data ...interface{}) {
			channel.Publish(map[string]interface{}{
				"id":         data[0],
				"index":      data[1],
				"downloaded": data[2],
				"size":       data[3],
				"progress":   data[4],
				"done":       false,
			})
		})
	}

	if err := dl.Resume(entry); err != nil {
		log.Printf("error downloading %s: %s", entry.Name(), err.Error())
		return
	}

	s.memstore.Delete(entry.ID())
	channel.Publish(map[string]interface{}{
		"id":         entry.ID(),
		"index":      0,
		"downloaded": entry.Size(),
		"size":       entry.Size(),
		"progress":   100,
		"done":       true,
	})

	status := "Completed"
	err := s.store.Update(entry.ID(), entryApi.UpdateDownload{
		Status: &status,
	})
	if err != nil {
		log.Println(err.Error())
	}
}

func (s *downloaderService) restart(ctx *fiber.Ctx) error {
	client := ctx.Params("client")

	id := ctx.Params("id")
	entry := s.memstore.Get(id)
	if entry == nil {
		return response.NotFound(ctx)
	}

	status := "Downloading"
	err := s.store.Update(entry.ID(), entryApi.UpdateDownload{
		Status: &status,
	})

	if err != nil {
		return response.Error(ctx, err.Error(), fiber.StatusBadGateway)
	}

	go s.doRestart(entry, client)

	return response.Success(ctx, nil)
}

func (s *downloaderService) doRestart(entry entry.Entry, client string) {
	channel := api.CreateChannel(client)
	setting := setting.Get()

	dl := downloader.New(entry.Downloader(),
		downloader.UseSetting(setting),
	)

	if watcher, ok := dl.(downloader.Watcher); ok {
		watcher.Watch(func(data ...interface{}) {
			channel.Publish(map[string]interface{}{
				"id":         data[0],
				"index":      data[1],
				"downloaded": data[2],
				"size":       data[3],
				"progress":   data[4],
				"done":       false,
			})
		})
	}

	defer s.memstore.Delete(entry.ID())

	if err := dl.Restart(entry); err != nil {
		log.Printf("error restarting %s: %s", entry.Name(), err.Error())
		return
	}

	channel.Publish(map[string]interface{}{
		"id":         entry.ID(),
		"index":      0,
		"downloaded": entry.Size(),
		"size":       entry.Size(),
		"progress":   100,
		"done":       true,
	})

	status := "Completed"
	err := s.store.Update(entry.ID(), entryApi.UpdateDownload{
		Status: &status,
	})
	if err != nil {
		log.Println(err.Error())
	}
}

func (s *downloaderService) pause(ctx *fiber.Ctx) error {
	id := ctx.Params("id")
	entry := s.memstore.Get(id)
	if entry == nil {
		return response.NotFound(ctx)
	}

	status := "Paused"
	err := s.store.Update(entry.ID(), entryApi.UpdateDownload{
		Status: &status,
	})

	if err != nil {
		return response.Error(ctx, err.Error(), fiber.StatusBadGateway)
	}

	return s.doStop(entry, ctx)
}

func (s *downloaderService) stop(ctx *fiber.Ctx) error {
	id := ctx.Params("id")
	entry := s.memstore.Get(id)
	if entry == nil {
		return response.NotFound(ctx)
	}

	s.memstore.Delete(entry.ID())

	status := "Stoped"
	err := s.store.Update(entry.ID(), entryApi.UpdateDownload{
		Status: &status,
	})

	if err != nil {
		return response.Error(ctx, err.Error(), fiber.StatusBadGateway)
	}

	return s.doStop(entry, ctx)
}

func (s *downloaderService) doStop(entry entry.Entry, ctx *fiber.Ctx) error {
	setting := setting.Get()

	dl := downloader.New(entry.Downloader(),
		downloader.UseSetting(setting),
	)

	if err := dl.Stop(entry); err != nil {
		return response.Error(ctx, fmt.Sprint("error stopping download:", err.Error()))
	}

	return response.Success(ctx, nil)
}

func (s *downloaderService) progressBar(c *websocket.Conn) {
	channel := api.CreateChannel(c.Params("client"))
	done := make(chan bool)

	ping := time.NewTicker(time.Second)

	go func() {
		for {
			t, _, err := c.ReadMessage()
			if err != nil {
				logger.Println("error reading message:", err)
				return
			}

			if t == websocket.CloseMessage {
				done <- true
			}
		}
	}()

	for {
		select {
		case <-done:
			return
		case data, ok := <-channel.Subscribe():
			if !ok {
				return
			}

			payload, err := json.Marshal(data)
			if err != nil {
				logger.Println("error marshalling data:", err)
				break
			}

			c.SetWriteDeadline(time.Now().Add(10 * time.Second))
			if err := c.WriteMessage(websocket.TextMessage, payload); err != nil {
				logger.Println("error sending progress data:", err)
				return
			}
		case <-ping.C:
			if err := c.WriteMessage(websocket.PingMessage, nil); err != nil {
				logger.Println("error ping:", err)
				continue
			}
		}
	}
}

func (s *downloaderService) Routes() []api.Route {
	return []api.Route{
		{
			Path:    "/:client/download/:id",
			Method:  "GET",
			Handler: s.download,
		},
		{
			Path:    "/:client/restart/:id",
			Method:  "PUT",
			Handler: s.restart,
		},
		{
			Path:    "/:client/resume/:id",
			Method:  "PUT",
			Handler: s.resume,
		},
		{
			Path:    "/pause/:id",
			Method:  "PUT",
			Handler: s.pause,
		},
		{
			Path:    "/stop/:id",
			Method:  "PUT",
			Handler: s.stop,
		},
		{
			Path:    "/ws/:client",
			Method:  "GET",
			Handler: websocket.New(s.progressBar),
		},
	}
}

func init() {
	api.RegisterService(newService)
}
