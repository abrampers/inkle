package http2

import (
	"reflect"
	"testing"
)

func TestHeadersState(t *testing.T) {
	tests := []struct {
		srcip, dstip   string
		srctcp, dsttcp uint16
		initialstate   map[ipTcpConn]map[string]string
		want           map[string]string
	}{
		{
			srcip:        "::1",
			srctcp:       58000,
			dstip:        "::1",
			dsttcp:       8000,
			initialstate: map[ipTcpConn]map[string]string{},
			want:         map[string]string{},
		},
		{
			srcip:  "::1",
			srctcp: 58000,
			dstip:  "::1",
			dsttcp: 8000,
			initialstate: map[ipTcpConn]map[string]string{
				ipTcpConn{}: map[string]string{"hello": "aloha"},
			},
			want: map[string]string{},
		},
		{
			srcip:  "::1",
			srctcp: 58000,
			dstip:  "::1",
			dsttcp: 8000,
			initialstate: map[ipTcpConn]map[string]string{
				ipTcpConn{}:                          map[string]string{"hello": "aloha"},
				ipTcpConn{"::1", 58001, "::1", 8000}: map[string]string{"hello": "aloha"},
			},
			want: map[string]string{},
		},
		{
			srcip:  "::1",
			srctcp: 58000,
			dstip:  "::1",
			dsttcp: 8000,
			initialstate: map[ipTcpConn]map[string]string{
				ipTcpConn{"::1", 58000, "::1", 8000}: map[string]string{"hello": "aloha"},
			},
			want: map[string]string{"hello": "aloha"},
		},
	}

	for i, test := range tests {
		State := &HeadersState{state: test.initialstate}
		if ret := State.Headers(test.srcip, test.srctcp, test.dstip, test.dsttcp); !reflect.DeepEqual(ret, test.want) {
			t.Errorf("State.Headers (testcase %d): returns incorrect map", i)
		}
	}
}

