package logging

import (
	"sync"
	"time"

	"github.com/google/uuid"
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
	m.mutex.RLock()
	defer m.mutex.RUnlock()
	i := 0
	for {
		// m.mutex.RLock()
		if i >= len(m.events) {
			// m.mutex.RUnlock()
			return nil, -1
		} else if event := m.events[i]; event.isMatchingRequest(ipdest, tcpdest) {
			// m.mutex.RUnlock()
			return event, i
		} else {
			i += 1
			// m.mutex.RUnlock()
		}
	}
}

func (m *eventLogManager) addEvent(event *EventLog) {
	m.mutex.Lock()
	m.events = append(m.events, event)
	m.mutex.Unlock()
}

func (m *eventLogManager) removeEvent(id uuid.UUID) {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	var idx int
	var event *EventLog
	found := false

	for idx, event = range m.events {
		if event.id == id {
			found = true
			break
		}
	}

	if found {
		m.events = append(m.events[:idx], m.events[idx+1:]...)
	}
}

func (m *eventLogManager) CleanupExpiredRequests() {
	for currtime := range m.tticker.C {
		m.removeExpiredEvents(currtime)
	}
}

func (m *eventLogManager) removeExpiredEvents(currtime time.Time) {
}
