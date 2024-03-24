package api

import "sync"

type (
	OnPublished func(data interface{})

	Channel interface {
		Publish(data interface{})
		Subscribe(callback ...OnPublished) <-chan interface{}
		Close() error
	}

	channel struct {
		ch   chan interface{}
		once sync.Once
		name string
	}
)

func NewChannel(name ...string) Channel {
	n := ""
	if len(name) > 0 {
		n = name[0]
	}

	return &channel{
		make(chan interface{}, 100),
		sync.Once{},
		n,
	}
}

func (c *channel) Publish(data interface{}) {
	c.ch <- data
}

func (c *channel) Subscribe(callback ...OnPublished) <-chan interface{} {
	if len(callback) > 0 {
		for data := range c.ch {
			callback[0](data)
		}

		return nil
	}

	return c.ch
}

func (c *channel) Close() error {
	c.once.Do(func() {
		close(c.ch)

		if c.name != "" {
			delete(channels, c.name)
		}
	})

	return nil
}

var channels = map[string]Channel{
	"gui": NewChannel(),
	"cli": NewChannel(),
}

func CreateChannel(name string) Channel {
	if channel, ok := channels[name]; ok {
		return channel
	}

	channel := NewChannel(name)
	channels[name] = channel

	return channel
}
