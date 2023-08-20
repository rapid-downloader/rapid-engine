package api

import "github.com/gofiber/contrib/websocket"

type (
	Socket struct {
		Path    string
		Method  string
		Handler func(c *websocket.Conn)
	}

	WebSocket interface {
		Sockets() []Socket
	}

	Channel interface {
		Post(data interface{})
		Signal() <-chan interface{}
	}

	channel struct {
		ch chan interface{}
	}
)

func NewChannel() Channel {
	return &channel{
		make(chan interface{}, 100),
	}
}

func (c *channel) Post(data interface{}) {
	c.ch <- data
}

func (c *channel) Signal() <-chan interface{} {
	return c.ch
}

var channels = map[string]Channel{
	"gui": NewChannel(),
	"cli": NewChannel(),
}

func CreateChannel(client string) Channel {
	return channels[client]
}
