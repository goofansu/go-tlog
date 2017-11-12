package logrus_udp

import (
	"fmt"
	"net"
	"os"

	"github.com/sirupsen/logrus"
)

// Hook to send logs via net.Conn
type Hook struct {
	Writer  net.Conn
	Network string
	Address string
}

// NewHook ...
func NewHook(network, address string) (*Hook, error) {
	conn, err := net.Dial(network, address)
	return &Hook{conn, network, address}, err
}

// Fire ...
func (hook *Hook) Fire(entry *logrus.Entry) error {
	line, err := entry.String()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to read entry, %v", err)
		return err
	}

	_, err = hook.Writer.Write([]byte(line))
	if err != nil {
		return err
	}
	return nil
}

// Levels ...
func (hook *Hook) Levels() []logrus.Level {
	return logrus.AllLevels
}
