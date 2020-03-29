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
	Stop()
}

type eventLogManager struct {
	events  []*EventLog
	tticker *time.Ticker
	timeout time.Duration
	mutex   sync.RWMutex
}

func NewEventLogManager(timeout time.Duration) EventLogManager {
	return &eventLogManager{timeout: timeout, tticker: time.NewTicker(timeout)}
}

// TODO: Print all remaining events as timeout
func (m *eventLogManager) Stop() {
	m.tticker.Stop()
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
	for i, event := range m.events {
		if event.isMatchingRequest(ipdest, tcpdest) {
			return event, i
		}
	}
	return nil, -1
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
		expiredevents := m.expiredEvents(currtime)
		m.removeEvents(expiredevents)
		m.printEvents(expiredevents)
	}
}

// This should return the events in the same order with events in the array
func (m *eventLogManager) expiredEvents(currtime time.Time) []*EventLog {
	m.mutex.RLock()
	defer m.mutex.RUnlock()
	expiredevents := []*EventLog{}

	for _, event := range m.events {
		if currtime.Sub(event.tstart) >= m.timeout {
			expiredevents = append(expiredevents, event)
		}
	}

	return expiredevents
}

// This should remove the records in order
func (m *eventLogManager) removeEvents(events []*EventLog) {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	idx := 0
	for _, rmevent := range events {
		for idx < len(m.events) {
			if m.events[idx].id == rmevent.id {
				m.events = append(m.events[:idx], m.events[idx+1:]...)
				break
			} else {
				idx++
			}
		}
	}
}

// TODO: Use ELK stack or for testing purposes, write to file
func (m *eventLogManager) printEvents(events []*EventLog) {
}
