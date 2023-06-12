package http

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/txchat/dtalk/api/proto/chat"
)

func TestGetMid(t *testing.T) {
	_, err := GetMid(&SendResp{
		Code:     int32(chat.SendMessageReply_InnerError),
		Mid:      "1",
		Datetime: 0,
		Repeat:   false,
	})
	assert.True(t, errors.As(err, &ErrSendMessageReplyInnerError))
}
