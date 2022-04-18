package internal

import (
	"github.com/golang/protobuf/proto"
	"github.com/google/uuid"
	bizProto "github.com/txchat/imparse/proto"
)

//
func TextMsg(chType int32, from, target, msg string) *bizProto.Proto {
	m := bizProto.TextMsg{
		Content: msg,
	}
	data, err := proto.Marshal(&m)
	if err != nil {
		panic(err)
	}
	body, err := proto.Marshal(&bizProto.CommonMsg{
		ChannelType: chType,
		LogId:       0,
		MsgId:       uuid.New().String(),
		From:        from,
		Target:      target,
		MsgType:     int32(bizProto.MsgType_Text),
		Msg:         data,
		Datetime:    0,
	})
	if err != nil {
		panic(err)
	}
	return &bizProto.Proto{
		EventType: bizProto.EventType_commonMsg,
		Body:      body,
	}
}
