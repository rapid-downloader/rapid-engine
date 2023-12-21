package websocket

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/goccy/go-json"
	"github.com/gorilla/websocket"
	"github.com/rapid-downloader/rapid/log"
)

type ListenFunc = func(msg []byte)

type Websocket interface {
	Listen(callback ListenFunc)
	Write(payload interface{}) error
	Close() error
}

type wsClient struct {
	url     string
	sendBuf chan []byte
	ctx     context.Context
	cancel  context.CancelFunc

	mutex sync.RWMutex
	conn  *websocket.Conn
}

func Connect(ctx context.Context, url string) Websocket {
	conn := &wsClient{
		url: url,
	}

	conn.ctx, conn.cancel = context.WithCancel(ctx)

	go conn.ping()
	go conn.listenWrite()
	return conn
}

func (ws *wsClient) connect() *websocket.Conn {
	ws.mutex.Lock()
	defer ws.mutex.Unlock()

	if ws.conn != nil {
		return ws.conn
	}

	conn, _, err := websocket.DefaultDialer.Dial(ws.url, nil)
	if err != nil {
		return nil
	}

	ws.conn = conn
	return conn
}

func (ws *wsClient) Listen(callback ListenFunc) {
	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ws.ctx.Done():
			return
		case <-ticker.C:
			for {
				conn := ws.connect()
				if conn == nil {
					break // break the inner loop to reconnect when the server is down
				}

				_, msg, err := conn.ReadMessage()
				if err != nil {
					log.Println("error reading websocket message:", err)
					ws.close()
					break
				}

				callback(msg)
			}
		}
	}
}

// listen for write so that it can handle write asyncronously
func (ws *wsClient) listenWrite() {
	for data := range ws.sendBuf {
		conn := ws.connect()
		if conn == nil {
			continue
		}

		if err := conn.WriteMessage(websocket.TextMessage, data); err != nil {
			log.Println("error sending data to websocket server:", err)
		}
	}
}

func (ws *wsClient) Write(payload interface{}) error {
	data, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	select {
	case ws.sendBuf <- data:
		return nil
	case <-ws.ctx.Done():
		return fmt.Errorf("context canceled")
	}
}

const pingPeriod = 10 * time.Second

func (ws *wsClient) ping() {
	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			conn := ws.connect()
			if conn == nil {
				continue
			}

			if err := conn.WriteControl(websocket.PingMessage, []byte{}, time.Now().Add(pingPeriod/2)); err != nil {
				ws.close()
			}

		case <-ws.ctx.Done():
			return
		}
	}
}

func (ws *wsClient) close() {
	ws.mutex.Lock()

	if ws.conn != nil {
		ws.conn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
		ws.conn.Close()
		ws.conn = nil
	}

	ws.mutex.Unlock()
}

func (ws *wsClient) Close() error {
	ws.cancel()
	ws.close()

	return nil
}
