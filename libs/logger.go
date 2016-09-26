package libs

import (
	"io"
	"os"
	"path/filepath"
	"sync"
	"time"
	"webcron-agent/config"

	logging "github.com/op/go-logging"
)

type Logger struct {
	log *logging.Logger
}

var lock sync.RWMutex

func NewTaskLogger() *Logger {
	now := time.Now()
	f, _ := CreateLogFile("/var/logs/task_" + now.Format("20060102") + ".log")
	return NewLogger(f)
}

func NewSocketLogger() *Logger {
	now := time.Now()
	f, _ := CreateLogFile("/var/logs/socket_" + now.Format("20060102") + ".log")
	return NewLogger(f)
}

func NewAppLogger() *Logger {
	f, _ := CreateLogFile("/var/logs/app.log")
	return NewLogger(f)
}

func CreateLogFile(file string) (io.Writer, error) {
	lock.Lock()
	defer lock.Unlock()
	config, _ := config.GetConfig()

	logpath := filepath.Join(config.App.Path, file)
	logdir := filepath.Dir(logpath)

	f, _ := os.Lstat(logpath)
	dir, _ := os.Lstat(logdir)

	if f == nil && dir != nil {
		f, err := os.Create(logpath)
		defer f.Close()
		if err != nil {
			return nil, err
		}
	}

	if dir == nil {
		os.MkdirAll(logdir, 0777)
		f, err := os.Create(logpath)
		defer f.Close()
		if err != nil {
			return nil, err
		}
	}

	return os.OpenFile(logpath, os.O_APPEND|os.O_WRONLY, 0600)
}

func NewLogger(out io.Writer) *Logger {
	var log = logging.MustGetLogger("example")
	var format = logging.MustStringFormatter(
		`%{color}%{time:15:04:05.000} %{shortfunc} -> %{level:.4s} %{id:03x}%{color:reset} %{message}`,
	)
	var formatfile = logging.MustStringFormatter(
		`%{time:15:04:05.000} %{shortfunc} -> %{level:.4s} %{color:reset} %{message}`,
	)

	//one
	stdout := logging.NewLogBackend(os.Stdout, "", 0)
	stdoutFormat := logging.NewBackendFormatter(stdout, format)

	//two
	backend := logging.NewLogBackend(out, "", 0)
	backendFormat := logging.NewBackendFormatter(backend, formatfile)

	logging.SetBackend(backendFormat, stdoutFormat)

	return &Logger{log: log}
}

func (l *Logger) Info(args ...interface{}) {
	l.log.Info(args)
}

func (l *Logger) Warning(args ...interface{}) {
	l.log.Warning(args)
}

func (l *Logger) Notice(args ...interface{}) {
	l.log.Notice(args)
}

func (l *Logger) Error(args ...interface{}) {
	l.log.Error(args)
}
