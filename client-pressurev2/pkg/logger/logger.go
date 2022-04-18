package logger

import (
	"github.com/rs/zerolog"
	"io"
	"os"
	"time"
)

func NewSysLog(out io.Writer) zerolog.Logger {
	var logger zerolog.Logger

	zerolog.TimeFieldFormat = time.RFC3339Nano
	out = io.MultiWriter(os.Stderr, out)
	logger = zerolog.New(out).With().Timestamp().Logger()
	return logger
}

func NewMsgLog(out io.Writer) zerolog.Logger {
	var logger zerolog.Logger

	zerolog.TimeFieldFormat = time.RFC3339Nano
	out = io.MultiWriter(out)
	logger = zerolog.New(out).With().Timestamp().Logger()
	return logger
}
