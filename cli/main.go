package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/joho/godotenv"
	"github.com/rapid-downloader/rapid/client"
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

var once = sync.Once{}

func (p *progressBar) update(progress client.Progress) {
	once.Do(func() {
		for i, chunk := range progress.Chunks {
			bar := p.mpb.AddBar(chunk.Size,
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
	})

	for i, chunk := range progress.Chunks {
		if val, ok := p.barMap.Load(i); ok {
			bar := val.(*mpb.Bar)

			if chunk.Done {
				bar.SetTotal(chunk.Size, true)
			} else {
				bar.IncrBy(int(chunk.Downloaded - bar.Current()))
			}
		}

	}
}

func init() {
	godotenv.Load("../.env")
}

func main() {
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, []os.Signal{syscall.SIGINT, syscall.SIGKILL, syscall.SIGTERM, syscall.SIGSTOP, os.Interrupt}...)
	ctx, cancel := context.WithCancel(context.Background())

	rapid, err := NewRapid(ctx, helper.Id(5))
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

		progressBar.update(progress)
	})

	for {
		select {
		case <-ctx.Done():
			rapid.Close()
			return
		case <-interrupt:
			stop(rapid)
			return
		}
	}
}

func stop(rapid *rapidClient) {
	entry, ok := loadStored()
	if !ok {
		return
	}

	rapid.Stop(entry.ID)
	rapid.Close()
}
