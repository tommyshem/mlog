// Package mlog is simple logging module for go, with a rotating file feature
// and console logging.
package mlog

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path"
	"sync/atomic"
)

// LogLevel type
type LogLevel int32

const (
	// LevelTrace logs everything.
	LevelTrace LogLevel = (1 << iota)

	// LevelInfo logs Info, Warnings and Errors.
	LevelInfo

	// LevelWarn logs Warning and Errors.
	LevelWarn

	// LevelError logs just Errors.
	LevelError

	// LevelOnlyFile logs only to file.
	LevelOnlyFile
)

// defaultMaxBytes default log file max size limit set at 10 Mega bytes.
const defaultMaxBytes int = 10 * 1024 * 1024

// defaultBackupCount default log file max backup's made.
const defaultBackupCount int = 5

type mlog struct {
	LogLevel int32

	Trace   *log.Logger
	Info    *log.Logger
	Warning *log.Logger
	Error   *log.Logger
	Fatal   *log.Logger

	LogFile *rotatingFileHandler
}

// Logger instance of the mlog.
var Logger mlog

// DefaultFlags used by created loggers
var DefaultFlags = log.Ldate | log.Ltime | log.Lshortfile

//rotatingFileHandler writes log to file, if file size exceeds maxBytes,
//it will backup current file and open a new one.
//
//backupCount -> max number of backup file's created, it will delete the oldest if backup made is greater than backupCount.
type rotatingFileHandler struct {
	fd *os.File

	fileName    string
	maxBytes    int
	backupCount int
}

// newRotatingFileHandler creates directorys if needed and creates a log file if needed
// else opens the logfile if there is one already.
func newRotatingFileHandler(fileName string, maxBytes int, backupCount int) (*rotatingFileHandler, error) {
	dir := path.Dir(fileName)
	os.Mkdir(dir, 0777)

	h := new(rotatingFileHandler)

	if maxBytes <= 0 {
		return nil, fmt.Errorf("invalid max bytes")
	}

	h.fileName = fileName
	h.maxBytes = maxBytes
	h.backupCount = backupCount

	var err error
	h.fd, err = os.OpenFile(fileName, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		return nil, err
	}

	return h, nil
}

func (h *rotatingFileHandler) Write(p []byte) (n int, err error) {
	h.doRollover()
	return h.fd.Write(p)
}

// Close simply closes the File
func (h *rotatingFileHandler) Close() error {
	if h.fd != nil {
		return h.fd.Close()
	}
	return nil
}

