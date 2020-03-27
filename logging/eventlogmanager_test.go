package logging

import (
	"github.com/google/uuid"
	"testing"
	"time"
)

func isEventsEqual(a, b []*EventLog) bool {
	lena, lenb := len(a), len(b)

	if lena != lenb {
		return false
	}

	for i := 0; i < lena; i++ {
		if !isEventEqualValue(a[i], b[i]) {
			return false
		}
	}
	return true
}

func TestIsEventsEqual(t *testing.T) {
	tests := []struct {
		a    []*EventLog
		b    []*EventLog
		want bool
	}{
		{
			a:    []*EventLog{},
			b:    []*EventLog{},
			want: true,
		},
		{
			a: []*EventLog{},
			b: []*EventLog{
				&EventLog{},
			},
			want: false,
		},
		{
			a: []*EventLog{
				&EventLog{},
			},
			b:    []*EventLog{},
			want: false,
		},
		{
			a: []*EventLog{
				&EventLog{},
			},
			b: []*EventLog{
				&EventLog{},
			},
			want: true,
		},
		{
			a: []*EventLog{
				&EventLog{},
			},
			b: []*EventLog{
				&EventLog{ipsource: "::1"},
			},
			want: false,
		},
	}

	for i, test := range tests {
		if ret := isEventsEqual(test.a, test.b); ret != test.want {
			t.Errorf("isEventsEqual(a, b) (testcase %d): returns '%t' while it should be '%t'", i, ret, test.want)
		}
	}
}

func TestAddEvent(t *testing.T) {
	currtime := time.Now()
	tests := []struct {
		event         *EventLog
		initialevents []*EventLog
		finalevents   []*EventLog
	}{
		{
			event: &EventLog{
				tstart:      currtime,
				servicename: "helloworld.Greeter",
				methodname:  "SayHello",
				ipsource:    "::1",
				tcpsource:   58108,
				ipdest:      "::1",
				tcpdest:     8000,
				info:        "Request",
			},
			initialevents: []*EventLog{},
			finalevents: []*EventLog{
				&EventLog{
					tstart:      currtime,
					servicename: "helloworld.Greeter",
					methodname:  "SayHello",
					ipsource:    "::1",
					tcpsource:   58108,
					ipdest:      "::1",
					tcpdest:     8000,
					info:        "Request",
				},
			},
		},
		{
			event: &EventLog{
				tstart:      currtime,
				servicename: "helloworld.Greeter",
				methodname:  "SayHello",
				ipsource:    "::1",
				tcpsource:   58108,
				ipdest:      "::1",
				tcpdest:     8000,
				info:        "Request",
			},
			initialevents: []*EventLog{
				&EventLog{},
			},
			finalevents: []*EventLog{
				&EventLog{},
				&EventLog{
					tstart:      currtime,
					servicename: "helloworld.Greeter",
					methodname:  "SayHello",
					ipsource:    "::1",
					tcpsource:   58108,
					ipdest:      "::1",
					tcpdest:     8000,
					info:        "Request",
				},
			},
		},
	}

	for i, test := range tests {
		elm := &eventLogManager{events: test.initialevents}
		elm.addEvent(test.event)
		if !isEventsEqual(elm.events, test.finalevents) {
			t.Errorf("addEvent (testcase %d): doesn't add event as expected", i)
		}
	}
}

