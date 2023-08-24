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

	LoggerFactory func(*setting.Setting) Logger
)

var loggermap = make(map[string]LoggerFactory)
var instance sync.Map

// TODO: rethink how we should properly close the logger if we are going to provide file base log

func New(setting *setting.Setting) Logger {
	val, ok := instance.Load(setting.LoggerProvider)
	if ok {
		return val.(Logger)
	}

	logger, ok := loggermap[setting.LoggerProvider]
	if !ok {
		log.Panicf("Provider %s is not implemented", setting.LoggerProvider)
		return nil
	}

	l := logger(setting)
	instance.Store(setting.LoggerProvider, l)

	return l
}

func registerLogger(name string, impl LoggerFactory) {
	loggermap[name] = impl
}
