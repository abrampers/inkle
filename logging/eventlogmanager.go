package logging

import (
	"fmt"
	"sync"
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
	mutex   sync.RWMutex
}

func NewEventLogManager(tticker *time.Ticker) EventLogManager {
	return &eventLogManager{tticker: tticker}
}

func (m *eventLogManager) CreateEvent(timestamp time.Time, servicename string, methodname string, ipsource string, tcpsource uint16, ipdest string, tcpdest uint16) {
	e := NewEventLog(timestamp, servicename, methodname, ipsource, tcpsource, ipdest, tcpdest, "Request")
	m.addEvent(e)
}

func (m *eventLogManager) InsertResponse(timestamp time.Time, ipsource string, tcpsource uint16, ipdest string, tcpdest uint16, grpcstatuscode string) {
}

func (m *eventLogManager) getEvent(ipsource string, tcpsource uint16, ipdest string, tcpdest uint16) (event *EventLog, idx int) {
	i := 0
	for {
		m.mutex.RLock()
		if i >= len(m.events) {
			m.mutex.RUnlock()
			return nil, -1
		} else if event := m.events[i]; event.isMatchingRequest(ipdest, tcpdest) {
			m.mutex.RUnlock()
			return event, i
		} else {
			i += 1
			m.mutex.RUnlock()
		}
	}
}

func (m *eventLogManager) addEvent(event *EventLog) {
	m.mutex.Lock()
	m.events = append(m.events, event)
	m.mutex.Unlock()
}

func (m *eventLogManager) removeEvent(idx int) error {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	lenevents := len(m.events)
	if idx >= lenevents || idx < 0 {
		return fmt.Errorf("Index out of range. idx='%d' while len(events)='%d'", idx, lenevents)
	}
	m.events = append(m.events[:idx], m.events[idx+1:]...)
	return nil
}

func (m *eventLogManager) CleanupExpiredRequests() {
	for currtime := range m.tticker.C {
		m.removeExpiredEvents(currtime)
	}
}

func (m *eventLogManager) removeExpiredEvents(currtime time.Time) {
}
