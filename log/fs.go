package log

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/rapid-downloader/rapid/setting"
)

type fsLogger struct {
	*sync.Mutex
	path string
}

const FS = "fs"

// stdoutLogger will log into std out
func FSLogger(s *setting.Setting) Logger {
	const DDMMYYYY = "02-01-2006"

	dir := fmt.Sprintf("%s/logs", s.DataLocation)
	os.MkdirAll(dir, os.ModePerm)

	filename := filepath.Join(dir, time.Now().Format(DDMMYYYY)+".txt")

	return &fsLogger{
		Mutex: &sync.Mutex{},
		path:  filename,
	}
}

func prefix() string {
	const FORMAT = "01-02-2006 15:04:05"
	timestamp := time.Now().Format(FORMAT)

	return timestamp + " "
}

func formatMessage(args ...interface{}) string {
	return prefix() + fmt.Sprintln(args...)
}

func (l *fsLogger) println(args ...interface{}) string {
	l.Lock()
	defer l.Unlock()

	file, err := os.OpenFile(l.path, os.O_CREATE|os.O_APPEND|os.O_RDWR, 0644)
	if err != nil {
		log.Fatal("Error creating or opening file log:", err.Error())
	}

	defer file.Close()

	msg := formatMessage(args...)
	if _, err := file.WriteString(msg); err != nil {
		log.Println("Error writing into log file:", err.Error())
	}

	return msg
}

func (l *fsLogger) Println(args ...interface{}) {
	l.println(args...)
}

func (l *fsLogger) Panicln(args ...interface{}) {
	msg := l.println(args...)
	fmt.Print(msg)
	os.Exit(1)
}

func init() {
	registerLogger(FS, FSLogger)
}
