package msggenerator

import (
	"github.com/txchat/im-util/internal/device"
	"github.com/txchat/im-util/internal/user"
	xproto "github.com/txchat/imparse/proto"
	"math/rand"
	"time"
)

type MsgGenerator struct {
	users     []*user.User
	sendClose chan bool
	ackClose  chan bool
}

func NewMsgGenerator(users []*user.User) *MsgGenerator {
	if len(users) < 2 {
		panic("system users less than 2")
	}
	return &MsgGenerator{
		users:     users,
		sendClose: make(chan bool),
		ackClose:  make(chan bool),
	}
}

func (m *MsgGenerator) randomTarget(userClient *user.User) string {
	//将时间戳设置成种子数
	rand.Seed(time.Now().UnixNano())
	index := rand.Intn(len(m.users) - 1)
	if m.users[index].GetUID() == userClient.GetUID() {
		if index >= len(m.users)-1 {
			return m.users[0].GetUID()
		}
		return m.users[index+1].GetUID()
	}
	return m.users[index].GetUID()
}

func (m *MsgGenerator) RangeSend(device *device.Device, rate time.Duration) {
	ticker := time.NewTicker(rate)

	for {
		select {
		case <-m.sendClose:
			return
		case <-ticker.C:
			device.SendMsg(xproto.Channel_name[int32(xproto.Channel_ToUser)], m.randomTarget(device.GetUser()), "1")
		}
	}
}

func (m *MsgGenerator) StopSend() {
	close(m.sendClose)
}

func (m *MsgGenerator) StopAck() {
	close(m.ackClose)
}
