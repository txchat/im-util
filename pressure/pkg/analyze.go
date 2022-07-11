package pkg

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/txchat/dtalk/pkg/util"
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
	indexMid:       make(map[int64]string),
	tempRev:        make(map[int64]map[string]interface{}),
	tempSource:     make(map[int64]string),
}

type TransmitMsgStatic struct {
	//key=connId+seq
	allTransmitMsg map[string]*TransmitMsg
	//key=mid val=connId+seq
	indexMid map[int64]string
	//key=mid val=line
	tempRev    map[int64]map[string]interface{}
	tempSource map[int64]string
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

func (ts *TransmitMsgStatic) SetMidIndex(mid int64, connId, seq string) {
	if _, ok := ts.indexMid[mid]; !ok {
		ts.indexMid[mid] = keyConnSeq(connId, seq)
	}
}

func (ts *TransmitMsgStatic) GetTransmitMsgByMidIndex(mid int64) *TransmitMsg {
	index := ts.indexMid[mid]
	return ts.allTransmitMsg[index]
}

func (ts *TransmitMsgStatic) SetTempRev(mid int64, data map[string]interface{}, source string) {
	ts.tempRev[mid] = data
	ts.tempSource[mid] = source
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

//
type analyze struct {
	lines  []string
	failed []string
}

func NewAnalyze(lines []string) *analyze {
	return &analyze{
		lines:  lines,
		failed: make([]string, 0),
	}
}

func (a *analyze) Start() error {
	for i, line := range a.lines {
		if line == "" {
			continue
		}
		item := make(map[string]interface{})
		err := json.Unmarshal([]byte(line), &item)
		if err != nil {
			return fmt.Errorf("line %d unmarshal faild: %v", i, err)
		}

		switch item[actionKey] {
		case actionSend:
			err = a.send(item)
		case actionAck:
			err = a.ack(item)
		case actionReceive:
			err = a.receive(item, line)
		}

		if err != nil {
			a.failed = append(a.failed, line)
		}
	}

	//将缓存的内容执行完毕
	ts := GetTransmitMsgStatic()
	for mid, m := range ts.tempRev {
		line := ts.tempSource[mid]
		err := a.receive(m, line)
		if err != nil {
			a.failed = append(a.failed, line)
		}
	}
	return nil
}
func (a *analyze) FailedCount() int {
	return len(a.failed)
}

func (a *analyze) send(item map[string]interface{}) error {
	userId := item[userIdKey].(string)
	connId := item[connIdKey].(string)
	seq := util.ToString(item[seqKey])
	timestr := item[timeKey].(string)
	sendTime, err := parseTime(timestr)
	if err != nil {
		return err
	}

	ts := GetTransmitMsgStatic()
	tm := ts.GetTransmitMsgByConnIdSeq(connId, seq)
	err = tm.LoadSend(userId, connId, seq, sendTime)
	if err != nil {
		return err
	}
	return nil
}

func (a *analyze) ack(item map[string]interface{}) error {
	connId := item[connIdKey].(string)
	ack := util.ToString(item[ackKey])
	mid := util.ToInt64(item[midKey])
	timestr := item[timeKey].(string)
	ackTime, err := parseTime(timestr)
	if err != nil {
		return err
	}

	ts := GetTransmitMsgStatic()
	tm := ts.GetTransmitMsgByConnIdSeq(connId, ack)
	err = tm.LoadAck(connId, ack, mid, ackTime)
	if err != nil {
		return err
	}
	ts.SetMidIndex(mid, connId, ack)
	return nil
}

func (a *analyze) receive(item map[string]interface{}, line string) error {
	mid := util.ToInt64(item[midKey])
	timestr := item[timeKey].(string)
	revTime, err := parseTime(timestr)
	if err != nil {
		return err
	}

	ts := GetTransmitMsgStatic()
	tm := ts.GetTransmitMsgByMidIndex(mid)
	if tm == nil {
		//cache
		ts.SetTempRev(mid, item, line)
		return nil
	}
	err = tm.LoadReceive(mid, revTime)
	if err != nil {
		return err
	}
	return nil
}

func parseTime(timeStr string) (time.Time, error) {
	return time.Parse(time.RFC3339, timeStr)
}
