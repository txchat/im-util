package analyze

import "time"

type SendItem struct {
	time time.Time
}

func (s *SendItem) ReLoad(time time.Time) {
	if s.time.After(time) {
		s.time = time
	}
}

func (s *SendItem) GetTime() time.Time {
	return s.time
}

type AckItem struct {
	time time.Time
	mid  int64
}

func (a *AckItem) ReLoad(time time.Time) {
	if a.time.After(time) {
		a.time = time
	}
}

func (a *AckItem) GetTime() time.Time {
	return a.time
}

func (a *AckItem) GetMid() int64 {
	return a.mid
}

type RevItem struct {
	time time.Time
	mid  int64
}

func (a *RevItem) ReLoad(time time.Time) {
	if a.time.After(time) {
		a.time = time
	}
}

func (a *RevItem) GetTime() time.Time {
	return a.time
}

func (a *RevItem) GetMid() int64 {
	return a.mid
}

type Connection struct {
	connId string

	//key:seq
	allSend map[string]*SendItem
	//key:ack
	allAck map[string]*AckItem
	//key:mid
	allRev map[int64]*RevItem
}

func NewConnection(connId string) *Connection {
	return &Connection{
		connId:  connId,
		allSend: make(map[string]*SendItem),
		allAck:  make(map[string]*AckItem),
		allRev:  make(map[int64]*RevItem),
	}
}

func (c *Connection) LoadSend(seq string, time time.Time) {
	if item, ok := c.allSend[seq]; !ok {
		c.allSend[seq] = &SendItem{
			time: time,
		}
	} else {
		item.ReLoad(time)
	}
}

func (c *Connection) LoadAck(ack string, mid int64, time time.Time) {
	if item, ok := c.allAck[ack]; !ok {
		c.allAck[ack] = &AckItem{
			time: time,
			mid:  mid,
		}
	} else {
		item.ReLoad(time)
	}
}

func (c *Connection) LoadRev(mid int64, time time.Time) {
	if item, ok := c.allRev[mid]; !ok {
		c.allRev[mid] = &RevItem{
			time: time,
			mid:  mid,
		}
	} else {
		item.ReLoad(time)
	}
}
