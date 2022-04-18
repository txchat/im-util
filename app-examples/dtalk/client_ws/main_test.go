package main

import (
	"encoding/base64"
	"testing"

	"github.com/Terry-Mao/goim/pkg/encoding/binary"
	"github.com/golang/protobuf/proto"
	comet "github.com/txchat/im/api/comet/grpc"
)

const (
	// size
	_packSize      = 4
	_headerSize    = 2
	_verSize       = 2
	_opSize        = 4
	_seqSize       = 4
	_ackSize       = 4
	_heartSize     = 4
	_rawHeaderSize = _packSize + _headerSize + _verSize + _opSize + _seqSize + _ackSize
	// offset
	_packOffset   = 0
	_headerOffset = _packOffset + _packSize
	_verOffset    = _headerOffset + _headerSize
	_opOffset     = _verOffset + _verSize
	_seqOffset    = _opOffset + _opSize
	_ackOffset    = _seqOffset + _seqSize
	_heartOffset  = _ackOffset + _ackSize
)

func testWt(p *comet.Proto) []byte {
	var (
		buf     []byte
		packLen int32
	)
	packLen = _rawHeaderSize + int32(len(p.Body))
	buf = make([]byte, packLen)
	binary.BigEndian.PutInt32(buf[_packOffset:], packLen)
	binary.BigEndian.PutInt16(buf[_headerOffset:], int16(_rawHeaderSize))
	binary.BigEndian.PutInt16(buf[_verOffset:], int16(p.Ver))
	binary.BigEndian.PutInt32(buf[_opOffset:], p.Op)
	binary.BigEndian.PutInt32(buf[_seqOffset:], p.Seq)
	binary.BigEndian.PutInt32(buf[_ackOffset:], p.Ack)
	copy(buf[_heartOffset:], p.Body)
	return buf
}

func Test_AuthFlame(t *testing.T) {
	authMsg := &comet.AuthMsg{
		AppId: "dtalk",
		Token: "WYXIXva0Q2vPSTECSMR7Shw/uZIYDE88ostB5DOe7QpCZOZ2LlzUPK0HhyhK5ffo1l9XygADG6+qzoE+3mvH8AE=#1613988388311*UEcIzO3v#0375610055c57e011a0a51457e0ce451849a4ca588b0ff0beb0ba5d929ca2dd82b",
	}
	p := new(comet.Proto)
	p.Ver = 1
	p.Op = int32(comet.Op_Auth)
	p.Seq = seq
	p.Body, _ = proto.Marshal(authMsg)
	t.Log(base64.StdEncoding.EncodeToString(testWt(p)))
}
