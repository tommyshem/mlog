package main

import (
	"fmt"
	"mlog"
)

func main() {
	mlog.Start(mlog.LevelTrace, "app.log")

	mlog.Trace("Hello world")
	fmt.Println("Welcome to SSS")
	mlog.Info("Hello world")
	mlog.Warning("Hello world")
}
