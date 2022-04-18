package logger

import (
	"github.com/txchat/im-util/client-pressurev2/pkg/filehelper"
	"testing"
)

func TestNewSysLog(t *testing.T) {
	fd, closer, err := filehelper.WriteFile("./testlog.txt")
	if err != nil {
		t.Error(err)
		return
	}
	defer closer()
	logger := NewSysLog(fd)
	logger.Info().Str("foo", "bee").Msg("hello world")
	t.Log("success")
}

func TestNewMsgLog(t *testing.T) {
	fd, closer, err := filehelper.WriteFile("./testlog2.txt")
	if err != nil {
		t.Error(err)
		return
	}
	defer closer()
	logger := NewMsgLog(fd)
	logger.Info().Str("foo", "bee").Msg("hello world")
	t.Log("success")
}
