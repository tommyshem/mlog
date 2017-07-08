### This is a fork from https://github.com/TranDuyThanh/mlog and changed to my needs.

A simple logging module for go, with a rotating file feature and optional console logging.

### Installation
	
	go get github.com/tommyshem/mlog

### Setup logging
You can set the logging level on a Logger, then it will only log entries with that severity or anything above it:
	
```go
	mlog.Start(mlog.LevelTrace,   "/tmp/app.log")   // Logs Trace, Warning, Info, Error and Fetal - (All logs) 
	mlog.Start(mlog.LevelInfo,    "/tmp/app.log")   // Logs Info, Warning, Error and Fetal 
	mlog.Start(mlog.LevelWarning, "/tmp/app.log")   // Logs Warning, Error and Fetal
	mlog.Start(mlog.LevelError,   "/tmp/app.log")   // Logs Error and Fatal only
	mlog.Start(mlog.LevelOnlyFile, "/tmp/app.log")  // Logs all but only to file
```
### Level logging

```go
	mlog.Trace("Hello World !")
	mlog.Info("Hello World !")
	mlog.Warning("Hello World !")
	mlog.Error("Hello World !")
	mlog.Fatal("Hello World !")
```

### Rotation

Log rotation is provided with `mlog`. Default max size of log file is: `10 MB`

Note:
	
	Log rotation should be done by an external program (like logrotate(8)) that can compress and delete old log entries

### Example

# Write to stdout/stderr and create a rotating log file

```go
package main

import (
	"github.com/TranDuyThanh/mlog"
)

func main() {
	mlog.Start(mlog.LevelInfo, "/tmp/app.log")

	mlog.Info("Hello World !")

	ipsum := "ipsum"
	mlog.Warning("Lorem %s", ipsum)
}
```


# Write to stdout/stderr only

```go
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
```

By default, the log will be rolled over to a backup file when its size reaches 10Mb and 10 such files will be created (and eventually reused).

Alternatively, you can specify the max size of the log file before it gets rotated, 
and the number of backup files you want to create, with the StartEx function.

```go
package main

import "github.com/TranDuyThanh/mlog"

func main() {
    mlog.StartEx(mlog.LevelInfo, "app.log", 5*1024*1024, 5)

    mlog.Info("Hello World !")

    ipsum := "ipsum"
    mlog.Warning("Lorem %s", ipsum)
}
```
This will rotate the file when it reaches 5Mb and 5 backup files will eventually be created.

Setting logger flags:
```go
mlog.DefaultFlags = log.Ldate|log.Ltime|log.Lmicroseconds|log.Lshortfile
```

## Output

```bash
I: 2015/05/15 07:09:45 main.go:10: Hello World !
W: 2015/05/15 07:09:45 main.go:13: Lorem ipsum
```
