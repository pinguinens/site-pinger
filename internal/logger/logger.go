package logger

import (
	"fmt"
	"io"
	"os"
	"strings"
	"time"

	log "github.com/rs/zerolog"
)

const (
	timeLayout = "2006-01-02T15-04-05"

	PrettyFormat = "pretty"
	JsonFormat   = "json"
)

type Logger struct {
	log.Logger
	file *os.File
}

func New(logFileName, format string) (*Logger, error) {
	if logFileName == "" {
		return &Logger{
			Logger: log.New(makeConsoleWriter(format)).With().Timestamp().Logger(),
		}, nil
	}

	file, err := openLogFile(logFileName)
	if err != nil {
		return nil, err
	}
	multi := log.MultiLevelWriter(makeFileWriter(file), makeConsoleWriter(format))

	return &Logger{
		Logger: log.New(multi).With().Timestamp().Logger(),
		file:   file,
	}, nil
}

func (l *Logger) Close() error {
	return l.file.Close()
}

func openLogFile(fileName string) (*os.File, error) {
	fn := fileName
	if strings.Contains(fileName, "%v") {
		fn = fmt.Sprintf(fileName, time.Now().Format(timeLayout))
	}

	return os.OpenFile(fn, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0755)
}

func makeConsoleWriter(format string) io.Writer {
	if format == PrettyFormat {
		output := log.ConsoleWriter{Out: os.Stdout, TimeFormat: time.RFC3339}
		output.FormatLevel = func(i interface{}) string {
			if i == nil {
				return ""
			}
			return strings.ToUpper(fmt.Sprintf("| %-6s|", i))
		}
		output.FormatMessage = func(i interface{}) string {
			if i == nil {
				return ""
			}
			return fmt.Sprintf("***%s****", i)
		}
		output.FormatFieldName = func(i interface{}) string {
			return fmt.Sprintf("| %s:", i)
		}
		output.FormatFieldValue = func(i interface{}) string {
			return fmt.Sprintf("%s", i)
		}
		return output
	}

	return os.Stdout
}

func makeFileWriter(file *os.File) io.Writer {
	return os.NewFile(file.Fd(), file.Name())
}
