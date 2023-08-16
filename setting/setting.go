package setting

import (
	"os"
	"path/filepath"
)

type (
	Setting struct {
		DownloadLocation string `json:"downloadLocation"`
		DataLocation     string `json:"dataLocation"`
		MaxRetry         int    `json:"maxRetry"`
		LoggerProvider   string `json:"loggerProvider"`
		MinChunkSize     int64  `json:"minChunkSize"`
	}
)

func Default() *Setting {
	home, _ := os.UserHomeDir()

	// location
	data := filepath.Join(home, ".rapid")
	download := filepath.Join(home, "Downloads")

	os.MkdirAll(data, os.ModePerm)

	return &Setting{
		DownloadLocation: download,
		DataLocation:     data,
		MaxRetry:         3,
		LoggerProvider:   "stdout",
		MinChunkSize:     1024 * 1024 * 5, // 5 MB
	}
}
