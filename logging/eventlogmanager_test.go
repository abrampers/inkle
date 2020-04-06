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
		if !isEventEqualValue(*a[i], *b[i]) {
			return false
		}
	}
	return true
}

func Test_isEventsEqual(t *testing.T) {
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

func Test_addEvent(t *testing.T) {
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

func Test_removeEvent(t *testing.T) {
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
					tstart:         currtime,
					servicename:    "helloworld.Greeter",
					methodname:     "SayHello",
					ipsource:       "::1",
					tcpsource:      58108,
					ipdest:         "::1",
					tcpdest:        8000,
					grpcstatuscode: "NULL",
					duration:       0,
					info:           "Request",
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
					tstart:         currtime,
					servicename:    "helloworld.Greeter",
					methodname:     "SayHello",
					ipsource:       "::1",
					tcpsource:      58108,
					ipdest:         "::1",
					tcpdest:        8000,
					grpcstatuscode: "NULL",
					duration:       0,
					info:           "Request",
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

func TestInsertResponse(t *testing.T) {
	currtime := time.Now()
	tests := []struct {
		timestamp                        time.Time
		ipsource, ipdest, grpcstatuscode string
		tcpsource, tcpdest               uint16
		initialevents, finalevents       []*EventLog
		want                             EventLog
	}{
		{
			timestamp:      currtime.Add(50 * time.Millisecond),
			ipsource:       "::1",
			tcpsource:      8000,
			ipdest:         "::1",
			tcpdest:        58108,
			grpcstatuscode: "0",
			initialevents: []*EventLog{
				&EventLog{},
				&EventLog{
					id:          uuid.MustParse("d96763c9-a9a4-49d0-9008-b63befa85b6d"),
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
			finalevents: []*EventLog{
				&EventLog{},
			},
			want: EventLog{
				id:             uuid.MustParse("d96763c9-a9a4-49d0-9008-b63befa85b6d"),
				tstart:         currtime,
				tfinish:        currtime.Add(50 * time.Millisecond),
				servicename:    "helloworld.Greeter",
				methodname:     "SayHello",
				ipsource:       "::1",
				tcpsource:      58108,
				ipdest:         "::1",
				tcpdest:        8000,
				grpcstatuscode: "0",
				duration:       50 * time.Millisecond,
				info:           "Request - Response",
			},
		},
		{
			timestamp:      currtime.Add(50 * time.Millisecond),
			ipsource:       "::1",
			tcpsource:      8000,
			ipdest:         "::1",
			tcpdest:        58108,
			grpcstatuscode: "0",
			initialevents: []*EventLog{
				&EventLog{},
				&EventLog{
					id:          uuid.MustParse("d96763c9-a9a4-49d0-9008-b63befa85b6d"),
					tstart:      currtime,
					servicename: "helloworld.Greeter",
					methodname:  "SayHello",
					ipsource:    "::1",
					tcpsource:   58108,
					ipdest:      "::1",
					tcpdest:     8000,
					info:        "Request",
				},
				&EventLog{
					id:          uuid.MustParse("14a9bb09-23c9-49ad-994c-de1a7f503e12"),
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
			finalevents: []*EventLog{
				&EventLog{},
				&EventLog{
					id:          uuid.MustParse("14a9bb09-23c9-49ad-994c-de1a7f503e12"),
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
			want: EventLog{
				// id:             uuid.MustParse("d96763c9-a9a4-49d0-9008-b63befa85b6d"),
				tstart:         currtime,
				tfinish:        currtime.Add(50 * time.Millisecond),
				servicename:    "helloworld.Greeter",
				methodname:     "SayHello",
				ipsource:       "::1",
				tcpsource:      58108,
				ipdest:         "::1",
				tcpdest:        8000,
				grpcstatuscode: "0",
				duration:       50 * time.Millisecond,
				info:           "Request - Response",
			},
		},
		{
			timestamp:      currtime,
			ipsource:       "::1",
			tcpsource:      8000,
			ipdest:         "::1",
			tcpdest:        58108,
			grpcstatuscode: "0",
			initialevents: []*EventLog{
				&EventLog{},
				&EventLog{},
				&EventLog{},
			},
			finalevents: []*EventLog{
				&EventLog{},
				&EventLog{},
				&EventLog{},
			},
			want: EventLog{
				tstart:         time.Time{},
				tfinish:        currtime,
				servicename:    "NULL",
				methodname:     "NULL",
				ipsource:       "::1",
				tcpsource:      58108,
				ipdest:         "::1",
				tcpdest:        8000,
				grpcstatuscode: "0",
				duration:       currtime.Sub(time.Time{}),
				info:           "NO REQUEST - Response",
			},
		},
	}

	for i, test := range tests {
		elm := &eventLogManager{events: test.initialevents}
		if ret := elm.InsertResponse(test.timestamp, test.ipsource, test.tcpsource, test.ipdest, test.tcpdest, test.grpcstatuscode); !isEventEqualValue(ret, test.want) {
			t.Errorf("InsertResponse (testcase %d): printed incorrect event", i)
		} else if !isEventsEqual(elm.events, test.finalevents) {
			t.Errorf("InsertResponse (testcase %d): doesn't remove event as expected", i)
		}
	}
}

func Test_getEvent(t *testing.T) {
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
		event, idx := elm.getEvent(test.ipdest, test.tcpdest)
		if idx != test.idx {
			t.Errorf("getEvent (testcase %d): returns incorrect index. Expected '%d' got '%d'.", i, test.idx, idx)
		} else if idx != -1 && event != elm.events[idx] {
			t.Errorf("getEvent (testcase %d): returns incorrect pointer.", i)
		}
	}
}

func Test_expiredEvents(t *testing.T) {
	currtime := time.Now()
	tests := []struct {
		timeout  time.Duration
		currtime time.Time
		events   []*EventLog
		want     []*EventLog
	}{
		{
			timeout:  100 * time.Millisecond,
			currtime: currtime,
			events:   []*EventLog{},
			want:     []*EventLog{},
		},
		{
			timeout:  100 * time.Millisecond,
			currtime: currtime,
			events: []*EventLog{
				&EventLog{
					tstart: currtime.Add(-110 * time.Millisecond),
					info:   "Request",
				},
			},
			want: []*EventLog{
				&EventLog{
					tstart:         currtime.Add(-110 * time.Millisecond),
					tfinish:        currtime,
					grpcstatuscode: "NULL",
					duration:       110 * time.Millisecond,
					info:           "Request - TIMEOUT",
				},
			},
		},
		{
			timeout:  100 * time.Millisecond,
			currtime: currtime,
			events: []*EventLog{
				&EventLog{
					tstart: currtime.Add(-110 * time.Millisecond),
					info:   "Request",
				},
				&EventLog{
					tstart: currtime.Add(-80 * time.Millisecond),
					info:   "Request",
				},
			},
			want: []*EventLog{
				&EventLog{
					tstart:         currtime.Add(-110 * time.Millisecond),
					tfinish:        currtime,
					grpcstatuscode: "NULL",
					duration:       110 * time.Millisecond,
					info:           "Request - TIMEOUT",
				},
			},
		},
		{
			timeout:  100 * time.Millisecond,
			currtime: currtime,
			events: []*EventLog{
				&EventLog{
					tstart: currtime.Add(-80 * time.Millisecond),
					info:   "Request",
				},
				&EventLog{
					tstart: currtime.Add(-110 * time.Millisecond),
					info:   "Request",
				},
			},
			want: []*EventLog{
				&EventLog{
					tstart:         currtime.Add(-110 * time.Millisecond),
					tfinish:        currtime,
					grpcstatuscode: "NULL",
					duration:       110 * time.Millisecond,
					info:           "Request - TIMEOUT",
				},
			},
		},
		{
			timeout:  100 * time.Millisecond,
			currtime: currtime,
			events: []*EventLog{
				&EventLog{
					tstart: currtime.Add(-150 * time.Millisecond),
					info:   "Request",
				},
				&EventLog{
					tstart: currtime.Add(-110 * time.Millisecond),
					info:   "Request",
				},
				&EventLog{
					tstart: currtime.Add(-90 * time.Millisecond),
					info:   "Request",
				},
				&EventLog{
					tstart: currtime.Add(-60 * time.Millisecond),
					info:   "Request",
				},
			},
			want: []*EventLog{
				&EventLog{
					tstart:         currtime.Add(-150 * time.Millisecond),
					tfinish:        currtime,
					grpcstatuscode: "NULL",
					duration:       150 * time.Millisecond,
					info:           "Request - TIMEOUT",
				},
				&EventLog{
					tstart:         currtime.Add(-110 * time.Millisecond),
					tfinish:        currtime,
					grpcstatuscode: "NULL",
					duration:       110 * time.Millisecond,
					info:           "Request - TIMEOUT",
				},
			},
		},
	}

	for i, test := range tests {
		elm := &eventLogManager{timeout: test.timeout, events: test.events}
		if ret := elm.expiredEvents(test.currtime); !isEventsEqual(ret, test.want) {
			t.Errorf("expiredEvents('%v') (testcase %d): doesn't return events as expected", test.timeout, i)
		}
	}
}

func Test_removeEvents(t *testing.T) {
	tests := []struct {
		expiredevents, initialevents, finalevents []*EventLog
	}{
		{
			expiredevents: []*EventLog{},
			initialevents: []*EventLog{},
			finalevents:   []*EventLog{},
		},
		{
			expiredevents: []*EventLog{},
			initialevents: []*EventLog{
				&EventLog{id: uuid.MustParse("d96763c9-a9a4-49d0-9008-b63befa85b6d")},
			},
			finalevents: []*EventLog{
				&EventLog{id: uuid.MustParse("d96763c9-a9a4-49d0-9008-b63befa85b6d")},
			},
		},
		{
			expiredevents: []*EventLog{
				&EventLog{id: uuid.MustParse("d96763c9-a9a4-49d0-9008-b63befa85b6d")},
			},
			initialevents: []*EventLog{
				&EventLog{id: uuid.MustParse("d96763c9-a9a4-49d0-9008-b63befa85b6d")},
			},
			finalevents: []*EventLog{},
		},
		{
			expiredevents: []*EventLog{
				&EventLog{id: uuid.MustParse("c8166a03-984f-450c-94bc-3c976f60c6a9")},
				&EventLog{id: uuid.MustParse("14a9bb09-23c9-49ad-994c-de1a7f503e12")},
			},
			initialevents: []*EventLog{
				&EventLog{id: uuid.MustParse("d96763c9-a9a4-49d0-9008-b63befa85b6d")},
				&EventLog{id: uuid.MustParse("e48468ad-c173-48ef-be3e-81e1b3fa48b8")},
				&EventLog{id: uuid.MustParse("c8166a03-984f-450c-94bc-3c976f60c6a9")},
				&EventLog{id: uuid.MustParse("a4baa908-952e-41a8-97e1-89a31e365184")},
				&EventLog{id: uuid.MustParse("14a9bb09-23c9-49ad-994c-de1a7f503e12")},
			},
			finalevents: []*EventLog{
				&EventLog{id: uuid.MustParse("d96763c9-a9a4-49d0-9008-b63befa85b6d")},
				&EventLog{id: uuid.MustParse("e48468ad-c173-48ef-be3e-81e1b3fa48b8")},
				&EventLog{id: uuid.MustParse("a4baa908-952e-41a8-97e1-89a31e365184")},
			},
		},
		{
			expiredevents: []*EventLog{},
			initialevents: []*EventLog{
				&EventLog{id: uuid.MustParse("d96763c9-a9a4-49d0-9008-b63befa85b6d")},
				&EventLog{id: uuid.MustParse("e48468ad-c173-48ef-be3e-81e1b3fa48b8")},
				&EventLog{id: uuid.MustParse("c8166a03-984f-450c-94bc-3c976f60c6a9")},
				&EventLog{id: uuid.MustParse("a4baa908-952e-41a8-97e1-89a31e365184")},
				&EventLog{id: uuid.MustParse("14a9bb09-23c9-49ad-994c-de1a7f503e12")},
			},
			finalevents: []*EventLog{
				&EventLog{id: uuid.MustParse("d96763c9-a9a4-49d0-9008-b63befa85b6d")},
				&EventLog{id: uuid.MustParse("e48468ad-c173-48ef-be3e-81e1b3fa48b8")},
				&EventLog{id: uuid.MustParse("c8166a03-984f-450c-94bc-3c976f60c6a9")},
				&EventLog{id: uuid.MustParse("a4baa908-952e-41a8-97e1-89a31e365184")},
				&EventLog{id: uuid.MustParse("14a9bb09-23c9-49ad-994c-de1a7f503e12")},
			},
		},
		{
			expiredevents: []*EventLog{
				&EventLog{id: uuid.MustParse("d96763c9-a9a4-49d0-9008-b63befa85b6d")},
				&EventLog{id: uuid.MustParse("e48468ad-c173-48ef-be3e-81e1b3fa48b8")},
				&EventLog{id: uuid.MustParse("c8166a03-984f-450c-94bc-3c976f60c6a9")},
				&EventLog{id: uuid.MustParse("a4baa908-952e-41a8-97e1-89a31e365184")},
				&EventLog{id: uuid.MustParse("14a9bb09-23c9-49ad-994c-de1a7f503e12")},
			},
			initialevents: []*EventLog{
				&EventLog{id: uuid.MustParse("d96763c9-a9a4-49d0-9008-b63befa85b6d")},
				&EventLog{id: uuid.MustParse("e48468ad-c173-48ef-be3e-81e1b3fa48b8")},
				&EventLog{id: uuid.MustParse("c8166a03-984f-450c-94bc-3c976f60c6a9")},
				&EventLog{id: uuid.MustParse("a4baa908-952e-41a8-97e1-89a31e365184")},
				&EventLog{id: uuid.MustParse("14a9bb09-23c9-49ad-994c-de1a7f503e12")},
			},
			finalevents: []*EventLog{},
		},
	}

	for i, test := range tests {
		elm := &eventLogManager{events: test.initialevents}
		elm.removeEvents(test.expiredevents)
		if !isEventsEqual(elm.events, test.finalevents) {
			t.Errorf("removeEvents (testcase %d): doesn't remove events as expected", i)
		}
	}
}
