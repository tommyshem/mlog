package main

import (
	"fmt"

	"github.com/tommyshem/mlog"
)

func main() {
	var testValue = 10

	// Set logging level and create logging file.
	mlog.Start(mlog.LevelTrace, "/tmp/app.log") // logs all levels

	fmt.Println("App Running") // Pretend app code under test.

	// Standard loging functions.
	mlog.Trace("App has started running")
	mlog.Info("This is a info message")
	mlog.Warning("Text not passed")

	// Any Fetal logs will log the message and then stop the program with error code 255
	mlog.Fatal("Fetal error testValue = ", testValue) // comment out to get to next logs.
	// mlog.Fatalf("This is a Fetal message ", testValue) // uncomment to log fetal and Fetal will stop the program with error code 255.

	fmt.Println("App still running") // Pretend app code under test.

	// Misc Functions
	// Use Append and Appendln to add text to any file passed in to function.
	mlog.Append("/tmp/test.log", "%s - %s", "Hello", "World")
	mlog.Appendln("/tmp/test.log", "%s - %s", "Hello", "World 2")
	fmt.Println("App Stopped")
}
