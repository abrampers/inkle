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

func isEventEqual(a, b EventLog) bool {
	if a.tstart != b.tstart ||
		a.tfinish != b.tfinish ||
		a.servicename != b.servicename ||
		a.methodname != b.methodname ||
		a.ipsource != b.ipsource ||
		a.tcpsource != b.tcpsource ||
		a.ipdest != b.ipdest ||
		a.tcpdest != b.tcpdest ||
		a.grpcstatuscode != b.grpcstatuscode ||
		a.duration != b.duration ||
		a.info != b.info {
		return false
	}
	return true
}
