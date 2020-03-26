package logging

import (
	"time"
)

type EventLog struct {
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

func (e *EventLog) InsertResponse(timestamp time.Time, grpcstatuscode string, responseinfo string) {
	e.tfinish = timestamp
	e.grpcstatuscode = grpcstatuscode
	e.duration = e.tfinish.Sub(e.tstart)
	e.info += responseinfo
}
