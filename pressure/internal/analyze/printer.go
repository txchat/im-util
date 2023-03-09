package analyze

import (
	"github.com/rs/zerolog"
)

type Printer struct {
	ts *TransmitMsgStatic

	successCount int
	failedCount  int
	users        map[string]bool
	log          zerolog.Logger
}

func NewPrinter(ts *TransmitMsgStatic, log zerolog.Logger) *Printer {
	return &Printer{
		ts:    ts,
		users: make(map[string]bool),
		log:   log,
	}
}

func (p *Printer) PrintAllLevel1() {
	for _, tm := range p.ts.allTransmitMsg {
		p.users[tm.from] = true
		state := tm.state == (sendFlag | revFlag | ackFlag)
		if state {
			p.successCount++
		} else {
			p.failedCount++
		}
		p.log.Info().Int64("mid", tm.mid).
			Time("sendTime", tm.sendTime).
			Time("receiveTime", tm.receiveTime).
			Time("responseTime", tm.responseTime).
			Bool("status", state).
			Int("flags", tm.state).
			Msg("")
	}
}

func (p *Printer) GetSuccessCount() int {
	return p.successCount
}

func (p *Printer) GetFailedCount() int {
	return p.failedCount
}
