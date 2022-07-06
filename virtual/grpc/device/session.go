package device

import xproto "github.com/txchat/imparse/proto"

type Session struct {
	channel xproto.Channel
	target  string
}

func (s *Session) SetChannel(ch xproto.Channel) {
	s.channel = ch
}

func (s *Session) GetChannel() xproto.Channel {
	return s.channel
}

func (s *Session) SetTarget(target string) {
	s.target = target
}

func (s *Session) GetTarget() string {
	return s.target
}
