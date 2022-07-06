package log

import (
	"github.com/txchat/im-util/pkg/util"
	"os"
	"testing"
)

func Test_ConsoleLogger(t *testing.T) {
	logger := NewLogger(os.Stderr)
	logger.Info().Str("foo", "bee").Msg("hello world")
	t.Log("success")
}

func Test_FileLogger(t *testing.T) {
	fd, closer, err := util.WriteFile("./log.txt")
	if err != nil {
		t.Error(err)
		return
	}
	defer closer()
	logger := NewLogger(fd)
	logger.Info().Str("foo", "bee").Msg("hello world")
	t.Log("success")
}
