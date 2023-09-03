package logger

import (
	"log"
	"sync"

	"github.com/rapid-downloader/rapid/setting"
)

type (
	Logger interface {
		Print(...interface{})
	}

	LogCloser interface {
		Close() error
	}

	LoggerFactory func(*setting.Setting) Logger
)

var loggermap = make(map[string]LoggerFactory)
var instance sync.Map

// TODO: rethink how we should properly close the logger if we are going to provide file base log

func New(provider string, s *setting.Setting) Logger {
	val, ok := instance.Load(provider)
	if ok {
		return val.(Logger)
	}

	logger, ok := loggermap[provider]
	if !ok {
		log.Panicf("Provider %s is not implemented", provider)
		return nil
	}

	l := logger(s)
	instance.Store(provider, l)

	return l
}

func registerLogger(name string, impl LoggerFactory) {
	loggermap[name] = impl
}