func TestGetHeaderState(t *testing.T) {
	tests := []struct {
		srcip, dstip             string
		srctcp, dsttcp           uint16
		key, value               string
		initialstate, finalstate map[ipTcpConn]map[string]string
	}{
		{
			srcip:        "::1",
			srctcp:       58000,
			dstip:        "::1",
			dsttcp:       8000,
			key:          ":method",
			value:        "POST",
			initialstate: map[ipTcpConn]map[string]string{},
			finalstate: map[ipTcpConn]map[string]string{
				ipTcpConn{"::1", 58000, "::1", 8000}: map[string]string{
					":method": "POST",
				},
			},
		},
		{
			srcip:  "::1",
			srctcp: 58000,
			dstip:  "::1",
			dsttcp: 8000,
			key:    ":method",
			value:  "POST",
			initialstate: map[ipTcpConn]map[string]string{
				ipTcpConn{}: map[string]string{"hello": "aloha"},
			},
			finalstate: map[ipTcpConn]map[string]string{
				ipTcpConn{}: map[string]string{"hello": "aloha"},
				ipTcpConn{"::1", 58000, "::1", 8000}: map[string]string{
					":method": "POST",
				},
			},
		},
		{
			srcip:  "::1",
			srctcp: 58000,
			dstip:  "::1",
			dsttcp: 8000,
			key:    ":path",
			value:  "/helloworld.Greeter/SayHello",
			initialstate: map[ipTcpConn]map[string]string{
				ipTcpConn{}:                          map[string]string{"hello": "aloha"},
				ipTcpConn{"::1", 58001, "::1", 8000}: map[string]string{"hello": "aloha"},
				ipTcpConn{"::1", 58000, "::1", 8000}: map[string]string{
					":method": "POST",
				},
			},
			finalstate: map[ipTcpConn]map[string]string{
				ipTcpConn{}:                          map[string]string{"hello": "aloha"},
				ipTcpConn{"::1", 58001, "::1", 8000}: map[string]string{"hello": "aloha"},
				ipTcpConn{"::1", 58000, "::1", 8000}: map[string]string{
					":method": "POST",
					":path":   "/helloworld.Greeter/SayHello",
				},
			},
		},
		{
			srcip:  "::1",
			srctcp: 58000,
			dstip:  "::1",
			dsttcp: 8000,
			key:    ":path",
			value:  "/helloworld.Greeter/SayHello",
			initialstate: map[ipTcpConn]map[string]string{
				ipTcpConn{}:                          map[string]string{"hello": "aloha"},
				ipTcpConn{"::1", 58001, "::1", 8000}: map[string]string{"hello": "aloha"},
				ipTcpConn{"::1", 58000, "::1", 8000}: map[string]string{
					":method": "POST",
					":path":   "/helloworld.Greeter/SayHello",
				},
			},
			finalstate: map[ipTcpConn]map[string]string{
				ipTcpConn{}:                          map[string]string{"hello": "aloha"},
				ipTcpConn{"::1", 58001, "::1", 8000}: map[string]string{"hello": "aloha"},
				ipTcpConn{"::1", 58000, "::1", 8000}: map[string]string{
					":method": "POST",
					":path":   "/helloworld.Greeter/SayHello",
				},
			},
		},
		{
			srcip:  "::1",
			srctcp: 58000,
			dstip:  "::1",
			dsttcp: 8000,
			key:    ":path",
			value:  "/helloworld.Greeter/SayHello",
			initialstate: map[ipTcpConn]map[string]string{
				ipTcpConn{}:                          map[string]string{"hello": "aloha"},
				ipTcpConn{"::1", 58001, "::1", 8000}: map[string]string{"hello": "aloha"},
				ipTcpConn{"::1", 58000, "::1", 8000}: map[string]string{
					":method": "POST",
					":path":   "/datetime.Datetime/GetDatetime",
				},
			},
			finalstate: map[ipTcpConn]map[string]string{
				ipTcpConn{}:                          map[string]string{"hello": "aloha"},
				ipTcpConn{"::1", 58001, "::1", 8000}: map[string]string{"hello": "aloha"},
				ipTcpConn{"::1", 58000, "::1", 8000}: map[string]string{
					":method": "POST",
					":path":   "/helloworld.Greeter/SayHello",
				},
			},
		},
	}

	for i, test := range tests {
		State := &HeadersState{state: test.initialstate}
		State.SetHeaders(test.srcip, test.srctcp, test.dstip, test.dsttcp, test.key, test.value)
		if !reflect.DeepEqual(State.state, test.finalstate) {
			t.Errorf("State.SetHeaders (testcase %d): doesn't mutate state as expected", i)
		}
	}
}

