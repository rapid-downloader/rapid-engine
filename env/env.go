package env

import (
	"os"

	"github.com/rapid-downloader/rapid/utils"
)

func Get(key string) utils.Parser {
	env := os.Getenv(key)
	return utils.Parse(env)
}
