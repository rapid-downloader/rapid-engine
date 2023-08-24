package api

type (
	Channel interface {
		Publish(data interface{})
		Subscribe() <-chan interface{}
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

func (c *channel) Publish(data interface{}) {
	c.ch <- data
}

func (c *channel) Subscribe() <-chan interface{} {
	return c.ch
}

var channels = map[string]Channel{
	"gui": NewChannel(),
	"cli": NewChannel(),
}

func CreateChannel(name string) Channel {
	if channel, ok := channels[name]; ok {
		return channel
	}

	channel := NewChannel()
	channels[name] = channel

	return channel
}
