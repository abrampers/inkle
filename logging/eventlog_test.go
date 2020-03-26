package logging

import (
	"testing"
	"time"
)

func isEventEqual(a, b *EventLog) bool {
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

func TestIsEventEqual(t *testing.T) {
	tests := []struct {
		a    *EventLog
		b    *EventLog
		want bool
	}{
		{
			a:    &EventLog{},
			b:    &EventLog{},
			want: true,
		},
		{
			a:    &EventLog{},
			b:    &EventLog{tstart: time.Now()},
			want: false,
		},
		{
			a:    &EventLog{},
			b:    &EventLog{tfinish: time.Now()},
			want: false,
		},
		{
			a:    &EventLog{},
			b:    &EventLog{servicename: "helloworld.Greeter"},
			want: false,
		},
		{
			a:    &EventLog{},
			b:    &EventLog{methodname: "SayHello"},
			want: false,
		},
		{
			a:    &EventLog{},
			b:    &EventLog{ipsource: "::1"},
			want: false,
		},
		{
			a:    &EventLog{},
			b:    &EventLog{tcpsource: "50100"},
			want: false,
		},
		{
			a:    &EventLog{},
			b:    &EventLog{ipdest: "127.0.0.1"},
			want: false,
		},
		{
			a:    &EventLog{},
			b:    &EventLog{tcpdest: "8000"},
			want: false,
		},
		{
			a:    &EventLog{},
			b:    &EventLog{grpcstatuscode: "0"},
			want: false,
		},
		{
			a:    &EventLog{},
			b:    &EventLog{duration: 200 * time.Millisecond},
			want: false,
		},
		{
			a:    &EventLog{},
			b:    &EventLog{info: "Request - "},
			want: false,
		},
	}

	for _, test := range tests {
		if ret := isEventEqual(test.a, test.b); ret != test.want {
			t.Errorf("isEventEqual(a, b) returns '%t' while it should be '%t'", ret, test.want)
		}
	}
}

func TestInsertResponse(t *testing.T) {
	stimestamp := time.Now()
	etimestamp := stimestamp.Add(200 * time.Millisecond)
	tests := []struct {
		endtimestamp   time.Time
		grpcstatuscode string
		responseinfo   string
		initialevent   *EventLog
		finalevent     *EventLog
	}{
		{
			endtimestamp:   etimestamp,
			grpcstatuscode: "0",
			responseinfo:   " - Response",
			initialevent: &EventLog{
				tstart:      stimestamp,
				servicename: "helloworld.Greeter",
				methodname:  "SayHello",
				ipsource:    "::1",
				tcpsource:   "58108",
				ipdest:      "::1",
				tcpdest:     "8000",
				info:        "Request",
			},
			finalevent: &EventLog{
				tstart:         stimestamp,
				tfinish:        etimestamp,
				servicename:    "helloworld.Greeter",
				methodname:     "SayHello",
				ipsource:       "::1",
				tcpsource:      "58108",
				ipdest:         "::1",
				tcpdest:        "8000",
				grpcstatuscode: "0",
				info:           "Request",
			},
		},
	}

	for _, test := range tests {
		test.initialevent.InsertResponse(test.endtimestamp, test.grpcstatuscode, test.responseinfo)
		if !isEventEqual(test.initialevent, test.finalevent) {
			t.Error("InsertResponse is not working properly")
		}
	}
}