func (h *rotatingFileHandler) doRollover() {
	f, err := h.fd.Stat()
	if err != nil {
		return
	}

	// log.Println("size: ", f.Size())

	if h.maxBytes <= 0 {
		return
	} else if f.Size() < int64(h.maxBytes) {
		return
	}

	if h.backupCount > 0 {
		h.fd.Close()

		for i := h.backupCount - 1; i > 0; i-- {
			sfn := fmt.Sprintf("%s.%d", h.fileName, i)
			dfn := fmt.Sprintf("%s.%d", h.fileName, i+1)

			os.Rename(sfn, dfn)
		}

		dfn := fmt.Sprintf("%s.1", h.fileName)
		os.Rename(h.fileName, dfn)

		h.fd, _ = os.OpenFile(h.fileName, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	}
}

// Start starts the logging
func Start(level LogLevel, path string) {
	doLogging(level, path, defaultMaxBytes, defaultBackupCount)
}

// StartEx starts the logging with maxBytes and backupCount changed from the default.
func StartEx(level LogLevel, path string, maxBytes, backupCount int) {
	doLogging(level, path, maxBytes, backupCount)
}

// Stop stops the logging
func Stop() error {
	if Logger.LogFile != nil {
		return Logger.LogFile.Close()
	}

	return nil
}

func Append(filename string, format string, a ...interface{}) {
	logInstance, _ := os.OpenFile(filename, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	defer logInstance.Close()
	_, err := logInstance.WriteString(fmt.Sprintf(format, a...))
	check(err)
}

func Appendln(filename string, format string, a ...interface{}) {
	logInstance, _ := os.OpenFile(filename, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	defer logInstance.Close()
	_, err := logInstance.WriteString(fmt.Sprintf(format+"\n", a...))
	check(err)
}

//Sync commits the current contents of the file to stable storage.
//Typically, this means flushing the file system's in-memory copy
//of recently written data to disk.
func Sync() {
	if Logger.LogFile != nil {
		Logger.LogFile.fd.Sync()
	}
}

func doLogging(logLevel LogLevel, fileName string, maxBytes, backupCount int) {
	traceHandle := ioutil.Discard
	infoHandle := ioutil.Discard
	warnHandle := ioutil.Discard
	errorHandle := ioutil.Discard
	fatalHandle := ioutil.Discard

	var fileHandle *rotatingFileHandler

	switch logLevel {
	case LevelOnlyFile:
		break
	case LevelTrace:
		traceHandle = os.Stdout
		fallthrough
	case LevelInfo:
		infoHandle = os.Stdout
		fallthrough
	case LevelWarn:
		warnHandle = os.Stdout
		fallthrough
	case LevelError:
		errorHandle = os.Stderr
		fatalHandle = os.Stderr
	}

	if fileName != "" {
		var err error
		fileHandle, err = newRotatingFileHandler(fileName, maxBytes, backupCount)
		if err != nil {
			log.Fatal("mlog: unable to create rotatingFileHandler: ", err)
		}

		if traceHandle == os.Stdout {
			traceHandle = io.MultiWriter(fileHandle, traceHandle)
		}

		if infoHandle == os.Stdout {
			infoHandle = io.MultiWriter(fileHandle, infoHandle)
		}

		if warnHandle == os.Stdout {
			warnHandle = io.MultiWriter(fileHandle, warnHandle)
		}

		if errorHandle == os.Stderr {
			errorHandle = io.MultiWriter(fileHandle, errorHandle)
		}

		if fatalHandle == os.Stderr {
			fatalHandle = io.MultiWriter(fileHandle, fatalHandle)
		}
	}

	Logger = mlog{
		Trace:   log.New(traceHandle, "T: ", DefaultFlags),
		Info:    log.New(infoHandle, "I: ", DefaultFlags),
		Warning: log.New(warnHandle, "W: ", DefaultFlags),
		Error:   log.New(errorHandle, "E: ", DefaultFlags),
		Fatal:   log.New(errorHandle, "F: ", DefaultFlags),
		LogFile: fileHandle,
	}

	atomic.StoreInt32(&Logger.LogLevel, int32(logLevel))
}

//** TRACE

// Trace writes to the Trace destination
func Trace(format string, a ...interface{}) {
	Logger.Trace.Output(2, "\033[34m"+fmt.Sprintf(format, a...)+"\033[0m")
}

//** INFO

// Info writes to the Info destination
func Info(format string, a ...interface{}) {
	Logger.Info.Output(2, "\033[32m"+fmt.Sprintf(format, a...)+"\033[0m")
}

//** WARNING

// Warning writes to the Warning destination
func Warning(format string, a ...interface{}) {
	Logger.Warning.Output(2, "\033[35m"+fmt.Sprintf(format, a...)+"\033[0m")
}

//** ERROR

// Error writes to the Error destination and accepts an err
func Error(err error) {
	Logger.Error.Output(2, "\033[33m"+fmt.Sprintf("%s\n", err)+"\033[0m")
}

// IfError is a shortcut function for log.Error if error
func IfError(err error) {
	if err != nil {
		Logger.Error.Output(2, fmt.Sprintf("%s\n", err))
	}
}

//** Fatal

// Fatalf writes to the Fatal destination and exits with an error code 255
func Fatalf(format string, a ...interface{}) {
	Logger.Fatal.Output(2, "\033[31m"+fmt.Sprintf(format, a...)+"\033[0m")
	Sync()
	os.Exit(255)
}

func check(err error) {
	if err != nil {
		log.Fatal("mlog: cannot append string to file: ", err)
	}
}

// FatalIfError is a shortcut function for log.Fatalf if error and
// exits with an error 255 code
func FatalIfError(err error) {
	if err != nil {
		Logger.Fatal.Output(2, fmt.Sprintf("%s\n", err))
		Sync()
		os.Exit(255)
	}
}
