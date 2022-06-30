package analyze

import (
	"fmt"
	"time"
)

const (
	actionKey = "action"
	userIdKey = "user_id"
	connIdKey = "conn_id"
	seqKey    = "seq"
	ackKey    = "ack"
	midKey    = "mid"
	timeKey   = "time"
)

const (
	actionSend    = "send"
	actionAck     = "ack"
	actionReceive = "receive"
)

const (
	sendFlag = 1 << iota
	ackFlag
	revFlag
)

var transmitMsgStatic = &TransmitMsgStatic{
	allTransmitMsg: make(map[string]*TransmitMsg),
}

type TransmitMsgStatic struct {
	//key=connId+seq
	allTransmitMsg map[string]*TransmitMsg
}

func keyConnSeq(conId, seq string) string {
	return fmt.Sprintf("%s-%s", conId, seq)
}

func GetTransmitMsgStatic() *TransmitMsgStatic {
	return transmitMsgStatic
}

func (ts *TransmitMsgStatic) GetAllTransmitMsgCount() int {
	return len(ts.allTransmitMsg)
}

func (ts *TransmitMsgStatic) GetTransmitMsgByConnIdSeq(connId, seq string) *TransmitMsg {
	var tm *TransmitMsg
	var ok bool
	if tm, ok = ts.allTransmitMsg[keyConnSeq(connId, seq)]; !ok {
		tm = &TransmitMsg{
			connId: connId,
			seq:    seq,
		}
		ts.allTransmitMsg[keyConnSeq(connId, seq)] = tm
	}
	return tm
}

//
type TransmitMsg struct {
	from string

	connId string
	seq    string
	mid    int64

	sendTime     time.Time
	receiveTime  time.Time
	responseTime time.Time

	state int
}

func (tm *TransmitMsg) LoadSend(from, connId, seq string, sendTime time.Time) error {
	if connId != tm.connId || seq != tm.seq {
		return fmt.Errorf("")
	}
	tm.from = from
	tm.sendTime = sendTime
	tm.state |= sendFlag
	return nil
}

func (tm *TransmitMsg) LoadAck(connId, ack string, mid int64, ackTime time.Time) error {
	if connId != tm.connId || ack != tm.seq {
		return fmt.Errorf("err connid seq")
	}
	if tm.mid == 0 {
		tm.mid = mid
		tm.responseTime = ackTime
	} else if tm.responseTime.After(ackTime) {
		tm.responseTime = ackTime
	}
	tm.state |= ackFlag
	return nil
}

func (tm *TransmitMsg) LoadReceive(mid int64, revTime time.Time) error {
	tm.receiveTime = revTime
	tm.state |= revFlag
	return nil
}
