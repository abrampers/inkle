package logging

import (
	"fmt"
	"github.com/google/uuid"
	"time"
)

type EventLog struct {
	id             uuid.UUID
	tstart         time.Time
	tfinish        time.Time
	servicename    string
	methodname     string
	ipsource       string
	tcpsource      uint16
	ipdest         string
	tcpdest        uint16
	grpcstatuscode string
	duration       time.Duration
	info           string
}

func NewEventLog(timestamp time.Time, servicename string, methodname string, ipsource string, tcpsource uint16, ipdest string, tcpdest uint16, info string) *EventLog {
	return &EventLog{
		id:          uuid.New(),
		tstart:      timestamp,
		servicename: servicename,
		methodname:  methodname,
		ipsource:    ipsource,
		tcpsource:   tcpsource,
		ipdest:      ipdest,
		tcpdest:     tcpdest,
		info:        info,
	}
}

func (e *EventLog) insertResponse(timestamp time.Time, grpcstatuscode string, responseinfo string) {
	e.tfinish = timestamp
	e.grpcstatuscode = grpcstatuscode
	e.duration = e.tfinish.Sub(e.tstart)
	e.info += responseinfo
}

func (e *EventLog) isMatchingRequest(ipdest string, tcpdest uint16) bool {
	return e.ipsource == ipdest && e.tcpsource == tcpdest
}

func (e EventLog) String() string {
	return fmt.Sprintf("%s, %s, %s, %s, %d, %s, %d, %s, %s, %s", e.tstart, e.servicename, e.methodname, e.ipsource, e.tcpsource, e.ipdest, e.tcpdest, e.grpcstatuscode, e.duration, e.info)
}
