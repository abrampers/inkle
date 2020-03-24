package logging

import (
	"time"
)

type EventLogManager interface {
	CreateEvent(servicename string, methodname string, ipsource string, tcpsource string, ipdest string, tcpdest string)
	InsertResponse(ipsource string, tcpsource string, ipdest string, tcpdest string, grpcstatuscode string)
	CleanupExpiredRequests()
}

type eventLogManager struct {
	events  []EventLog
	timeout time.Duration
}

func NewEventLogManager(timeout time.Duration) EventLogManager {
	return &eventLogManager{timeout: timeout}
}

func (m *eventLogManager) CreateEvent(servicename string, methodname string, ipsource string, tcpsource string, ipdest string, tcpdest string) {
}

func (m *eventLogManager) InsertResponse(ipsource string, tcpsource string, ipdest string, tcpdest string, grpcstatuscode string) {
}

func (m *eventLogManager) CleanupExpiredRequests() {}
