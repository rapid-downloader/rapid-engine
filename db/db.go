package db

import (
	"fmt"
	"time"

	"github.com/rapid-downloader/rapid/setting"
	bolt "go.etcd.io/bbolt"
)

var k = "default"
var instances = make(map[string]*bolt.DB)

func DB(key ...string) *bolt.DB {
	instance := k
	if len(key) > 0 {
		instance = key[0]
	}

	return instances[instance]
}

func Open(key ...string) {
	setting := setting.Get()

	path := fmt.Sprintf("%s/%s", setting.DataLocation, "entries.db")
	options := bolt.DefaultOptions
	options.Timeout = time.Duration(time.Second)

	db, err := bolt.Open(path, 0600, options)
	if err != nil {
		panic(err)
	}

	instance := k
	if len(key) > 0 {
		instance = key[0]
	}

	instances[instance] = db
}

func Close(key ...string) error {
	return DB(key...).Close()
}
