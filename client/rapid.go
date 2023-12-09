package client

import (
	"bytes"
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/goccy/go-json"

	"github.com/gorilla/websocket"
	"github.com/rapid-downloader/rapid/env"
	"github.com/rapid-downloader/rapid/helper"
)

type rapidClient struct {
	id     string
	url    string
	wsUrl  string
	ws     *websocket.Conn
	ctx    context.Context
	cancel context.CancelFunc
}

func New(ctx context.Context, id string) (Rapid, error) {

	host := env.Get("API_HOST").String("localhost")
	port := env.Get("API_PORT").String(":9999")

	url := fmt.Sprintf("http://%s%s", host, port)
	wsUrl := fmt.Sprintf("ws://%s%s/ws/%s", host, port, id)

	ws, res, err := websocket.DefaultDialer.DialContext(ctx, wsUrl, nil)
	if err != nil {
		return nil, fmt.Errorf("error dialing websocket: %v. Status code %d", err, res.StatusCode)
	}

	ctx, cancel := context.WithCancel(ctx)
	return &rapidClient{
		id:     id,
		url:    url,
		wsUrl:  wsUrl,
		ctx:    ctx,
		cancel: cancel,
		ws:     ws,
	}, nil
}

func (r *rapidClient) Listen(progressbar OnProgress) {
	onError := make(chan bool)

	go func() {
		for {
			select {
			case <-onError:
				return
			default:
				_, message, err := r.ws.ReadMessage()
				if err != nil {
					progressbar(Progress{}, err)
					return
				}

				var progress Progress
				if err := json.Unmarshal(message, &progress); err != nil {
					fmt.Println("Error unmarshalling message:", err)
					return
				}

				progressbar(progress, nil)
			}
		}
	}()

	ping := time.NewTicker(time.Second)

	go func() {
		for {
			select {
			case <-ping.C:
				if err := r.ws.WriteMessage(websocket.PingMessage, nil); err != nil {
					onError <- true
					progressbar(Progress{}, err)
					return
				}
			case <-r.ctx.Done():
				return
			}
		}
	}()

	for {
		select {
		case <-onError:
			r.closeConn()
			return
		case <-r.ctx.Done():
			r.closeConn()
			return
		}
	}
}

func (r *rapidClient) closeConn() {
	if err := r.ws.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, "")); err != nil {
		fmt.Println("Error sending close signal to server:", err)
		return
	}
}

func (r *rapidClient) Download(request Request) (*Download, error) {
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

	var result Download
	// TODO: check if this is working or not
	if err := helper.UnmarshalBody(res, &result); err != nil {
		return nil, fmt.Errorf("error unmarshalling buffer: %s", err)
	}

	download := fmt.Sprintf("%s/%s/download/%s", r.url, r.id, result.ID)
	req, err = http.NewRequestWithContext(r.ctx, "GET", download, nil)
	if err != nil {
		return nil, fmt.Errorf("error preparing download request: %s", err.Error())
	}

	res, err = http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error creating download request: %s", err)

	}

	res.Body.Close()

	return &result, nil
}

func (r *rapidClient) Resume(id string) error {
	return r.docontinue("resume", id)
}

func (r *rapidClient) Restart(id string) error {
	return r.docontinue("restart", id)
}

func (r *rapidClient) docontinue(t, id string) error {
	resume := fmt.Sprintf("%s/%s/%s/%s", r.url, r.id, t, id)

	req, err := http.NewRequest("PUT", resume, nil)
	if err != nil {
		return fmt.Errorf("error preparing %s request: %s", t, err)
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("error %sing download: %s", t, err)
	}

	return res.Body.Close()
}

func (r *rapidClient) Stop(id string) error {
	return r.stop("stop", id)
}

func (r *rapidClient) Pause(id string) error {
	return r.stop("pause", id)
}

func (r *rapidClient) stop(t string, id string) error {
	stop := fmt.Sprintf("%s/%s/%s", r.url, t, id)

	req, err := http.NewRequest("PUT", stop, nil)
	if err != nil {
		return fmt.Errorf("error preparing %s request: %s", t, err)
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("error %sing download: %s", t, err)
	}

	return res.Body.Close()
}

func (r *rapidClient) Close() error {
	r.cancel()
	return nil
}
