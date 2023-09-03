package downloader

import (
	"log"

	"github.com/rapid-downloader/rapid/entry"
	"github.com/rapid-downloader/rapid/logger"
	"github.com/rapid-downloader/rapid/setting"
)

type (
	Downloader interface {
		Download(entry entry.Entry) error
		Resume(entry entry.Entry) error
		Restart(entry entry.Entry) error
		Stop(entry entry.Entry) error
	}

	Watcher interface {
		Watch(update OnProgress)
	}

	DownloaderFactory func(o *option) Downloader

	OnProgress func(data ...interface{})

	option struct {
		setting *setting.Setting
		logger  logger.Logger
	}

	Options func(o *option)
)

func UseSetting(setting *setting.Setting) Options {
	return func(o *option) {
		o.setting = setting
	}
}

func UseLogger(l logger.Logger) Options {
	return func(o *option) {
		o.logger = l
	}

}

var downloadermap = make(map[string]DownloaderFactory)

func New(provider string, options ...Options) Downloader {
	opt := &option{
		setting: setting.Default(),
	}

	for _, option := range options {
		option(opt)
	}

	downloader, ok := downloadermap[provider]
	if !ok {
		log.Panicf("Provider %s is not implemented", provider)
		return nil
	}

	return downloader(opt)
}

func registerDownloader(name string, impl DownloaderFactory) {
	downloadermap[name] = impl
}
