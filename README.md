A simple logging module for go, with a rotating file feature and console logging.

### Installation
	
	go get github.com/TranDuyThanh/mlog

### Level logging

	mlog.Trace("Hello World !")
	mlog.Info("Hello World !")
	mlog.Warning("Hello World !")
	mlog.Error("Hello World !")
	mlog.Fatal("Hello World !")

You can set the logging level on a Logger, then it will only log entries with that severity or anything above it:
	
	mlog.Start(mlog.LevelTrace,   "app.log")
	mlog.Start(mlog.LevelInfo,    "app.log")
	mlog.Start(mlog.LevelWarning, "app.log")
	mlog.Start(mlog.LevelError,   "app.log")
	
### Rotation

Log rotation is provided with `mlog`. Default max size of log file is: `1GB`

Note:
	
	Log rotation should be done by an external program (like logrotate(8)) that can compress and delete old log entries

### Example

Write to stdout/stderr and create a rotating logfile

	package main

	import (
		"github.com/TranDuyThanh/mlog"
	)

	func main() {
		mlog.Start(mlog.LevelInfo, "app.log")

		mlog.Info("Hello World !")

		ipsum := "ipsum"
		mlog.Warning("Lorem %s", ipsum)
	}

Write to stdout/stderr only

	package main

	import (
		"github.com/TranDuyThanh/mlog"
	)

	func main() {
		mlog.Start(mlog.LevelInfo, "")

		mlog.Info("Hello World !")

		ipsum := "ipsum"
		mlog.Warning("Lorem %s", ipsum)
	}

## Output

```
I: 2015/10/13 09:14:40.212225 std_mlog.go:10: Hello World !
W: 2015/10/13 09:14:40.212614 std_mlog.go:13: Lorem ipsum
```
