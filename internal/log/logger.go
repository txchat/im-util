package log

import (
	"io"
	"time"

	"github.com/rs/zerolog"
)

func NewLogger(out ...io.Writer) zerolog.Logger {
	var logger zerolog.Logger

	zerolog.TimeFieldFormat = time.RFC3339Nano
	w := io.MultiWriter(out...)
	logger = zerolog.New(w).With().Timestamp().Logger()
	return logger
}
