package main

import "github.com/rapid-downloader/rapid/client"

var globalStore = make(map[string]client.Download)

func store(id string, entry client.Download) {
	globalStore[id] = entry
}

func loadStored() (client.Download, bool) {
	for _, entry := range globalStore {
		return entry, true
	}

	return client.Download{}, false
}
