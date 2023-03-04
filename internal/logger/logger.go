package logger

import (
	"os"

	log "github.com/rs/zerolog"
)

type Logger struct {
	log.Logger
	file *os.File
}

var (
	logFile *os.File
)

func New(logFileName string) (*Logger, error) {
	if logFileName == "" {
		return &Logger{
			Logger: log.New(os.Stdout).With().Timestamp().Logger(),
		}, nil
	}

	var err error
	logFile, err = os.OpenFile(logFileName, os.O_WRONLY|os.O_CREATE, 0755)
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
