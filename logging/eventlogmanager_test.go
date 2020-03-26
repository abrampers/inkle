package logging

import (
	"testing"
)

func TestIsEventsEqual(t *testing.T) {
	tests := []struct {
		a    []EventLog
		b    []EventLog
		want bool
	}{
		{
			a:    []EventLog{},
			b:    []EventLog{},
			want: true,
		},
		{
			a: []EventLog{},
			b: []EventLog{
				EventLog{},
			},
			want: false,
		},
		{
			a: []EventLog{
				EventLog{},
			},
			b:    []EventLog{},
			want: false,
		},
		{
			a: []EventLog{
				EventLog{},
			},
			b: []EventLog{
				EventLog{},
			},
			want: true,
		},
		{
			a: []EventLog{
				EventLog{},
			},
			b: []EventLog{
				EventLog{ipsource: "::1"},
			},
			want: false,
		},
	}

	for _, test := range tests {
		if ret := isEventsEqual(test.a, test.b); ret != test.want {
			t.Errorf("isEventsEqual(a, b) returns '%t' while it should be '%t'", ret, test.want)
		}
	}
}

// func TestCreateEvent(t *testing.T) {
// 	currtime := time.Now()
// 	tests := []struct {
// 		currtime      time.Time
// 		initialevents []EventLog
// 		finalevents   []EventLog
// 	}{
// 		{
// 			initialevents: []EventLog{},
// 			finalevents: []EventLog{
// 				*NewEventLog(currtime, "", "", "", "", "", "", ""),
// 			},
// 		},
// 		{
// 			initialevents: []EventLog{
// 				*NewEventLog(time.Time{}, "", "", "", "", "", "", ""),
// 			},
// 			finalevents: []EventLog{
// 				*NewEventLog(time.Time{}, "", "", "", "", "", "", ""),
// 				*NewEventLog(currtime, "", "", "", "", "", "", ""),
// 			},
// 		},
// 	}
//
// }
//
// func TestInsertResponse(t *testing.T) {
// 	tests := []struct {
// 		initialevents []EventLog
// 		finalevents   []EventLog
// 	}{}
//
// 	for _, test := range tests {
// 	}
// }
//
// func TestCleanupExpiredRequests(t *testing.T) {
// 	tests := []struct {
// 		initialevents []EventLog
// 		finalevents   []EventLog
// 	}{}
//
// 	for _, test := range tests {
// 	}
// }
