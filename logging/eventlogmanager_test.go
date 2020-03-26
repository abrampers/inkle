package logging

import (
	"testing"
	"time"
)

func isEventsEqual(a, b []*EventLog) bool {
	lena, lenb := len(a), len(b)

	if lena != lenb {
		return false
	}

	for i := 0; i < lena; i++ {
		if !isEventEqual(a[i], b[i]) {
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
	tests := []struct {
		idx                        int
		initialevents, finalevents []*EventLog
		want                       bool
	}{
		{
			idx:           -1,
			initialevents: []*EventLog{},
			finalevents:   []*EventLog{},
			want:          false,
		},
		{
			idx:           0,
			initialevents: []*EventLog{},
			finalevents:   []*EventLog{},
			want:          false,
		},
		{
			idx:           1,
			initialevents: []*EventLog{},
			finalevents:   []*EventLog{},
			want:          false,
		},
		{
			idx: -1,
			initialevents: []*EventLog{
				&EventLog{},
			},
			finalevents: []*EventLog{
				&EventLog{},
			},
			want: false,
		},
		{
			idx: 0,
			initialevents: []*EventLog{
				&EventLog{},
			},
			finalevents: []*EventLog{},
			want:        true,
		},
		{
			idx: 1,
			initialevents: []*EventLog{
				&EventLog{},
			},
			finalevents: []*EventLog{
				&EventLog{},
			},
			want: false,
		},
		{
			idx: -1,
			initialevents: []*EventLog{
				&EventLog{},
				&EventLog{servicename: "helloworld.Greeter"},
			},
			finalevents: []*EventLog{
				&EventLog{},
				&EventLog{servicename: "helloworld.Greeter"},
			},
			want: false,
		},
		{
			idx: 0,
			initialevents: []*EventLog{
				&EventLog{},
				&EventLog{servicename: "helloworld.Greeter"},
			},
			finalevents: []*EventLog{
				&EventLog{servicename: "helloworld.Greeter"},
			},
			want: true,
		},
		{
			idx: 1,
			initialevents: []*EventLog{
				&EventLog{},
				&EventLog{servicename: "helloworld.Greeter"},
			},
			finalevents: []*EventLog{
				&EventLog{},
			},
			want: true,
		},
		{
			idx: 2,
			initialevents: []*EventLog{
				&EventLog{},
				&EventLog{servicename: "helloworld.Greeter"},
			},
			finalevents: []*EventLog{
				&EventLog{},
				&EventLog{servicename: "helloworld.Greeter"},
			},
			want: false,
		},
		{
			idx: -1,
			initialevents: []*EventLog{
				&EventLog{},
				&EventLog{servicename: "datetime.Datetime"},
				&EventLog{servicename: "helloworld.Greeter"},
			},
			finalevents: []*EventLog{
				&EventLog{},
				&EventLog{servicename: "datetime.Datetime"},
				&EventLog{servicename: "helloworld.Greeter"},
			},
			want: false,
		},
		{
			idx: 0,
			initialevents: []*EventLog{
				&EventLog{},
				&EventLog{servicename: "datetime.Datetime"},
				&EventLog{servicename: "helloworld.Greeter"},
			},
			finalevents: []*EventLog{
				&EventLog{servicename: "datetime.Datetime"},
				&EventLog{servicename: "helloworld.Greeter"},
			},
			want: true,
		},
		{
			idx: 1,
			initialevents: []*EventLog{
				&EventLog{},
				&EventLog{servicename: "datetime.Datetime"},
				&EventLog{servicename: "helloworld.Greeter"},
			},
			finalevents: []*EventLog{
				&EventLog{},
				&EventLog{servicename: "helloworld.Greeter"},
			},
			want: true,
		},
		{
			idx: 2,
			initialevents: []*EventLog{
				&EventLog{},
				&EventLog{servicename: "datetime.Datetime"},
				&EventLog{servicename: "helloworld.Greeter"},
			},
			finalevents: []*EventLog{
				&EventLog{},
				&EventLog{servicename: "datetime.Datetime"},
			},
			want: true,
		},
		{
			idx: 3,
			initialevents: []*EventLog{
				&EventLog{},
				&EventLog{servicename: "datetime.Datetime"},
				&EventLog{servicename: "helloworld.Greeter"},
			},
			finalevents: []*EventLog{
				&EventLog{},
				&EventLog{servicename: "datetime.Datetime"},
				&EventLog{servicename: "helloworld.Greeter"},
			},
			want: false,
		},
	}

	for i, test := range tests {
		elm := &eventLogManager{events: test.initialevents}
		err := elm.removeEvent(test.idx)
		if test.want == true && err != nil {
			t.Errorf("removeEvent('%d') (testcase %d): returns err='%v', where there should be no error", test.idx, i, err)
		} else if test.want == false && err == nil {
			t.Errorf("removeEvent('%d') (testcase %d): returns no err, where there should be error", test.idx, i)
		} else if !isEventsEqual(elm.events, test.finalevents) {
			t.Errorf("removeEvent('%d') (testcase %d): doesn't remove events as expected", test.idx, i)
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

// func TestCleanupExpiredRequests(t *testing.T) {
// 	tests := []struct {
// 		initialevents []EventLog
// 		finalevents   []EventLog
// 	}{}
//
// 	for _, test := range tests {
// 	}
// }
