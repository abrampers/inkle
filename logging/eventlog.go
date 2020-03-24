package logging

import (
	"regexp"
	"time"
)

var re *regexp.Regexp = regexp.MustCompile(`/([a-zA-Z\.]+)/([a-zA-Z\.]+)`)

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

func NewEventLog(servicename string, methodname string, srcip string, srctcp string, destip string, desttcp string, info string) *EventLog {
	return &EventLog{
		tstart:      time.Now(),
		servicename: servicename,
		methodname:  methodname,
		ipsource:    srcip,
		tcpsource:   srctcp,
		ipdest:      destip,
		tcpdest:     desttcp,
		info:        info,
	}
}
