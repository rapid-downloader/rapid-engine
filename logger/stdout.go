package logger

import (
	"log"

	"github.com/rapid-downloader/rapid/setting"
)

type stdLogger struct{}

const StdOut = "stdout"

// stdoutLogger will log into std out
func stdoutLogger(_ *setting.Setting) Logger {
	return &stdLogger{}
}

func (l *stdLogger) Print(args ...interface{}) {
	log.Println(args...)
}

func (l *stdLogger) Panic(args ...interface{}) {
	log.Panic(args...)
}

func init() {
	registerLogger(StdOut, stdoutLogger)
}
