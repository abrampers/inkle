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
	timeout time.Duration
}

func NewEventLogManager(timeout time.Duration) EventLogManager {
	return &eventLogManager{timeout: timeout}
}

func (m *eventLogManager) CreateEvent(timestamp time.Time, servicename string, methodname string, ipsource string, tcpsource string, ipdest string, tcpdest string) {
	e := NewEventLog(timestamp, servicename, methodname, ipsource, tcpsource, ipdest, tcpdest, "Request")
	m.events = append(m.events, e)
}

func (m *eventLogManager) InsertResponse(timestamp time.Time, ipsource string, tcpsource string, ipdest string, tcpdest string, grpcstatuscode string) {
}

func (m *eventLogManager) CleanupExpiredRequests() {}

func isEventsEqual(a, b []*EventLog) bool {
	lena, lenb := len(a), len(b)

	if lena != lenb {
		return false
	}

	for i := 0; i < lena; i++ {
		if !isEventEqual(*a[i], *b[i]) {
			return false
		}
	}
	return true
}
