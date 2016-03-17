package log

import (
	"encoding/json"
	"fmt"
	"os"
	"sync"

	"github.com/xvtom/logs"
)

const (
	DefaultBufferSize = 1024
)

type LogLevel int

const (
	LogLevelTrace LogLevel = iota
	LogLevelDebug
	LogLevelInfo
	LogLevelWarn
	LogLevelError
	LogLevelCritical
)

var (
	loggers = make(map[string]*logs.BeeLogger)
	mutex   sync.Mutex
)

var Levels = map[string]int{
	"Trace":    0,
	"Debug":    1,
	"Info":     2,
	"Warn":     3,
	"Error":    4,
	"Critical": 5,
}

// Trace log trace level message
func Trace(v ...interface{}) {
	Tracef("%s", fmt.Sprint(v...))
}

// Traceln log trace level message
func Traceln(v ...interface{}) {
	Tracef("%s", fmt.Sprintln(v...))
}

// Tracef log trace level message
func Tracef(format string, v ...interface{}) {
	for _, logger := range loggers {
		logger.Trace(format, v...)
	}
}

// Debug log debug level message
func Debug(v ...interface{}) {
	Debugf("%s", fmt.Sprint(v...))
}

// Debugln log debug level message
func Debugln(v ...interface{}) {
	Debugf("%s", fmt.Sprintln(v...))
}

// Debugf log debug level message
func Debugf(format string, v ...interface{}) {
	for _, logger := range loggers {
		logger.Debug(format, v...)
	}
}

// Info log info level message
func Info(v ...interface{}) {
	Infof("%s", fmt.Sprint(v...))
}

// Infoln log info level message
func Infoln(v ...interface{}) {
	Infof("%s", fmt.Sprintln(v...))
}

// Infof log info level message
func Infof(format string, v ...interface{}) {
	for _, logger := range loggers {
		logger.Info(format, v...)
	}
}

// Print log info level message
func Print(v ...interface{}) {
	Printf("%s", fmt.Sprint(v...))
}

// Println log info level message
func Println(v ...interface{}) {
	Printf("%s", fmt.Sprintln(v...))
}

// Printf log info level message
func Printf(format string, v ...interface{}) {
	for _, logger := range loggers {
		logger.Info(format, v...)
	}
}

// Warn log warn level message
func Warn(v ...interface{}) {
	Warnf("%s", fmt.Sprint(v...))
}

// Warnln log warn level message
func Warnln(v ...interface{}) {
	Warnf("%s", fmt.Sprintln(v...))
}

// Warnf log warn level message
func Warnf(format string, v ...interface{}) {
	for _, logger := range loggers {
		logger.Warn(format, v...)
	}
}

// Error log error level message
func Error(v ...interface{}) {
	Errorf("%s", fmt.Sprint(v...))
}

// Errorln log error level message
func Errorln(v ...interface{}) {
	Errorf("%s", fmt.Sprintln(v...))
}

// Errorf log error level message
func Errorf(format string, v ...interface{}) {
	if len(loggers) == 0 {
		fmt.Fprintf(os.Stderr, format, v...)
		return
	}
	for _, logger := range loggers {
		logger.Error(format, v...)
	}
}

// Critical log critical level message
func Critical(v ...interface{}) {
	Criticalf("%s", fmt.Sprint(v...))
}

// Criticalln log critical level message
func Criticalln(v ...interface{}) {
	Criticalf("%s", fmt.Sprintln(v...))
}

// Criticalf log critical level message
func Criticalf(format string, v ...interface{}) {
	if len(loggers) == 0 {
		fmt.Fprintf(os.Stderr, format, v...)
		return
	}
	for _, logger := range loggers {
		logger.Critical(format, v...)
	}
}

// Fatal is equivalent to Critical() followed by a call to os.Exit(1)
func Fatal(v ...interface{}) {
	Fatalf("%s", fmt.Sprint(v...))
}

// Fatalln is equivalent to Criticalln() followed by a call to os.Exit(1)
func Fatalln(v ...interface{}) {
	Fatalf("%s", fmt.Sprintln(v...))
}

// Fatalf is equivalent to Criticalf() followed by a call to os.Exit(1)
func Fatalf(format string, v ...interface{}) {
	Criticalf(format, v...)
	for _, l := range loggers {
		l.Close()
	}
	os.Exit(1)
}

// Close close all loggers
func Close() {
	for _, l := range loggers {
		l.Close()
	}
	loggers = make(map[string]*logs.BeeLogger)
}

// RegisterLogger register a logger to the global logger list
func RegisterLogger(mode string, bufSize int64, config interface{}) error {
	logger := logs.NewLogger(bufSize)
	logger.SetLogFuncCallDepth(3)

	c, err := json.Marshal(config)
	if err != nil {
		return err
	}

	if err := logger.SetLogger(mode, string(c)); err != nil {
		return err
	}
	loggers[mode] = logger
	return nil
}

// Console returns a function which can be used on the Init function
// to register a 'console' logger
func Console(level LogLevel) func() error {
	type consoleConf struct {
		Level LogLevel `json:"level"`
	}
	return func() error {
		conf := &consoleConf{
			Level: level,
		}
		return RegisterLogger("console", DefaultBufferSize, conf)
	}
}

// File returns a function which can be used on the Init function
// to register a 'file' logger
func File(level LogLevel, filename string) func() error {
	type fileConf struct {
		Level       LogLevel `json:"level"`
		Filename    string   `json:"filename"`
		MaxLines    int      `json:"maxlines"`
		MaxSize     int64    `json:"maxsize"`
		DailyRotate bool     `json:"daily"`
		MaxDays     int      `json:"maxdays"`
		Rotate      bool     `json:"rotate"`
	}

	return func() error {
		conf := &fileConf{
			Level:       level,
			Filename:    filename,
			MaxLines:    0,
			MaxSize:     0,
			DailyRotate: true,
			MaxDays:     7,
			Rotate:      true,
		}
		return RegisterLogger("file", DefaultBufferSize, conf)
	}
}

// Init initializes loggers with given options
// Each option is a function which registers a logger. For example:
//
// err := Init(Console(LogLevelDebug), File(LogLevelInfo, "/tmp/myapp.log"))
func Init(options ...func() error) error {
	mutex.Lock()
	defer mutex.Unlock()

	// close all existing loggers first
	Close()

	// register loggers
	for _, option := range options {
		if err := option(); err != nil {
			return err
		}
	}
	return nil
}
