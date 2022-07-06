package log

import (
	"github.com/rs/zerolog"
	"io"
	"time"
)

func NewLogger(out ...io.Writer) zerolog.Logger {
	var logger zerolog.Logger

	zerolog.TimeFieldFormat = time.RFC3339Nano
	w := io.MultiWriter(out...)
	logger = zerolog.New(w).With().Timestamp().Logger()
	return logger
}
