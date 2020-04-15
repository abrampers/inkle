package main

import (
	"fmt"
	"testing"
)

func Test_validateRequestFrameHeaders(t *testing.T) {
	tests := []struct {
		input map[string]string
		want  error
	}{
		{
			input: map[string]string{},
			want:  fmt.Errorf("No :method header in frame"),
		},
		{
			input: map[string]string{":method": "GET"},
			want:  fmt.Errorf(":method is not supported"),
		},
		{
			input: map[string]string{":method": "POST"},
			want:  fmt.Errorf("No :scheme header in frame"),
		},
		{
			input: map[string]string{":method": "POST", ":scheme": "https"},
			want:  fmt.Errorf(":scheme is not supported"),
		},
		{
			input: map[string]string{":method": "POST", ":scheme": "http"},
			want:  nil,
		},
	}

	for i, test := range tests {
		err := validateRequestFrameHeaders(test.input)
		if err == nil && test.want == nil {
			continue
		} else if err.Error() != test.want.Error() {
			t.Errorf("validateRequestFrameHeaders (testcase %d): returns incorrect error", i)
		}
	}
}

func Test_validateResponseFrameHeaders(t *testing.T) {
	tests := []struct {
		input map[string]string
		want  error
	}{
		{
			input: map[string]string{},
			want:  fmt.Errorf("No :status header in frame"),
		},
		{
			input: map[string]string{":status": "100"},
			want:  fmt.Errorf("Incorrect status header"),
		},
		{
			input: map[string]string{":status": "200"},
			want:  nil,
		},
	}

	for i, test := range tests {
		err := validateResponseFrameHeaders(test.input)
		if err == nil && test.want == nil {
			continue
		} else if err.Error() != test.want.Error() {
			t.Errorf("validateResponseFrameHeaders (testcase %d): returns incorrect error", i)
		}
	}
}
