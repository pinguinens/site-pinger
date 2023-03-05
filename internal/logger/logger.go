package logger

import (
	"fmt"
	"os"
	"strings"
	"time"

	log "github.com/rs/zerolog"
)

const (
	timeLayout = "2006-01-02T15-04-05"
)

type Logger struct {
	log.Logger
	file *os.File
}

func New(logFileName string) (*Logger, error) {
	if logFileName == "" {
		return &Logger{
			Logger: log.New(os.Stdout).With().Timestamp().Logger(),
		}, nil
	}

	var err error
	fn := logFileName
	if strings.Contains(logFileName, "%v") {
		fn = fmt.Sprintf(logFileName, time.Now().Format(timeLayout))
	}
	logFile, err := os.OpenFile(fn, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0755)
	if err != nil {
		return nil, err
	}

	multi := log.MultiLevelWriter(os.NewFile(logFile.Fd(), logFile.Name()), os.Stdout)
	logger := log.New(multi).With().Timestamp().Logger()

	return &Logger{
		Logger: logger,
		file:   logFile,
	}, nil
}

func (l *Logger) Close() error {
	return l.file.Close()
}
