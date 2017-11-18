// Copyright 2013 The beego Authors.
// Copyright 2014 The Gogs Authors.
//
// Licensed under the Apache License, Version 2.0 (the "License"): you may
// not use this file except in compliance with the License. You may obtain
// a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS, WITHOUT
// WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the
// License for the specific language governing permissions and limitations
// under the License.

package logs

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"sync"
)

const (
	// log message levels
	LevelTrace = iota
	LevelDebug
	LevelInfo
	LevelWarn
	LevelError
	LevelCritical
	LevelFatal
)

type loggerType func() LoggerInterface

// LoggerInterface defines the behavior of a log provider.
type LoggerInterface interface {
	Init(config string) error
	WriteMsg(msg string, level int) error
	Destroy()
	Flush()
}

var adapters = make(map[string]loggerType)

// Register makes a log provide available by the provided name.
// If Register is called twice with the same name or if driver is nil,
// it panics.
func Register(name string, log loggerType) {
	if log == nil {
		panic("logs: Register provide is nil")
	}
	if _, dup := adapters[name]; dup {
		panic("logs: Register called twice for provider " + name)
	}
	adapters[name] = log
}

// BeeLogger is default logger in beego application.
// it can contain several providers and log message into all providers.
type BeeLogger struct {
	Adapter string

	lock                sync.Mutex
	level               int
	enableFuncCallDepth bool
	loggerFuncCallDepth int
	msg                 chan *logMsg
	outputs             map[string]LoggerInterface
	quit                chan bool
}

type logMsg struct {
	level int
	msg   string
}

// NewLogger returns a new BeeLogger.
// channellen means the number of messages in chan.
// if the buffering chan is full, logger adapters write to file or other way.
func NewLogger(channellen int64) *BeeLogger {
	bl := new(BeeLogger)
	bl.enableFuncCallDepth = true
	bl.loggerFuncCallDepth = 4
	bl.msg = make(chan *logMsg, channellen)
	bl.outputs = make(map[string]LoggerInterface)
	bl.quit = make(chan bool)
	go bl.StartLogger()
	return bl
}

// SetLogger provides a given logger adapter into BeeLogger with config string.
// config need to be correct JSON as string: {"interval":360}.
func (bl *BeeLogger) SetLogger(adaptername string, config string) error {
	bl.lock.Lock()
	defer bl.lock.Unlock()
	if log, ok := adapters[adaptername]; ok {
		lg := log()
		if err := lg.Init(config); err != nil {
			return err
		}
		bl.outputs[adaptername] = lg
		bl.Adapter = adaptername
		return nil
	} else {
		return fmt.Errorf("logs: unknown adaptername %q (forgotten Register?)", adaptername)
	}
}

// remove a logger adapter in BeeLogger.
func (bl *BeeLogger) DelLogger(adaptername string) error {
	bl.lock.Lock()
	defer bl.lock.Unlock()
	if lg, ok := bl.outputs[adaptername]; ok {
		lg.Destroy()
		delete(bl.outputs, adaptername)
		return nil
	} else {
		return fmt.Errorf("logs: unknown adaptername %q (forgotten Register?)", adaptername)
	}
}

func (bl *BeeLogger) writerMsg(loglevel int, msg string) error {
	if bl.level > loglevel {
		return nil
	}
	lm := new(logMsg)
	lm.level = loglevel
	if bl.enableFuncCallDepth {
		_, file, line, ok := runtime.Caller(bl.loggerFuncCallDepth)
		if ok {
			_, file = filepath.Split(file)
			lm.msg = fmt.Sprintf("[%s:%d] %s", file, line, msg)
		} else {
			lm.msg = msg
		}
	} else {
		lm.msg = msg
	}
	bl.msg <- lm
	return nil
}

// set log message level.
// if message level (such as LevelTrace) is less than logger level (such as LevelWarn), ignore message.
func (bl *BeeLogger) SetLevel(l int) {
	bl.level = l
}

// set log funcCallDepth
func (bl *BeeLogger) SetLogFuncCallDepth(d int) {
	bl.loggerFuncCallDepth = d
}

// enable log funcCallDepth
func (bl *BeeLogger) EnableFuncCallDepth(b bool) {
	bl.enableFuncCallDepth = b
}

// start logger chan reading.
// when chan is full, write logs.
func (bl *BeeLogger) StartLogger() {
	for {
		select {
		case bm := <-bl.msg:
			for _, l := range bl.outputs {
				l.WriteMsg(bm.msg, bm.level)
			}
		case <-bl.quit:
			return
		}
	}

}

// log trace level message.
func (bl *BeeLogger) Trace(format string, v ...interface{}) {
	msg := fmt.Sprintf("[T] "+format, v...)
	bl.writerMsg(LevelTrace, msg)
}

// log debug level message.
func (bl *BeeLogger) Debug(format string, v ...interface{}) {
	msg := fmt.Sprintf("[D] "+format, v...)
	bl.writerMsg(LevelDebug, msg)
}

// log info level message.
func (bl *BeeLogger) Info(format string, v ...interface{}) {
	msg := fmt.Sprintf("[I] "+format, v...)
	bl.writerMsg(LevelInfo, msg)
}

// log warn level message.
func (bl *BeeLogger) Warn(format string, v ...interface{}) {
	msg := fmt.Sprintf("[W] "+format, v...)
	bl.writerMsg(LevelWarn, msg)
}

// log error level message.
func (bl *BeeLogger) Error(format string, v ...interface{}) {
	msg := fmt.Sprintf("[E] "+format, v...)
	bl.writerMsg(LevelError, msg)
}

// log critical level message.
func (bl *BeeLogger) Critical(format string, v ...interface{}) {
	msg := fmt.Sprintf("[C] "+format, v...)
	bl.writerMsg(LevelCritical, msg)
}

func (bl *BeeLogger) Fatal(format string, v ...interface{}) {
	msg := fmt.Sprintf("[F] "+format, v...)
	bl.writerMsg(LevelFatal, msg)
	bl.Close()
	os.Exit(1)
}

// flush all chan data.
func (bl *BeeLogger) Flush() {
	for _, l := range bl.outputs {
		l.Flush()
	}
}

// close logger, flush all chan data and destroy all adapters in BeeLogger.
func (bl *BeeLogger) Close() {
	bl.quit <- true
	for {
		if len(bl.msg) > 0 {
			bm := <-bl.msg
			for _, l := range bl.outputs {
				l.WriteMsg(bm.msg, bm.level)
			}
		} else {
			break
		}
	}
	for _, l := range bl.outputs {
		l.Flush()
		l.Destroy()
	}
}
