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
	tcpsource      string
	ipdest         string
	tcpdest        string
	grpcstatuscode string
	duration       time.Duration
	info           string
}

func NewEventLog(timestamp time.Time, servicename string, methodname string, srcip string, srctcp string, destip string, desttcp string, info string) *EventLog {
	return &EventLog{
		tstart:      timestamp,
		servicename: servicename,
		methodname:  methodname,
		ipsource:    srcip,
		tcpsource:   srctcp,
		ipdest:      destip,
		tcpdest:     desttcp,
		info:        info,
	}
}

func (e *EventLog) InsertResponse(timestamp time.Time, grpcstatuscode string, responseinfo string) {
	e.tfinish = timestamp
	e.grpcstatuscode = grpcstatuscode
	e.duration = e.tfinish.Sub(e.tstart)
	e.info += responseinfo
}