func TestRemoveEvent(t *testing.T) {
	id, id2, id3 := uuid.New(), uuid.New(), uuid.New()
	tests := []struct {
		id                         uuid.UUID
		initialevents, finalevents []*EventLog
	}{
		{
			id:            id,
			initialevents: []*EventLog{},
			finalevents:   []*EventLog{},
		},
		{
			id: id,
			initialevents: []*EventLog{
				&EventLog{id: id},
			},
			finalevents: []*EventLog{},
		},
		{
			id: id,
			initialevents: []*EventLog{
				&EventLog{id: id2},
			},
			finalevents: []*EventLog{
				&EventLog{id: id2},
			},
		},
		{
			id: id,
			initialevents: []*EventLog{
				&EventLog{id: id2},
				&EventLog{id: id3},
			},
			finalevents: []*EventLog{
				&EventLog{id: id2},
				&EventLog{id: id3},
			},
		},
		{
			id: id,
			initialevents: []*EventLog{
				&EventLog{id: id},
				&EventLog{id: id2},
				&EventLog{id: id3},
			},
			finalevents: []*EventLog{
				&EventLog{id: id2},
				&EventLog{id: id3},
			},
		},
		{
			id: id,
			initialevents: []*EventLog{
				&EventLog{id: id2},
				&EventLog{id: id},
				&EventLog{id: id3},
			},
			finalevents: []*EventLog{
				&EventLog{id: id2},
				&EventLog{id: id3},
			},
		},
		{
			id: id,
			initialevents: []*EventLog{
				&EventLog{id: id2},
				&EventLog{id: id3},
				&EventLog{id: id},
			},
			finalevents: []*EventLog{
				&EventLog{id: id2},
				&EventLog{id: id3},
			},
		},
	}

	for i, test := range tests {
		elm := &eventLogManager{events: test.initialevents}
		elm.removeEvent(test.id)
		if !isEventsEqual(elm.events, test.finalevents) {
			t.Errorf("removeEvent('%d') (testcase %d): doesn't remove events as expected", test.id, i)
		}
	}
}

func TestCreateEvent(t *testing.T) {
	currtime := time.Now()
	tests := []struct {
		timestamp     time.Time
		servicename   string
		methodname    string
		ipsource      string
		tcpsource     uint16
		ipdest        string
		tcpdest       uint16
		initialevents []*EventLog
		finalevents   []*EventLog
	}{
		{
			timestamp:     currtime,
			servicename:   "helloworld.Greeter",
			methodname:    "SayHello",
			ipsource:      "::1",
			tcpsource:     58108,
			ipdest:        "::1",
			tcpdest:       8000,
			initialevents: []*EventLog{},
			finalevents: []*EventLog{
				&EventLog{
					tstart:      currtime,
					servicename: "helloworld.Greeter",
					methodname:  "SayHello",
					ipsource:    "::1",
					tcpsource:   58108,
					ipdest:      "::1",
					tcpdest:     8000,
					info:        "Request",
				},
			},
		},
		{
			timestamp:   currtime,
			servicename: "helloworld.Greeter",
			methodname:  "SayHello",
			ipsource:    "::1",
			tcpsource:   58108,
			ipdest:      "::1",
			tcpdest:     8000,
			initialevents: []*EventLog{
				&EventLog{},
			},
			finalevents: []*EventLog{
				&EventLog{},
				&EventLog{
					tstart:      currtime,
					servicename: "helloworld.Greeter",
					methodname:  "SayHello",
					ipsource:    "::1",
					tcpsource:   58108,
					ipdest:      "::1",
					tcpdest:     8000,
					info:        "Request",
				},
			},
		},
	}

	for i, test := range tests {
		elm := &eventLogManager{events: test.initialevents}
		elm.CreateEvent(test.timestamp, test.servicename, test.methodname, test.ipsource, test.tcpsource, test.ipdest, test.tcpdest)
		if !isEventsEqual(elm.events, test.finalevents) {
			t.Errorf("CreateEvent (testcase %d): doesn't create event as expected", i)
		}
	}
}

// func TestInsertResponse(t *testing.T) {
// 	tests := []struct {
// 		initialevents []EventLog
// 		finalevents   []EventLog
// 	}{}
//
// 	for _, test := range tests {
// 	}
// }

