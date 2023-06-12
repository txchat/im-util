package proto

import (
	"github.com/txchat/dtalk/api/proto/content"
	"github.com/txchat/dtalk/api/proto/message"

	"github.com/golang/protobuf/proto"
	"github.com/txchat/im/api/protocol"
)

func CreateProtoAck(ack int32) (*protocol.Proto, error) {
	p := new(protocol.Proto)
	p.Op = int32(protocol.Op_MessageReply)
	p.Seq = 0
	p.Ack = ack
	p.Body = nil
	return p, nil
}

func Text(msg string, ait ...string) (message.MsgType, []byte, error) {
	m := content.TextMsg{
		Content: msg,
		Ait:     ait,
	}
	data, err := proto.Marshal(&m)
	return message.MsgType_Text, data, err
}
