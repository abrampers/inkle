package logging

import (
	"testing"
	"time"
)

func TestIsEventEqual(t *testing.T) {
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
			b:    EventLog{tcpsource: "50100"},
			want: false,
		},
		{
			a:    EventLog{},
			b:    EventLog{ipdest: "127.0.0.1"},
			want: false,
		},
		{
			a:    EventLog{},
			b:    EventLog{tcpdest: "8000"},
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

	for _, test := range tests {
		if ret := isEventEqual(test.a, test.b); ret != test.want {
			t.Errorf("isEventEqual(a, b) returns '%t' while it should be '%t'", ret, test.want)
		}
	}
}
