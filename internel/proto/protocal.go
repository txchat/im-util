package proto

import (
	"time"

	"github.com/golang/protobuf/proto"
	"github.com/google/uuid"
	comet "github.com/txchat/im/api/comet/grpc"
	xproto "github.com/txchat/imparse/proto"
)

func CreateProtoSendMsg(seq int32, from, target string, channelType xproto.Channel, text string) (*comet.Proto, error) {
	p := new(comet.Proto)
	p.Op = int32(comet.Op_SendMsg)
	p.Seq = seq
	p.Ack = 0
	msgType, data, err := Text(text)
	if err != nil {
		return nil, err
	}
	commData, err := Common(channelType, from, target, msgType, data)
	if err != nil {
		return nil, err
	}
	protoData, err := Proto(xproto.Proto_common, commData)
	if err != nil {
		return nil, err
	}
	body, err := proto.Marshal(protoData)
	if err != nil {
		return nil, err
	}
	p.Body = body

	return p, nil
}

func CreateProtoAck(seq int32, Ack int32) (*comet.Proto, error) {
	p := new(comet.Proto)
	p.Op = int32(comet.Op_ReceiveMsgReply)
	p.Seq = seq
	p.Ack = Ack
	p.Body = nil

	return p, nil
}

func Text(msg string) (xproto.MsgType, []byte, error) {
	m := xproto.TextMsg{
		Content: msg,
	}
	data, err := proto.Marshal(&m)
	return xproto.MsgType_Text, data, err
}

func Common(channelType xproto.Channel, from, target string, msgType xproto.MsgType, msgData []byte) ([]byte, error) {
	msgId := uuid.New().String()
	return proto.Marshal(&xproto.Common{
		ChannelType: channelType,
		Mid:         0,
		Seq:         msgId,
		From:        from,
		Target:      target,
		MsgType:     msgType,
		Msg:         msgData,
		Datetime:    uint64(time.Now().UnixNano() / 1e6),
	})
}

func Proto(eventType xproto.Proto_EventType, body []byte) (*xproto.Proto, error) {
	return &xproto.Proto{
		EventType: eventType,
		Body:      body,
	}, nil
}

func ConvertCommonAck(body []byte) (*xproto.CommonAck, error) {
	var p xproto.Proto
	err := proto.Unmarshal(body, &p)
	if err != nil {
		return nil, err
	}
	var commonAck xproto.CommonAck
	err = proto.Unmarshal(p.GetBody(), &commonAck)
	return &commonAck, err
}

func ConvertBizProto(body []byte) (*xproto.Proto, error) {
	var p xproto.Proto
	return &p, proto.Unmarshal(body, &p)
}

func ConvertCommon(bizBody []byte) (*xproto.Common, error) {
	var common xproto.Common
	return &common, proto.Unmarshal(bizBody, &common)
}
