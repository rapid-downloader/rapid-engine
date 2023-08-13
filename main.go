package main

import (
	"fmt"
	"log"
	"time"

	"github.com/rapid-downloader/rapid/downloader"
	"github.com/rapid-downloader/rapid/entry"
)

func main() {
	url := "https://link.testfile.org/PDF50MB"
	entry, err := entry.Fetch(url)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(entry)

	dl := downloader.New(downloader.Default)
	dl.(downloader.Watcher).Watch(func(data ...interface{}) {
		log.Println(data...)
	})

	go func() {
		if err := dl.Download(entry); err != nil {
			log.Fatal(err)
		}
	}()

	time.Sleep(10 * time.Second)

	if err := dl.Stop(entry); err != nil {
		log.Fatal(err)
	}

	if err := dl.Resume(entry); err != nil {
		log.Fatal(err)
	}
}
