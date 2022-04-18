package device

import xproto "github.com/txchat/imparse/proto"

var defaultEmpty = NewSession()

type session struct {
	channelType xproto.Channel
	target      string
}

func NewSession() *session {
	return &session{}
}

func (s *session) ChangePage(chType xproto.Channel, target string) {
	s.channelType = chType
	s.target = target
}

func (s *session) GetPage() (xproto.Channel, string) {
	return s.channelType, s.target
}
