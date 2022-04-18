package logger

import (
	"github.com/natefinch/lumberjack"
	"github.com/rs/zerolog"
	"io"
	"os"
	"time"
)

var Log = zerolog.New(os.Stderr).With().Timestamp().Logger()

func getWriter(filename string) io.Writer {
	return &lumberjack.Logger{
		Filename:   filename,
		MaxSize:    10,    //最大M数，超过则切割
		MaxBackups: 5,     //最大文件保留数，超过就删除最老的日志文件
		MaxAge:     30,    //保存30天
		Compress:   false, //是否压缩
	}
}

func NewLog(filename string) zerolog.Logger {
	var logger zerolog.Logger
	var out io.Writer

	if zerolog.GlobalLevel() == zerolog.DebugLevel {
		out = io.MultiWriter(zerolog.ConsoleWriter{Out: os.Stderr, TimeFormat: time.RFC3339Nano}, getWriter(filename))
		logger = zerolog.New(out).With().Timestamp().Logger()
	} else {
		zerolog.TimeFieldFormat = time.RFC3339Nano
		out = io.MultiWriter(os.Stderr, getWriter(filename))
		logger = zerolog.New(out).With().Timestamp().Logger()
	}

	return logger
}

func NewMsgLog(filename string) zerolog.Logger {
	var logger zerolog.Logger
	var out io.Writer

	if zerolog.GlobalLevel() == zerolog.DebugLevel {
		out = io.MultiWriter(zerolog.ConsoleWriter{Out: os.Stderr, TimeFormat: time.RFC3339Nano}, getWriter(filename))
	} else {
		zerolog.TimeFieldFormat = time.RFC3339Nano
		out = io.MultiWriter(getWriter(filename))
	}
	logger = zerolog.New(out).With().Logger()
	return logger
}

func Init(debug bool) {
	zerolog.DurationFieldUnit = 10 * time.Second
	zerolog.TimeFieldFormat = time.RFC3339Nano
	zerolog.SetGlobalLevel(zerolog.InfoLevel)
	if debug {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
		Log = zerolog.New(zerolog.ConsoleWriter{Out: os.Stderr, TimeFormat: time.RFC3339Nano}).With().Timestamp().Logger()
	}
}
