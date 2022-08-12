package logging

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"io"
	"os"
	"path"
	"runtime"
)

const (
	logDir  = "/app_logs"
	logFile = "/all.log"
)

type writerHook struct {
	Writer    []io.Writer
	LogLevels []logrus.Level
}

func (hook *writerHook) Fire(entry *logrus.Entry) error {
	line, err := entry.String()
	if err != nil {
		return err
	}

	for _, w := range hook.Writer {
		_, err := w.Write([]byte(line))
		if err != nil {
			return err
		}
	}
	return err
}

func (hook *writerHook) Levels() []logrus.Level {
	return hook.LogLevels
}

var e *logrus.Entry

type Logger struct {
	*logrus.Entry
}

// GetLogger возращаем логгер-обертку над логгером logrus
func GetLogger() *Logger {
	return &Logger{e}
}

// GetLoggerWithField для создания логгера содержащего определённые поля для вывода
func (l *Logger) GetLoggerWithField(k string, v interface{}) *Logger {
	return &Logger{l.WithField(k, v)}
}

func Init() {
	l := logrus.New()
	l.SetReportCaller(true)
	//найстройка формата записи логов
	l.Formatter = &logrus.TextFormatter{
		CallerPrettyfier: func(frame *runtime.Frame) (function string, file string) {
			filename := path.Base(frame.File)
			return fmt.Sprintf("%s()", frame.Function), fmt.Sprintf("%s:%d", filename, frame.Line)
		},
		FullTimestamp: true,
	}

	prjDir := os.Getenv("PROJECT_DIR")
	logDirPath := prjDir + logDir

	//создание директории для хранения логов с правами -rw_-r__-r__
	err := os.MkdirAll(logDirPath, 0644)
	if err != nil {
		panic(any(err))
	}

	allFiles, err := os.OpenFile(logDirPath+logFile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0640)
	if err != nil {
		panic(any(err))
	}

	//запись в никуда
	l.SetOutput(io.Discard)

	//целевая история для внесения изменений в дефолтную реализацию Writer`ов logrus`а
	l.AddHook(&writerHook{
		Writer:    []io.Writer{allFiles, os.Stdout},
		LogLevels: logrus.AllLevels,
	})

	l.SetLevel(logrus.TraceLevel)
	e = logrus.NewEntry(l)
}
