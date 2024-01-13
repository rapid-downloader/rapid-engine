package main

import (
	"bytes"
	"context"
	"fmt"
	"net/http"

	"github.com/goccy/go-json"

	"github.com/rapid-downloader/rapid/client"
	"github.com/rapid-downloader/rapid/client/websocket"
	"github.com/rapid-downloader/rapid/env"
	"github.com/rapid-downloader/rapid/helper"
)

type rapidClient struct {
	id     string
	url    string
	wsUrl  string
	ws     websocket.Websocket
	ctx    context.Context
	cancel context.CancelFunc
}

func NewRapid(ctx context.Context, id string) (*rapidClient, error) {
	host := env.Get("API_HOST").String("localhost")
	port := env.Get("API_PORT").String(":9999")

	url := fmt.Sprintf("http://%s%s", host, port)
	wsUrl := fmt.Sprintf("ws://%s%s/ws/%s", host, port, id)

	ctx, cancel := context.WithCancel(ctx)
	ws := websocket.Connect(ctx, wsUrl)

	return &rapidClient{
		id:     id,
		url:    url,
		wsUrl:  wsUrl,
		ctx:    ctx,
		cancel: cancel,
		ws:     ws,
	}, nil

}

func (r *rapidClient) Listen(progressbar client.OnProgress) {
	r.ws.Listen(func(msg []byte) {
		var progress client.Progress
		if err := json.Unmarshal(msg, &progress); err != nil {
			progressbar(client.Progress{}, err)
			return
		}

		progressbar(progress, nil)
	})
}

func (r *rapidClient) Fetch(request client.Request) (*client.Download, error) {
	fetch := fmt.Sprintf("%s/fetch", r.url)

	payload, err := json.Marshal(request)
	if err != nil {
		return nil, fmt.Errorf("error marshalling request: %s", err)
	}

	req, err := http.NewRequestWithContext(r.ctx, "POST", fetch, bytes.NewBuffer(payload))
	if err != nil {
		return nil, fmt.Errorf("error preparing fetch request: %s", err.Error())
	}

	req.Header.Add("Content-Type", "application/json")

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error creating fetch request: %s", err)
	}

	defer res.Body.Close()
	var result client.Download
	// TODO: check if this is working or not
	if err := helper.UnmarshalBody(res, &result); err != nil {
		return nil, fmt.Errorf("error unmarshalling buffer: %s", err)
	}

	return &result, nil
}

func (r *rapidClient) Download(id string) error {
	download := fmt.Sprintf("%s/%s/download/%s", r.url, r.id, id)
	req, err := http.NewRequestWithContext(r.ctx, "GET", download, nil)
	if err != nil {
		return fmt.Errorf("error preparing download request: %s", err.Error())
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("error creating download request: %s", err)

	}

	defer res.Body.Close()

	return nil
}

func (r *rapidClient) Stop(id string) error {
	stop := fmt.Sprintf("%s/stop/%s", r.url, id)
	req, err := http.NewRequest("PUT", stop, nil)
	if err != nil {
		return fmt.Errorf("error preparing stop request: %s", err)
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("error stoping download: %s", err)
	}

	return res.Body.Close()
}

func (r *rapidClient) Close() error {
	r.cancel()
	return r.ws.Close()
}
