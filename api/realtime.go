package api

import "github.com/gofiber/contrib/websocket"

type (
	Socket struct {
		Path    string
		Method  string
		Handler func(c *websocket.Conn)
	}

	Realtime interface {
		Sockets() []Socket
	}

	Channel interface {
		Post(name string, data interface{})
		RegisterSignal(name string, ch chan interface{})
		Signal(name string) <-chan interface{}
	}

	channel struct {
		mapCh map[string]chan interface{}
	}
)

var def = "default"
var chans = map[string]Channel{
	def: defaultChannel(),
}

func NewChannel(name ...string) Channel {
	n := def
	if len(name) > 0 {
		n = name[0]
	}

	if ch, ok := chans[n]; ok {
		return ch
	}

	return defaultChannel()
}

func (c *channel) RegisterSignal(name string, ch chan interface{}) {
	c.mapCh[name] = ch
}

func defaultChannel() Channel {
	return &channel{
		make(map[string]chan interface{}),
	}
}

func (c *channel) Post(path string, data interface{}) {
	c.mapCh[path] <- data
}

func (c *channel) Signal(name string) <-chan interface{} {
	return c.mapCh[name]
}
