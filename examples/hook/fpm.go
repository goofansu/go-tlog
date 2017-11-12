package main

import (
	logrus_fpm "github.com/goofansu/go-tlog/hooks/fpm"
	"github.com/goofansu/go-tlog/tlogrus"
)

// LoginEvent ...
type LoginEvent struct {
	UserID   int     `tlog:"userid"`
	Username *string `tlog:"username"`
}

// RegEvent ...
type RegEvent struct {
	Username *string `tlog:"username"`
}

func init() {
	tlogrus.AddHook(logrus_fpm.NewHook())
}

func main() {
	username := "test"
	tlogrus.Log(&LoginEvent{1, &username})
	tlogrus.Log(&RegEvent{&username})
}
