package main

import (
	"log"

	logrus_udp "github.com/goofansu/go-tlog/hooks/udp"
	"github.com/goofansu/go-tlog/tlogrus"
)

// LoginEvent ...
type LoginEvent struct {
	UserID   int     `tlog:"userid"`
	Username *string `tlog:"username"`
}

func init() {
	hook, err := logrus_udp.NewHook("udp", "192.168.1.12:6667")
	if err != nil {
		log.Fatal("Unable to connect to UDP daemon")
	} else {
		tlogrus.AddHook(hook)
	}
}

func main() {
	username := "test"
	tlogrus.Log(&LoginEvent{1, &username})
}
