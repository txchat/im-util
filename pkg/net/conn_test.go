package net

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/txchat/im/api/protocol"
)

func TestDialIMAndServe(t *testing.T) {
	conn, err := DialIMAndServe("mock_server", &protocol.AuthBody{
		AppId: "mock",
		Token: "TestDialIMAndServe",
		Ext:   nil,
	}, time.Second, NewMockAuth)
	assert.Nil(t, err)
	defer conn.Close()
	p := conn.Read()
	assert.EqualValues(t, protocol.Op_Message, p.GetOp())
}
