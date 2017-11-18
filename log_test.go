package log

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"strings"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestClose(t *testing.T) {
	Init(Console(LogLevelTrace))
	assert.Equal(t, 1, len(loggers))

	Close()
	assert.Equal(t, 0, len(loggers))
}

func TestInit(t *testing.T) {
	err := Init(Console(LogLevelTrace))
	assert.NoError(t, err)
	assert.Equal(t, 1, len(loggers))

	tmpdir := os.TempDir()
	err2 := Init(Console(LogLevelDebug), File(LogLevelDebug, path.Join(tmpdir, "test_logging.log")))
	assert.NoError(t, err2)
	assert.Equal(t, 2, len(loggers))

	err3 := Init(Console(LogLevelDebug), Console(LogLevelWarn))
	assert.NoError(t, err3)
	assert.Equal(t, 1, len(loggers))
}

func TestConsoleLogger(t *testing.T) {
	err := Init(Console(LogLevelDebug))
	assert.NoError(t, err)

	Traceln("Trace")
	Debugln("Debug")
	Infoln("Info")
	Warnln("Warn")
	Errorln("Error")
	Criticalln("Critical")

	// close all loggers
	Close()
}

func TestFileLogger(t *testing.T) {
	logfile := filepath.Join(os.TempDir(), "test_logger.log")
	defer os.RemoveAll(logfile)

	err := Init(File(LogLevelDebug, logfile))
	assert.NoError(t, err)

	Traceln("Trace")
	Debugln("Debug")
	Infoln("Info")
	Warnln("Warn")
	Errorln("Error")
	Criticalln("Critical")

	// close all loggers
	Close()

	f, err := os.Open(logfile)
	assert.NoError(t, err)
	content, err := ioutil.ReadAll(f)
	assert.NoError(t, err)

	assert.False(t, strings.Contains(string(content), "Trace"))
	assert.True(t, strings.Contains(string(content), "Debug"))
	assert.True(t, strings.Contains(string(content), "Info"))
	assert.True(t, strings.Contains(string(content), "Warn"))
	assert.True(t, strings.Contains(string(content), "Error"))
	assert.True(t, strings.Contains(string(content), "Critical"))
	f.Close()
}

func TestLogFuncDepth(t *testing.T) {
	logfile := filepath.Join(os.TempDir(), "test_func_depth.log")
	defer os.RemoveAll(logfile)

	err := Init(File(LogLevelDebug, logfile))
	assert.NoError(t, err)

	var wg sync.WaitGroup
	wg.Add(1)

	// write logs from a goroutine
	go func() {
		defer wg.Done()

		Debug("Message A")
	}()

	Debug("Message B")

	wg.Wait()
	Close()

	f, err := os.Open(logfile)
	assert.NoError(t, err)
	content, err := ioutil.ReadAll(f)
	fmt.Print(string(content))
	assert.Contains(t, string(content), "Message A")
	assert.Contains(t, string(content), "Message B")
	assert.Contains(t, string(content), "log_test.go")
	assert.NotContains(t, string(content), "log.go")

	f.Close()
}
