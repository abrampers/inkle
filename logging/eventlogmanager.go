package logging

import (
	"time"
)

type EventLogManager interface {
	CreateEvent(timestamp time.Time, servicename string, methodname string, ipsource string, tcpsource uint16, ipdest string, tcpdest uint16)
	InsertResponse(timestamp time.Time, ipsource string, tcpsource uint16, ipdest string, tcpdest uint16, grpcstatuscode string)
	CleanupExpiredRequests()
}

type eventLogManager struct {
	events  []*EventLog
	tticker *time.Ticker
}

func NewEventLogManager(tticker *time.Ticker) EventLogManager {
	return &eventLogManager{tticker: tticker}
}

func (m *eventLogManager) CreateEvent(timestamp time.Time, servicename string, methodname string, ipsource string, tcpsource uint16, ipdest string, tcpdest uint16) {
	e := NewEventLog(timestamp, servicename, methodname, ipsource, tcpsource, ipdest, tcpdest, "Request")
	m.events = append(m.events, e)
}

func (m *eventLogManager) InsertResponse(timestamp time.Time, ipsource string, tcpsource uint16, ipdest string, tcpdest uint16, grpcstatuscode string) {
}

func (m *eventLogManager) getEvent(ipsource string, tcpsource uint16, ipdest string, tcpdest uint16) (event *EventLog, idx int) {
	for i, event := range m.events {
		if event.isMatchingRequest(ipdest, tcpdest) {
			return event, i
		}
	}
	return nil, -1
}

func (m *eventLogManager) removeEvent(idx int) error {
	return nil
}

func (m *eventLogManager) CleanupExpiredRequests() {
	for currtime := range m.tticker.C {
		m.removeExpiredEvents(currtime)
	}
}

func (m *eventLogManager) removeExpiredEvents(currtime time.Time) {
}
