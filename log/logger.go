package log

import (
	"fmt"
	"log"
	"sync"

	"github.com/rapid-downloader/rapid/setting"
)

type (
	Logger interface {
		Println(...interface{})
		Panicln(...interface{})
	}

	LogCloser interface {
		Close() error
	}

	LoggerFactory func(*setting.Setting) Logger
)

var loggermap = make(map[string]LoggerFactory)
var instance sync.Map

var logging Logger

func New(provider string) Logger {
	setting := setting.Get()

	val, ok := instance.Load(provider)
	if ok {
		return val.(Logger)
	}

	logger, ok := loggermap[provider]
	if !ok {
		log.Panicf("provider %s is not implemented", provider)
		return nil
	}

	l := logger(setting)
	instance.Store(provider, l)

	return l
}

func Println(args ...interface{}) {
	logging.Println(args...)
}

func Printf(format string, args ...interface{}) {
	str := fmt.Sprintf(format, args...)
	logging.Println(str)
}

func Panicln(args ...interface{}) {
	logging.Panicln(args...)
}

func registerLogger(name string, impl LoggerFactory) {
	loggermap[name] = impl
}

var once sync.Once

func init() {
	once.Do(func() {
		logging = New(FS)
	})
}
