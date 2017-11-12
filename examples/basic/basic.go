package main

import (
	"fmt"
	"os"

	"github.com/goofansu/go-tlog/tlogrus"
)

// LoginEvent ...
type LoginEvent struct {
	UserID   int     `tlog:"userid"`
	Username *string `tlog:"username" validate:"required"`
}

func init() {
	tlogrus.SetOutput(os.Stdout)

	// f, err := os.OpenFile("tlog.log", os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0644)
	// if err != nil {
	//	log.Fatal("Cannot open log file to write", err)
	// }
	// tlogrus.SetOutput(f)
}

func main() {
	if err := tlogrus.Log(&LoginEvent{}); err != nil {
		fmt.Println(err)
	}
}
