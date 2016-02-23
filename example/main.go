package main

import (
	"fmt"
	"github.com/TranDuyThanh/mlog"
)

func main() {
	mlog.Start(mlog.LevelTrace, "app.log")

	mlog.Trace("Hello world")
	fmt.Println("Welcome to SSS")
	mlog.Info("Hello world")
	mlog.Warning("Hello world")
	mlog.Append("test.log", "%s - %s", "Hello", "World")
	mlog.Appendln("test.log", "%s - %s", "Hello", "World 2")
}
