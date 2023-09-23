package setting

import (
	"os"
	"path/filepath"

	"github.com/BurntSushi/toml"
)

type (
	Setting struct {
		DownloadLocation string `toml:"download_location"`
		DataLocation     string `toml:"data_location"`
		MaxRetry         int    `toml:"max_retry"`
		MinChunkSize     int64  `toml:"min_chunk_size"`
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
		MinChunkSize:     1024 * 1024 * 5, // 5 MB
	}
}

func Get() *Setting {
	s := Default()
	location := filepath.Join(s.DataLocation, "setting.toml")

	file, err := os.Open(location)
	if err != nil {
		f, _ := os.Create(location)
		toml.NewEncoder(f).Encode(s)

		return s
	}

	defer file.Close()

	var setting Setting
	decoder := toml.NewDecoder(file)
	if _, err := decoder.Decode(&setting); err != nil {
		return s
	}

	return &setting
}