func TestGetEvent(t *testing.T) {
	currtime := time.Now()
	tests := []struct {
		ipsource, ipdest   string
		tcpsource, tcpdest uint16
		events             []*EventLog
		idx                int
	}{
		{
			ipdest:  "::1",
			tcpdest: 58108,
			events: []*EventLog{
				&EventLog{
					tstart:      currtime,
					servicename: "helloworld.Greeter",
					methodname:  "SayHello",
					ipsource:    "::1",
					tcpsource:   58108,
					ipdest:      "::1",
					tcpdest:     8000,
					info:        "Request",
				},
			},
			idx: 0,
		},
		{
			ipdest:  "127.0.0.1",
			tcpdest: 58108,
			events: []*EventLog{
				&EventLog{
					tstart:      currtime,
					servicename: "helloworld.Greeter",
					methodname:  "SayHello",
					ipsource:    "::1",
					tcpsource:   58108,
					ipdest:      "::1",
					tcpdest:     8000,
					info:        "Request",
				},
			},
			idx: -1,
		},
		{
			ipdest:  "::1",
			tcpdest: 58110,
			events: []*EventLog{
				&EventLog{
					tstart:      currtime,
					servicename: "helloworld.Greeter",
					methodname:  "SayHello",
					ipsource:    "::1",
					tcpsource:   58108,
					ipdest:      "::1",
					tcpdest:     8000,
					info:        "Request",
				},
			},
			idx: -1,
		},
		{
			ipdest:  "::1",
			tcpdest: 58108,
			events: []*EventLog{
				&EventLog{},
				&EventLog{
					tstart:      currtime,
					servicename: "helloworld.Greeter",
					methodname:  "SayHello",
					ipsource:    "::1",
					tcpsource:   58108,
					ipdest:      "::1",
					tcpdest:     8000,
					info:        "Request",
				},
			},
			idx: 1,
		},
		{
			ipdest:  "127.0.0.1",
			tcpdest: 58108,
			events: []*EventLog{
				&EventLog{},
				&EventLog{
					tstart:      currtime,
					servicename: "helloworld.Greeter",
					methodname:  "SayHello",
					ipsource:    "::1",
					tcpsource:   58108,
					ipdest:      "::1",
					tcpdest:     8000,
					info:        "Request",
				},
			},
			idx: -1,
		},
		{
			ipdest:  "::1",
			tcpdest: 58110,
			events: []*EventLog{
				&EventLog{},
				&EventLog{
					tstart:      currtime,
					servicename: "helloworld.Greeter",
					methodname:  "SayHello",
					ipsource:    "::1",
					tcpsource:   58108,
					ipdest:      "::1",
					tcpdest:     8000,
					info:        "Request",
				},
			},
			idx: -1,
		},
	}

	for i, test := range tests {
		elm := &eventLogManager{events: test.events}
		event, idx := elm.getEvent(test.ipsource, test.tcpsource, test.ipdest, test.tcpdest)
		if idx != test.idx {
			t.Errorf("getEvent (testcase %d): returns incorrect index. Expected '%d' got '%d'.", i, test.idx, idx)
		} else if idx != -1 && event != elm.events[idx] {
			t.Errorf("getEvent (testcase %d): returns incorrect pointer.", i)
		}
	}
}

func TestRemoveExpiredEvents(t *testing.T) {
	currtime := time.Now()
	tests := []struct {
		timeout       time.Duration
		currtime      time.Time
		initialevents []*EventLog
		finalevents   []*EventLog
	}{
		{
			timeout:       100 * time.Millisecond,
			currtime:      currtime,
			initialevents: []*EventLog{},
			finalevents:   []*EventLog{},
		},
		{
			timeout:  100 * time.Millisecond,
			currtime: currtime,
			initialevents: []*EventLog{
				&EventLog{tstart: currtime.Add(-200 * time.Millisecond)},
			},
			finalevents: []*EventLog{},
		},
		{
			timeout:  100 * time.Millisecond,
			currtime: currtime,
			initialevents: []*EventLog{
				&EventLog{tstart: currtime.Add(-200 * time.Millisecond)},
				&EventLog{tstart: currtime.Add(-150 * time.Millisecond)},
				&EventLog{tstart: currtime.Add(-90 * time.Millisecond)},
				&EventLog{tstart: currtime.Add(-80 * time.Millisecond)},
			},
			finalevents: []*EventLog{
				&EventLog{tstart: currtime.Add(-90 * time.Millisecond)},
				&EventLog{tstart: currtime.Add(-80 * time.Millisecond)},
			},
		},
	}

	for i, test := range tests {
		elm := &eventLogManager{events: test.initialevents}
		elm.removeExpiredEvents(test.currtime)
		if !isEventsEqual(elm.events, test.finalevents) {
			t.Errorf("removeExpired (testcase %d): didn't remove expired elements as expected", i)
		}
	}
}
