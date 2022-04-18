package logger

import (
	"testing"
)

func TestNewMsgLog(t *testing.T) {
	logger := NewMsgLog("./testlog.txt")
	logger.Info().Str("foo", "bee").Msg("hello world")
	t.Log("success")
}

func TestNewLog(t *testing.T) {
	logger := NewLog("./testlog2.txt")
	logger.Info().Str("foo", "bee").Msg("hello world")
	t.Log("success")
}

func TestInit(t *testing.T) {
	Init(true)
	Log.Info().Str("foo", "bee").Msg("hello world")
	t.Log("success")
}
