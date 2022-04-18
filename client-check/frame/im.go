package frame

import (
	"fmt"
	"github.com/Terry-Mao/goim/pkg/encoding/binary"
	"github.com/txchat/im-util/client-check/model"
	comet "github.com/txchat/im/api/comet/grpc"
)

const (
	// MaxBodySize max proto body size
	MaxBodySize = int32(1 << 12)
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
	_maxPackSize   = MaxBodySize + int32(_rawHeaderSize)
	// offset
	_packOffset   = 0
	_headerOffset = _packOffset + _packSize
	_verOffset    = _headerOffset + _headerSize
	_opOffset     = _verOffset + _verSize
	_seqOffset    = _opOffset + _opSize
	_ackOffset    = _seqOffset + _seqSize
	_heartOffset  = _ackOffset + _ackSize
)

func ToProto(buf []byte) (p comet.Proto, err error) {
	var (
		bodyLen   int
		headerLen int16
		packLen   int32
	)
	if len(buf) < _rawHeaderSize {
		return p, model.ErrProtoPackLen
	}
	packLen = binary.BigEndian.Int32(buf[_packOffset:_headerOffset])
	headerLen = binary.BigEndian.Int16(buf[_headerOffset:_verOffset])
	p.Ver = int32(binary.BigEndian.Int16(buf[_verOffset:_opOffset]))
	p.Op = binary.BigEndian.Int32(buf[_opOffset:_seqOffset])
	p.Seq = binary.BigEndian.Int32(buf[_seqOffset:_ackOffset])
	p.Ack = binary.BigEndian.Int32(buf[_ackOffset:])
	if packLen < 0 || packLen > _maxPackSize {
		return p, model.ErrProtoPackLen
	}
	if headerLen != _rawHeaderSize {
		return p, model.ErrProtoHeaderLen
	}
	if bodyLen = int(packLen - int32(headerLen)); bodyLen > 0 {
		p.Body = buf[headerLen:packLen]
	} else {
		p.Body = nil
	}
	if bodyLen != len(buf[headerLen:]) {
		fmt.Printf("got %v, need %v\n", len(buf[headerLen:]), bodyLen)
		return p, model.ErrProtoPackLen
	}
	return
}

func ToBytes(p comet.Proto) []byte {
	var (
		buf     []byte
		packLen int
	)
	packLen = _rawHeaderSize + len(p.Body)
	buf = make([]byte, packLen)
	binary.BigEndian.PutInt32(buf[_packOffset:], int32(packLen))
	binary.BigEndian.PutInt16(buf[_headerOffset:], int16(_rawHeaderSize))
	binary.BigEndian.PutInt16(buf[_verOffset:], int16(p.Ver))
	binary.BigEndian.PutInt32(buf[_opOffset:], p.Op)
	binary.BigEndian.PutInt32(buf[_seqOffset:], p.Seq)
	binary.BigEndian.PutInt32(buf[_ackOffset:], p.Ack)
	copy(buf[_heartOffset:], p.Body)
	return buf
}
