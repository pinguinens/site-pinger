package logger

import (
	"os"

	log "github.com/rs/zerolog"
)

func New(outputFile *os.File) *log.Logger {
	multi := log.MultiLevelWriter(os.NewFile(outputFile.Fd(), outputFile.Name()), os.Stdout)
	logger := log.New(multi).With().Timestamp().Logger()

	return &logger
}
