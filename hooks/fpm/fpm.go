// Package logrus_fpm represents File Per Message
package logrus_fpm

import (
	"fmt"
	"os"
	"strings"
	"sync"

	"github.com/sirupsen/logrus"
)

// Hook to save logs to seperate file by entry.Message
type Hook struct {
	mu *sync.Mutex
}

// NewHook ...
func NewHook() *Hook {
	hook := &Hook{
		mu: new(sync.Mutex),
	}
	return hook
}

// Fire ...
func (hook *Hook) Fire(entry *logrus.Entry) error {
	return hook.fileWrite(entry)
}

// Levels ...
func (hook *Hook) Levels() []logrus.Level {
	return logrus.AllLevels
}

// Write a log line directly to a file
func (hook *Hook) fileWrite(entry *logrus.Entry) error {
	var (
		fd  *os.File
		err error
	)

	hook.mu.Lock()
	defer hook.mu.Unlock()

	path := fmt.Sprintf("%s.log", entry.Message)
	fd, err = os.OpenFile(path, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0644)
	if err != nil {
		return err
	}
	defer fd.Close()

	line, err := entry.String()
	if err != nil {
		return err
	}

	line = strings.TrimPrefix(line, entry.Message+"|")
	_, err = fd.WriteString(line)
	if err != nil {
		return err
	}
	return nil
}
