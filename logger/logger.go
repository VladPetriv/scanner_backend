package logger

import (
	"fmt"
	"io"
	"os"
	"path"
	"runtime"

	"github.com/sirupsen/logrus"
)

var e *logrus.Entry // nolint

type Logger struct {
	*logrus.Entry
}

func Get() *Logger {
	Init()

	return &Logger{e}
}

func Init() {
	log := logrus.New()
	log.SetReportCaller(true)
	log.Formatter = &logrus.TextFormatter{ // nolint
		CallerPrettyfier: func(f *runtime.Frame) (string, string) {
			filename := path.Base(f.File)

			return fmt.Sprintf("%s:%d", filename, f.Line), fmt.Sprintf("%s()", f.Function)
		},
		DisableColors: false,
		FullTimestamp: true,
	}

	err := os.Mkdir("logs", 0o755) // nolint
	if os.IsNotExist(err) {
		panic(err)
	}

	allFile, err := os.OpenFile("logs/all.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0660)
	if err != nil {
		panic(fmt.Sprintf("[Error]: %s", err))
	}

	log.SetOutput(io.MultiWriter(allFile, os.Stdout)) // Send all logs to nowhere by default

	log.SetLevel(logrus.TraceLevel)

	e = logrus.NewEntry(log)
}
