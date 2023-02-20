package logger

import (
	"os"

	log "github.com/rs/zerolog"
)

var (
	logFile *os.File
)

func New(logFileName string) (*log.Logger, error) {
	var err error
	logFile, err = os.OpenFile(logFileName, os.O_WRONLY|os.O_CREATE, 0755)
	if err != nil {
		return nil, err
	}

	multi := log.MultiLevelWriter(os.NewFile(logFile.Fd(), logFile.Name()), os.Stdout)
	logger := log.New(multi).With().Timestamp().Logger()

	return &logger, nil
}

func CloseLogFile() error {
	return logFile.Close()
}
