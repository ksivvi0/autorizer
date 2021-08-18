package services

import (
	"log"
	"os"
)

type LoggerService interface {
	WriteError(data string)
	WriteNotice(data string)
}

type Logger struct {
	errors  *log.Logger
	notices *log.Logger
}

func NewLoggerInstance(file string) (*Logger, error) {
	f, err := os.OpenFile(file, os.O_APPEND|os.O_CREATE, 0640)
	if err != nil {
		return nil, err
	}

	eLogger := log.New(f, "ERROR: ", log.Ltime|log.Lshortfile)
	nLogger := log.New(f, "NOTICE: ", log.Ltime|log.Lshortfile)
	return &Logger{
		errors:  eLogger,
		notices: nLogger,
	}, nil
}

func (l *Logger) WriteError(data string) {
	l.errors.Println(data)
}

func (l *Logger) WriteNotice(data string) {
	l.notices.Println(data)
}
