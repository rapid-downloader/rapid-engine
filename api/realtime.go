package api

import "github.com/gofiber/contrib/websocket"

type (
	Channel struct {
		Path    string
		Method  string
		Handler func(c *websocket.Conn)
	}

	Realtime interface {
		Channels() []Channel
	}
)
