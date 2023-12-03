package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/rapid-downloader/rapid/client"
	_ "github.com/rapid-downloader/rapid/env"
	"github.com/rapid-downloader/rapid/helper"
	"github.com/vbauerster/mpb"
	"github.com/vbauerster/mpb/decor"
)

type progressBar struct {
	mpb    *mpb.Progress
	barMap sync.Map
}

func progressbar() progressBar {
	return progressBar{
		mpb:    mpb.New(),
		barMap: sync.Map{},
	}
}

func (p *progressBar) update(index int, downloaded int64, chunkSize int64) {
	i := fmt.Sprintf("%d", index)

	if val, ok := p.barMap.Load(i); ok {
		bar := val.(*mpb.Bar)
		bar.IncrBy(int(downloaded - bar.Current()))

		return
	}

	bar := p.mpb.AddBar(chunkSize,
		mpb.PrependDecorators(
			decor.CountersKiloByte("% .2f / % .2f"),
		),
		mpb.AppendDecorators(
			decor.AverageETA(decor.ET_STYLE_MMSS),
			decor.Name(" | "),
			decor.AverageSpeed(decor.UnitKB, "% .2f"),
		),
	)

	p.barMap.Store(i, bar)
}

func main() {
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, []os.Signal{syscall.SIGINT, syscall.SIGKILL, syscall.SIGTERM, syscall.SIGSTOP, os.Interrupt}...)
	ctx, cancel := context.WithCancel(context.Background())

	rapid, err := client.New(ctx, helper.ID(5))
	if err != nil {
		log.Fatal(err)
	}

	executeCommands(ctx, rapid)
	progressBar := progressbar()

	go rapid.Listen(func(progress client.Progress, err error) {
		if err != nil {
			if err.Error() == "websocket: close sent" {
				cancel()
				return
			}

			fmt.Println(err)
			cancel()
			return
		}

		if progress.Done {
			cancel()
			return
		}

		progressBar.update(progress.Index, progress.Downloaded, progress.Size)
	})

	for {
		select {
		case <-ctx.Done():
			return
		case <-interrupt:
			stop(rapid)
			return
		}
	}
}

func stop(rapid client.Rapid) {
	entry, ok := loadStored()
	if !ok {
		return
	}

	rapid.Stop(entry.ID)
}
