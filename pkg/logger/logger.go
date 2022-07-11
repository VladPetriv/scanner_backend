package logger

import (
	"fmt"
	"io"
	"os"
	"path"
	"runtime"

	"github.com/sirupsen/logrus"

	"github.com/VladPetriv/scanner_backend/pkg/config"
)

var e *logrus.Entry

type Logger struct {
	*logrus.Entry
}

func Get() *Logger {
	cfg, err := config.Get()
	if err != nil {
		panic(err)
	}

	Init(cfg.LogLevel)

	return &Logger{e}
}

func Init(logLevel string) {
	level, err := logrus.ParseLevel(logLevel)
	if err != nil {
		panic(err)
	}
	log := logrus.New()

	log.SetReportCaller(true)

	log.Formatter = &logrus.JSONFormatter{
		CallerPrettyfier: func(f *runtime.Frame) (string, string) {
			filename := path.Base(f.File)

			return fmt.Sprintf("%s:%d", filename, f.Line), fmt.Sprintf("%s()", f.Function)
		},
	}

	err = os.Mkdir("logs", 0o755)
	if os.IsNotExist(err) {
		panic(err)
	}

	file, err := os.OpenFile("logs/all.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0660)
	if err != nil {
		panic(fmt.Sprintf("[Error]: %s", err))
	}

	log.SetOutput(io.MultiWriter(file, os.Stdout))

	log.SetLevel(level)

	e = logrus.NewEntry(log)
}
