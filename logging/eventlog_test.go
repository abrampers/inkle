package logging

import (
	"testing"
	"time"
)

func isEventEqualValue(a, b EventLog) bool {
	if a.tstart != b.tstart ||
		a.tfinish != b.tfinish ||
		a.servicename != b.servicename ||
		a.methodname != b.methodname ||
		a.ipsource != b.ipsource ||
		a.tcpsource != b.tcpsource ||
		a.ipdest != b.ipdest ||
		a.tcpdest != b.tcpdest ||
		a.grpcstatuscode != b.grpcstatuscode ||
		a.duration != b.duration ||
		a.info != b.info {
		return false
	}
	return true
}

func Test_isEventEqualValue(t *testing.T) {
	tests := []struct {
		a    EventLog
		b    EventLog
		want bool
	}{
		{
			a:    EventLog{},
			b:    EventLog{},
			want: true,
		},
		{
			a:    EventLog{},
			b:    EventLog{tstart: time.Now()},
			want: false,
		},
		{
			a:    EventLog{},
			b:    EventLog{tfinish: time.Now()},
			want: false,
		},
		{
			a:    EventLog{},
			b:    EventLog{servicename: "helloworld.Greeter"},
			want: false,
		},
		{
			a:    EventLog{},
			b:    EventLog{methodname: "SayHello"},
			want: false,
		},
		{
			a:    EventLog{},
			b:    EventLog{ipsource: "::1"},
			want: false,
		},
		{
			a:    EventLog{},
			b:    EventLog{tcpsource: 50100},
			want: false,
		},
		{
			a:    EventLog{},
			b:    EventLog{ipdest: "127.0.0.1"},
			want: false,
		},
		{
			a:    EventLog{},
			b:    EventLog{tcpdest: 8000},
			want: false,
		},
		{
			a:    EventLog{},
			b:    EventLog{grpcstatuscode: "0"},
			want: false,
		},
		{
			a:    EventLog{},
			b:    EventLog{duration: 200 * time.Millisecond},
			want: false,
		},
		{
			a:    EventLog{},
			b:    EventLog{info: "Request - "},
			want: false,
		},
	}

	for i, test := range tests {
		if ret := isEventEqualValue(test.a, test.b); ret != test.want {
			t.Errorf("isEventEqualValue(a, b) (testcase %d): returns '%t' while it should be '%t'", i, ret, test.want)
		}
	}
}

func Test_insertResponse(t *testing.T) {
	stimestamp := time.Now()
	etimestamp := stimestamp.Add(200 * time.Millisecond)

	tests := []struct {
		endtimestamp   time.Time
		grpcstatuscode string
		responseinfo   string
		initialevent   EventLog
		finalevent     EventLog
	}{
		{
			endtimestamp:   etimestamp,
			grpcstatuscode: "0",
			responseinfo:   " - Response",
			initialevent: EventLog{
				tstart:      stimestamp,
				servicename: "helloworld.Greeter",
				methodname:  "SayHello",
				ipsource:    "::1",
				tcpsource:   58108,
				ipdest:      "::1",
				tcpdest:     8000,
				info:        "Request",
			},
			finalevent: EventLog{
				tstart:         stimestamp,
				tfinish:        etimestamp,
				servicename:    "helloworld.Greeter",
				methodname:     "SayHello",
				ipsource:       "::1",
				tcpsource:      58108,
				ipdest:         "::1",
				tcpdest:        8000,
				grpcstatuscode: "0",
				duration:       etimestamp.Sub(stimestamp),
				info:           "Request - Response",
			},
		},
		{
			endtimestamp:   etimestamp,
			grpcstatuscode: "0",
			responseinfo:   " - TIMEOUT",
			initialevent: EventLog{
				tstart:      stimestamp,
				servicename: "helloworld.Greeter",
				methodname:  "SayHello",
				ipsource:    "::1",
				tcpsource:   58108,
				ipdest:      "::1",
				tcpdest:     8000,
				info:        "Request",
			},
			finalevent: EventLog{
				tstart:         stimestamp,
				tfinish:        etimestamp,
				servicename:    "helloworld.Greeter",
				methodname:     "SayHello",
				ipsource:       "::1",
				tcpsource:      58108,
				ipdest:         "::1",
				tcpdest:        8000,
				grpcstatuscode: "0",
				duration:       etimestamp.Sub(stimestamp),
				info:           "Request - TIMEOUT",
			},
		},
	}

	for i, test := range tests {
		test.initialevent.insertResponse(test.endtimestamp, test.grpcstatuscode, test.responseinfo)
		if !isEventEqualValue(test.initialevent, test.finalevent) {
			t.Errorf("insertResponse (testcase %d): doesn't modify event as expected", i)
		}
	}
}

func Test_isMatchingRequest(t *testing.T) {
	tests := []struct {
		ipdest  string
		tcpdest uint16
		event   *EventLog
		want    bool
	}{
		{
			ipdest:  "::1",
			tcpdest: 58108,
			event: &EventLog{
				ipsource:  "::1",
				tcpsource: 58108,
			},
			want: true,
		},
		{
			ipdest:  "::1",
			tcpdest: 58108,
			event: &EventLog{
				ipsource:  "::1",
				tcpsource: 8000,
			},
			want: false,
		},
		{
			ipdest:  "127.0.0.1",
			tcpdest: 58108,
			event: &EventLog{
				ipsource:  "192.168.0.1",
				tcpsource: 58108,
			},
			want: false,
		},
	}

	for i, test := range tests {
		if ret := test.event.isMatchingRequest(test.ipdest, test.tcpdest); ret != test.want {
			t.Errorf("isMatchingRequest('%s', '%d') (testcase %d): expected '%t' got '%t'", test.ipdest, test.tcpdest, i, test.want, ret)
		}
	}
}

func TestString(t *testing.T) {
	stimestamp := time.Date(2000, 2, 1, 12, 13, 14, 0, time.UTC)
	duration := 2 * time.Millisecond
	etimestamp := stimestamp.Add(duration)
	tests := []struct {
		input EventLog
		want  string
	}{
		{
			input: EventLog{
				tstart:         stimestamp,
				tfinish:        etimestamp,
				servicename:    "helloworld.Greeter",
				methodname:     "SayHello",
				ipsource:       "::1",
				tcpsource:      58108,
				ipdest:         "::1",
				tcpdest:        8000,
				grpcstatuscode: "0",
				duration:       duration,
				info:           "Request - TIMEOUT",
			},
			want: "2000-02-01 12:13:14 +0000 UTC, helloworld.Greeter, SayHello, ::1, 58108, ::1, 8000, 0, 2ms, Request - TIMEOUT",
		},
	}

	for i, test := range tests {
		if test.input.String() != test.want {
			t.Errorf("String (testcase %d): returns incorrect string", i)
		}
	}
}
