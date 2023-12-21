package env

import (
	"os"
)

func Get(key string) Parser {
	env := os.Getenv(key)
	return parse(env)
}
