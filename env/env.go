package env

import (
	"os"

	"github.com/joho/godotenv"
)

func init() {
	godotenv.Load()
}

func Get(key string) Parser {
	env := os.Getenv(key)
	return parse(env)
}
