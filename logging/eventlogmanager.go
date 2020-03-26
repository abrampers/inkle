package logging

import (
	"time"
)

type EventLogManager interface {
	CreateEvent(timestamp time.Time, servicename string, methodname string, ipsource string, tcpsource string, ipdest string, tcpdest string)
	InsertResponse(timestamp time.Time, ipsource string, tcpsource string, ipdest string, tcpdest string, grpcstatuscode string)
	CleanupExpiredRequests()
}

type eventLogManager struct {
	events  []*EventLog
	tticker *time.Ticker
}

func NewEventLogManager(tticker *time.Ticker) EventLogManager {
	return &eventLogManager{tticker: tticker}
}

func (m *eventLogManager) CreateEvent(timestamp time.Time, servicename string, methodname string, ipsource string, tcpsource string, ipdest string, tcpdest string) {
	e := NewEventLog(timestamp, servicename, methodname, ipsource, tcpsource, ipdest, tcpdest, "Request")
	m.events = append(m.events, e)
}

func (m *eventLogManager) InsertResponse(timestamp time.Time, ipsource string, tcpsource string, ipdest string, tcpdest string, grpcstatuscode string) {
}

func (m *eventLogManager) CleanupExpiredRequests() {
	for _ = range m.tticker.C {
		m.removeExpiredEvents()
	}
}

func (m *eventLogManager) removeExpiredEvents() {
}