func TestUpdateState(t *testing.T) {
	tests := []struct {
		srcip, dstip             string
		srctcp, dsttcp           uint16
		input                    map[string]string
		initialstate, finalstate map[ipTcpConn]map[string]string
	}{
		{
			srcip:        "::1",
			srctcp:       58000,
			dstip:        "::1",
			dsttcp:       8000,
			input:        map[string]string{":method": "POST"},
			initialstate: map[ipTcpConn]map[string]string{},
			finalstate: map[ipTcpConn]map[string]string{
				ipTcpConn{"::1", 58000, "::1", 8000}: map[string]string{
					":method": "POST",
				},
			},
		},
		{
			srcip:  "::1",
			srctcp: 58000,
			dstip:  "::1",
			dsttcp: 8000,
			input:  map[string]string{":method": "POST"},
			initialstate: map[ipTcpConn]map[string]string{
				ipTcpConn{}: map[string]string{"hello": "aloha"},
			},
			finalstate: map[ipTcpConn]map[string]string{
				ipTcpConn{}: map[string]string{"hello": "aloha"},
				ipTcpConn{"::1", 58000, "::1", 8000}: map[string]string{
					":method": "POST",
				},
			},
		},
		{
			srcip:  "::1",
			srctcp: 58000,
			dstip:  "::1",
			dsttcp: 8000,
			input:  map[string]string{":path": "/helloworld.Greeter/SayHello"},
			initialstate: map[ipTcpConn]map[string]string{
				ipTcpConn{}:                          map[string]string{"hello": "aloha"},
				ipTcpConn{"::1", 58001, "::1", 8000}: map[string]string{"hello": "aloha"},
				ipTcpConn{"::1", 58000, "::1", 8000}: map[string]string{
					":method": "POST",
				},
			},
			finalstate: map[ipTcpConn]map[string]string{
				ipTcpConn{}:                          map[string]string{"hello": "aloha"},
				ipTcpConn{"::1", 58001, "::1", 8000}: map[string]string{"hello": "aloha"},
				ipTcpConn{"::1", 58000, "::1", 8000}: map[string]string{
					":method": "POST",
					":path":   "/helloworld.Greeter/SayHello",
				},
			},
		},
		{
			srcip:  "::1",
			srctcp: 58000,
			dstip:  "::1",
			dsttcp: 8000,
			input:  map[string]string{},
			initialstate: map[ipTcpConn]map[string]string{
				ipTcpConn{}:                          map[string]string{"hello": "aloha"},
				ipTcpConn{"::1", 58001, "::1", 8000}: map[string]string{"hello": "aloha"},
				ipTcpConn{"::1", 58000, "::1", 8000}: map[string]string{
					":method": "POST",
					":path":   "/helloworld.Greeter/SayHello",
				},
			},
			finalstate: map[ipTcpConn]map[string]string{
				ipTcpConn{}:                          map[string]string{"hello": "aloha"},
				ipTcpConn{"::1", 58001, "::1", 8000}: map[string]string{"hello": "aloha"},
				ipTcpConn{"::1", 58000, "::1", 8000}: map[string]string{
					":method": "POST",
					":path":   "/helloworld.Greeter/SayHello",
				},
			},
		},
		{
			srcip:  "::1",
			srctcp: 58000,
			dstip:  "::1",
			dsttcp: 8000,
			input:  map[string]string{":path": "/helloworld.Greeter/SayHello"},
			initialstate: map[ipTcpConn]map[string]string{
				ipTcpConn{}:                          map[string]string{"hello": "aloha"},
				ipTcpConn{"::1", 58001, "::1", 8000}: map[string]string{"hello": "aloha"},
				ipTcpConn{"::1", 58000, "::1", 8000}: map[string]string{
					":method": "POST",
					":path":   "/helloworld.Greeter/SayHello",
				},
			},
			finalstate: map[ipTcpConn]map[string]string{
				ipTcpConn{}:                          map[string]string{"hello": "aloha"},
				ipTcpConn{"::1", 58001, "::1", 8000}: map[string]string{"hello": "aloha"},
				ipTcpConn{"::1", 58000, "::1", 8000}: map[string]string{
					":method": "POST",
					":path":   "/helloworld.Greeter/SayHello",
				},
			},
		},
		{
			srcip:  "::1",
			srctcp: 58000,
			dstip:  "::1",
			dsttcp: 8000,
			input:  map[string]string{":path": "/helloworld.Greeter/SayHello"},
			initialstate: map[ipTcpConn]map[string]string{
				ipTcpConn{}:                          map[string]string{"hello": "aloha"},
				ipTcpConn{"::1", 58001, "::1", 8000}: map[string]string{"hello": "aloha"},
				ipTcpConn{"::1", 58000, "::1", 8000}: map[string]string{
					":method": "POST",
					":path":   "/datetime.Datetime/GetDatetime",
				},
			},
			finalstate: map[ipTcpConn]map[string]string{
				ipTcpConn{}:                          map[string]string{"hello": "aloha"},
				ipTcpConn{"::1", 58001, "::1", 8000}: map[string]string{"hello": "aloha"},
				ipTcpConn{"::1", 58000, "::1", 8000}: map[string]string{
					":method": "POST",
					":path":   "/helloworld.Greeter/SayHello",
				},
			},
		},
	}

	for i, test := range tests {
		State := &HeadersState{state: test.initialstate}
		if State.UpdateState(test.srcip, test.srctcp, test.dstip, test.dsttcp, test.input); !reflect.DeepEqual(State.state, test.finalstate) {
			t.Errorf("State.UpdateState (testcase %d): didn't change state as expected", i)
			t.Log(State.state)
			t.Log(test.finalstate)
		}
	}
}
