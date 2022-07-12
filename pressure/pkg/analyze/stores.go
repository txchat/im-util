package analyze

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/txchat/dtalk/pkg/util"
)

type tiledRev struct {
	connId map[string]int
}

func NewTiledRev() *tiledRev {
	return &tiledRev{
		connId: make(map[string]int),
	}
}

func (t *tiledRev) LoadConn(connId string) {
	t.connId[connId]++
}

func (t *tiledRev) ExceptConnId(connId string) string {
	for cid := range t.connId {
		if cid != connId {
			return cid
		}
	}
	return ""
}

type AnalyzeStore struct {
	lines  []string
	failed []string

	//key: connId
	connInfo map[string]*Connection
	tileRev  map[int64]*tiledRev
}

func NewAnalyzeStore(lines []string) *AnalyzeStore {
	return &AnalyzeStore{
		lines:    lines,
		failed:   make([]string, 0),
		connInfo: make(map[string]*Connection),
		tileRev:  make(map[int64]*tiledRev),
	}
}

func (t *AnalyzeStore) LoadAll() error {
	for i, line := range t.lines {
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
			err = t.parseSend(item)
		case actionAck:
			err = t.parseAck(item)
		case actionReceive:
			err = t.parseReceive(item)
		}
		if err != nil {
			t.failed = append(t.failed, line)
		}
	}
	t.TileRev()
	return nil
}

func (t *AnalyzeStore) Start() error {
	ts := GetTransmitMsgStatic()
	for connId, connection := range t.connInfo {
		for seq, item := range connection.allSend {
			// send
			tm := ts.GetTransmitMsgByConnIDSeq(connId, seq)
			err := tm.LoadSend("", connId, seq, item.GetTime())
			if err != nil {
				return err
			}
			// find ack
			ack := connection.allAck[seq]
			if ack == nil {
				continue
			}
			err = tm.LoadAck(connId, seq, ack.GetMid(), ack.GetTime())
			if err != nil {
				return err
			}
			// find rev
			rev := t.tileRev[ack.GetMid()]
			if rev == nil {
				continue
			}
			revItem := t.getRevItem(rev.ExceptConnId(connId), ack.GetMid())
			if revItem == nil {
				continue
			}
			err = tm.LoadReceive(revItem.GetMid(), revItem.GetTime())
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func (t *AnalyzeStore) FailedCount() int {
	return len(t.failed)
}

func (t *AnalyzeStore) TileRev() {
	// key:mid, val: connId
	for connId, connection := range t.connInfo {
		for mid := range connection.allRev {
			var xx *tiledRev
			var ok bool
			if xx, ok = t.tileRev[mid]; !ok {
				xx = NewTiledRev()
				t.tileRev[mid] = xx
			}
			xx.LoadConn(connId)
		}
	}
}

func (t *AnalyzeStore) getRevItem(connId string, mid int64) *RevItem {
	c := t.connInfo[connId]
	if c == nil {
		return nil
	}
	return c.allRev[mid]
}

func (t *AnalyzeStore) parseSend(item map[string]interface{}) error {
	//userId := item[userIdKey].(string)
	connId := item[connIdKey].(string)
	seq := util.ToString(item[seqKey])
	timestr := item[timeKey].(string)
	sendTime, err := parseTime(timestr)
	if err != nil {
		return err
	}

	c := t.getConn(connId)
	c.LoadSend(seq, sendTime)
	return nil
}

func (t *AnalyzeStore) parseAck(item map[string]interface{}) error {
	connId := item[connIdKey].(string)
	ack := util.ToString(item[ackKey])
	mid := util.ToInt64(item[midKey])
	timestr := item[timeKey].(string)
	ackTime, err := parseTime(timestr)
	if err != nil {
		return err
	}

	c := t.getConn(connId)
	c.LoadAck(ack, mid, ackTime)
	return nil
}

func (t *AnalyzeStore) parseReceive(item map[string]interface{}) error {
	connId := item[connIdKey].(string)
	mid := util.ToInt64(item[midKey])
	timestr := item[timeKey].(string)
	revTime, err := parseTime(timestr)
	if err != nil {
		return err
	}

	c := t.getConn(connId)
	c.LoadRev(mid, revTime)
	return nil
}

func (t *AnalyzeStore) getConn(connId string) *Connection {
	var c *Connection
	var ok bool
	if c, ok = t.connInfo[connId]; !ok {
		c = NewConnection(connId)
		t.connInfo[connId] = c
	}
	return c
}

func parseTime(timeStr string) (time.Time, error) {
	return time.Parse(time.RFC3339, timeStr)
}
