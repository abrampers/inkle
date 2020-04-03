package logging

import (
	"os"
	"sync"
	"time"

	"github.com/google/uuid"
)

type EventLogManager interface {
	CreateEvent(timestamp time.Time, servicename string, methodname string, ipsource string, tcpsource uint16, ipdest string, tcpdest uint16)
	InsertResponse(timestamp time.Time, ipsource string, tcpsource uint16, ipdest string, tcpdest uint16, grpcstatuscode string) EventLog
	CleanupExpiredRequests()
	Stop()
}

type eventLogManager struct {
	events  []*EventLog
	tticker *time.Ticker
	timeout time.Duration
	mutex   sync.RWMutex
	writer  *os.File
}

func NewEventLogManager(timeout time.Duration, isstdout bool) EventLogManager {
	if isstdout {
		f, err := os.OpenFile("text.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			panic(err)
		}
		return &eventLogManager{timeout: timeout, tticker: time.NewTicker(timeout), writer: f}
	}
	return &eventLogManager{timeout: timeout, tticker: time.NewTicker(timeout), writer: os.Stdout}
}

// TODO: Print all remaining events as timeout
func (m *eventLogManager) Stop() {
	m.tticker.Stop()
	defer m.writer.Close()
}

func (m *eventLogManager) CreateEvent(timestamp time.Time, servicename string, methodname string, ipsource string, tcpsource uint16, ipdest string, tcpdest uint16) {
	e := NewEventLog(timestamp, servicename, methodname, ipsource, tcpsource, ipdest, tcpdest, "Request")
	m.addEvent(e)
}

func (m *eventLogManager) InsertResponse(timestamp time.Time, ipsource string, tcpsource uint16, ipdest string, tcpdest uint16, grpcstatuscode string) EventLog {
	var event *EventLog
	var idx int
	event, idx = m.getEvent(ipdest, tcpdest)
	if idx == -1 {
		event = NewEventLog(time.Time{}, "NULL", "NULL", ipdest, tcpdest, ipsource, tcpsource, "NO REQUEST")
	} else {
		m.removeEvent(event.id)
	}

	event.insertResponse(timestamp, grpcstatuscode, " - Response")
	// m.printEvent(event) // Consider spawn goroutine
	return *event
}

func (m *eventLogManager) getEvent(ipdest string, tcpdest uint16) (event *EventLog, idx int) {
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
			event.insertResponse(currtime, "NULL", " - TIMEOUT")
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
