package main

import (
	"fmt"
	"io/ioutil"
	"net"
	"os"
	"testing"
	"time"

	"github.com/abrampers/inkle/http2"
	"github.com/abrampers/inkle/logging"
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

func Test_outputFile(t *testing.T) {
	file, err := ioutil.TempFile("", "inkle.log")
	if err != nil {
		t.Error("outputFile: failed to create temporary file")
	}

	defer file.Close()
	defer os.Remove("inkle.log")

	tests := []struct {
		isstdout bool
		filename string
		want     *os.File
		err      error
	}{
		{
			isstdout: true,
			want:     os.Stdout,
			err:      nil,
		},
		{
			isstdout: true,
			filename: "asdf/asdf",
			want:     os.Stdout,
			err:      nil,
		},
		{
			isstdout: false,
			filename: file.Name(),
			want:     file,
			err:      nil,
		},
		{
			isstdout: false,
			filename: "asdf/asdf",
			want:     nil,
			err:      fmt.Errorf("open asdf/asdf: no such file or directory"),
		},
	}

	for i, test := range tests {
		f, err := outputFile(test.isstdout, test.filename)
		if err == nil && test.err == nil {
			continue
		} else if (err == nil && test.err != nil) || (err != nil && test.err == nil) || (err.Error() != test.err.Error()) {
			t.Errorf("outputFile (testcase %d): returns incorrect error", i)
			t.Log(err)
			t.Log(test.err)
		} else if f != test.want {
			t.Errorf("outputFile (testcase %d): returns incorrect file", i)
		}
		f.Close()
	}
}

func Test_handlePacket(t *testing.T) {
	tests := []struct {
		bytes []byte
		cidr  *net.IPNet
		want  string
	}{
		{
			bytes: []byte{
				0x00, 0x00, 0x5e, 0x01, 0x04, 0x00, 0x00, 0x00,
				0x01, 0x83, 0x86, 0x45, 0x95, 0x62, 0x72, 0xd1,
				0x41, 0xfc, 0x1e, 0xca, 0x24, 0x5f, 0x15, 0x85,
				0x2a, 0x4b, 0x63, 0x1b, 0x87, 0xeb, 0x19, 0x68,
				0xa0, 0xff, 0x41, 0x8a, 0xa0, 0xe4, 0x1d, 0x13,
				0x9d, 0x09, 0xb8, 0xf0, 0x00, 0x0f, 0x5f, 0x8b,
				0x1d, 0x75, 0xd0, 0x62, 0x0d, 0x26, 0x3d, 0x4c,
				0x4d, 0x65, 0x64, 0x7a, 0x8d, 0x9a, 0xca, 0xc8,
				0xb4, 0xc7, 0x60, 0x2b, 0x89, 0xe5, 0xc0, 0xb4,
				0x85, 0xef, 0x40, 0x02, 0x74, 0x65, 0x86, 0x4d,
				0x83, 0x35, 0x05, 0xb1, 0x1f, 0x40, 0x89, 0x9a,
				0xca, 0xc8, 0xb2, 0x4d, 0x49, 0x4f, 0x6a, 0x7f,
				0x86, 0x7d, 0xf7, 0xdf, 0x71, 0xeb, 0x7f, 0x00,
				0x00, 0x0c, 0x00, 0x01, 0x00, 0x00, 0x00, 0x01,
				0x00, 0x00, 0x00, 0x00, 0x07, 0x0a, 0x05, 0x41,
				0x62, 0x72, 0x61, 0x6d,
			},
			cidr: &net.IPNet{IP: net.ParseIP("::"), Mask: net.CIDRMask(0, 128)},
			want: "helloworld.Greeter,SayHello,::1,58108,::1,8000,-1,0,Request\n",
		},
		{
			bytes: []byte{
				0x00, 0x00, 0x00, 0x04, 0x00, 0x00, 0x00, 0x00,
				0x00,
			},
			cidr: &net.IPNet{IP: net.ParseIP("::"), Mask: net.CIDRMask(0, 128)},
			want: "",
		},
		{
			bytes: []byte{
				0x00, 0x00, 0x04, 0x08, 0x00, 0x00, 0x00, 0x00,
				0x00, 0x00, 0x00, 0x00, 0x0c, 0x00, 0x00, 0x08,
				0x06, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0x04,
				0x10, 0x10, 0x09, 0x0e, 0x07, 0x07, 0x00, 0x00,
				0x0e, 0x01, 0x04, 0x00, 0x00, 0x00, 0x01, 0x88,
				0x5f, 0x8b, 0x1d, 0x75, 0xd0, 0x62, 0x0d, 0x26,
				0x3d, 0x4c, 0x4d, 0x65, 0x64, 0x00, 0x00, 0x12,
				0x00, 0x00, 0x00, 0x00, 0x00, 0x01, 0x00, 0x00,
				0x00, 0x00, 0x0d, 0x0a, 0x0b, 0x48, 0x65, 0x6c,
				0x6c, 0x6f, 0x20, 0x41, 0x62, 0x72, 0x61, 0x6d,
				0x00, 0x00, 0x18, 0x01, 0x05, 0x00, 0x00, 0x00,
				0x01, 0x40, 0x88, 0x9a, 0xca, 0xc8, 0xb2, 0x12,
				0x34, 0xda, 0x8f, 0x01, 0x30, 0x40, 0x89, 0x9a,
				0xca, 0xc8, 0xb5, 0x25, 0x42, 0x07, 0x31, 0x7f,
				0x00,
			},
			cidr: &net.IPNet{IP: net.ParseIP("::"), Mask: net.CIDRMask(0, 128)},
			want: "NULL,NULL,::1,8000,::1,58108,0,0,NO_REQUEST - Response\n",
		},
		{
			bytes: []byte{
				0x00, 0x00, 0x08, 0x06, 0x01, 0x00, 0x00, 0x00,
				0x00, 0x02, 0x04, 0x10, 0x10, 0x09, 0x0e, 0x07,
				0x07, 0x00, 0x00, 0x04, 0x08, 0x00, 0x00, 0x00,
				0x00, 0x00, 0x00, 0x00, 0x00, 0x13, 0x00, 0x00,
				0x08, 0x06, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02,
				0x04, 0x10, 0x10, 0x09, 0x0e, 0x07, 0x07,
			},
			cidr: &net.IPNet{IP: net.ParseIP("::"), Mask: net.CIDRMask(0, 128)},
			want: "",
		},
		{
			bytes: []byte{
				0x00, 0x00, 0x5e, 0x01, 0x04, 0x00, 0x00, 0x00,
				0x01, 0x83, 0x86, 0x45, 0x95, 0x62, 0x72, 0xd1,
				0x41, 0xfc, 0x1e, 0xca, 0x24, 0x5f, 0x15, 0x85,
				0x2a, 0x4b, 0x63, 0x1b, 0x87, 0xeb, 0x19, 0x68,
				0xa0, 0xff, 0x41, 0x8a, 0xa0, 0xe4, 0x1d, 0x13,
				0x9d, 0x09, 0xb8, 0xf0, 0x00, 0x0f, 0x5f, 0x8b,
				0x1d, 0x75, 0xd0, 0x62, 0x0d, 0x26, 0x3d, 0x4c,
				0x4d, 0x65, 0x64, 0x7a, 0x8d, 0x9a, 0xca, 0xc8,
				0xb4, 0xc7, 0x60, 0x2b, 0x89, 0xe5, 0xc0, 0xb4,
				0x85, 0xef, 0x40, 0x02, 0x74, 0x65, 0x86, 0x4d,
				0x83, 0x35, 0x05, 0xb1, 0x1f, 0x40, 0x89, 0x9a,
				0xca, 0xc8, 0xb2, 0x4d, 0x49, 0x4f, 0x6a, 0x7f,
				0x86, 0x7d, 0xf7, 0xdf, 0x71, 0xeb, 0x7f, 0x00,
				0x00, 0x0c, 0x00, 0x01, 0x00, 0x00, 0x00, 0x01,
				0x00, 0x00, 0x00, 0x00, 0x07, 0x0a, 0x05, 0x41,
				0x62, 0x72, 0x61, 0x6d,
			},
			cidr: &net.IPNet{IP: net.ParseIP("::"), Mask: net.CIDRMask(0, 128)},
			want: "helloworld.Greeter,SayHello,::1,58108,::1,8000,-1,0,Request\n",
		},
		{
			bytes: []byte{
				0x00, 0x00, 0x00, 0x04, 0x00, 0x00, 0x00, 0x00,
				0x00,
			},
			cidr: &net.IPNet{IP: net.ParseIP("::"), Mask: net.CIDRMask(0, 128)},
			want: "",
		},
		{
			bytes: []byte{
				0x00, 0x00, 0x04, 0x08, 0x00, 0x00, 0x00, 0x00,
				0x00, 0x00, 0x00, 0x00, 0x0c, 0x00, 0x00, 0x08,
				0x06, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0x04,
				0x10, 0x10, 0x09, 0x0e, 0x07, 0x07, 0x00, 0x00,
				0x0e, 0x01, 0x04, 0x00, 0x00, 0x00, 0x01, 0x88,
				0x5f, 0x8b, 0x1d, 0x75, 0xd0, 0x62, 0x0d, 0x26,
				0x3d, 0x4c, 0x4d, 0x65, 0x64, 0x00, 0x00, 0x12,
				0x00, 0x00, 0x00, 0x00, 0x00, 0x01, 0x00, 0x00,
				0x00, 0x00, 0x0d, 0x0a, 0x0b, 0x48, 0x65, 0x6c,
				0x6c, 0x6f, 0x20, 0x41, 0x62, 0x72, 0x61, 0x6d,
				0x00, 0x00, 0x18, 0x01, 0x05, 0x00, 0x00, 0x00,
				0x01, 0x40, 0x88, 0x9a, 0xca, 0xc8, 0xb2, 0x12,
				0x34, 0xda, 0x8f, 0x01, 0x30, 0x40, 0x89, 0x9a,
				0xca, 0xc8, 0xb5, 0x25, 0x42, 0x07, 0x31, 0x7f,
				0x00,
			},
			cidr: &net.IPNet{IP: net.ParseIP("::"), Mask: net.CIDRMask(0, 128)},
			want: "NULL,NULL,::1,8000,::1,58108,0,0,NO_REQUEST - Response\n",
		},
		{
			bytes: []byte{
				0x00, 0x00, 0x08, 0x06, 0x01, 0x00, 0x00, 0x00,
				0x00, 0x02, 0x04, 0x10, 0x10, 0x09, 0x0e, 0x07,
				0x07, 0x00, 0x00, 0x04, 0x08, 0x00, 0x00, 0x00,
				0x00, 0x00, 0x00, 0x00, 0x00, 0x13, 0x00, 0x00,
				0x08, 0x06, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02,
				0x04, 0x10, 0x10, 0x09, 0x0e, 0x07, 0x07,
			},
			cidr: &net.IPNet{IP: net.ParseIP("::"), Mask: net.CIDRMask(0, 128)},
			want: "",
		},
		{
			bytes: []byte{
				0x00, 0x00, 0x5e, 0x01, 0x04, 0x00, 0x00, 0x00,
				0x01, 0x83, 0x86, 0x45, 0x95, 0x62, 0x72, 0xd1,
				0x41, 0xfc, 0x1e, 0xca, 0x24, 0x5f, 0x15, 0x85,
				0x2a, 0x4b, 0x63, 0x1b, 0x87, 0xeb, 0x19, 0x68,
				0xa0, 0xff, 0x41, 0x8a, 0xa0, 0xe4, 0x1d, 0x13,
				0x9d, 0x09, 0xb8, 0xf0, 0x00, 0x0f, 0x5f, 0x8b,
				0x1d, 0x75, 0xd0, 0x62, 0x0d, 0x26, 0x3d, 0x4c,
				0x4d, 0x65, 0x64, 0x7a, 0x8d, 0x9a, 0xca, 0xc8,
				0xb4, 0xc7, 0x60, 0x2b, 0x89, 0xe5, 0xc0, 0xb4,
				0x85, 0xef, 0x40, 0x02, 0x74, 0x65, 0x86, 0x4d,
				0x83, 0x35, 0x05, 0xb1, 0x1f, 0x40, 0x89, 0x9a,
				0xca, 0xc8, 0xb2, 0x4d, 0x49, 0x4f, 0x6a, 0x7f,
				0x86, 0x7d, 0xf7, 0xdf, 0x71, 0xeb, 0x7f, 0x00,
				0x00, 0x0c, 0x00, 0x01, 0x00, 0x00, 0x00, 0x01,
				0x00, 0x00, 0x00, 0x00, 0x07, 0x0a, 0x05, 0x41,
				0x62, 0x72, 0x61, 0x6d,
			},
			cidr: &net.IPNet{IP: net.ParseIP("::"), Mask: net.CIDRMask(128, 128)},
			want: "helloworld.Greeter,SayHello,::1,58108,::1,8000,-1,0,Request\n",
		},
		{
			bytes: []byte{
				0x00, 0x00, 0x00, 0x04, 0x00, 0x00, 0x00, 0x00,
				0x00,
			},
			cidr: &net.IPNet{IP: net.ParseIP("::"), Mask: net.CIDRMask(128, 128)},
			want: "",
		},
		{
			bytes: []byte{
				0x00, 0x00, 0x04, 0x08, 0x00, 0x00, 0x00, 0x00,
				0x00, 0x00, 0x00, 0x00, 0x0c, 0x00, 0x00, 0x08,
				0x06, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0x04,
				0x10, 0x10, 0x09, 0x0e, 0x07, 0x07, 0x00, 0x00,
				0x0e, 0x01, 0x04, 0x00, 0x00, 0x00, 0x01, 0x88,
				0x5f, 0x8b, 0x1d, 0x75, 0xd0, 0x62, 0x0d, 0x26,
				0x3d, 0x4c, 0x4d, 0x65, 0x64, 0x00, 0x00, 0x12,
				0x00, 0x00, 0x00, 0x00, 0x00, 0x01, 0x00, 0x00,
				0x00, 0x00, 0x0d, 0x0a, 0x0b, 0x48, 0x65, 0x6c,
				0x6c, 0x6f, 0x20, 0x41, 0x62, 0x72, 0x61, 0x6d,
				0x00, 0x00, 0x18, 0x01, 0x05, 0x00, 0x00, 0x00,
				0x01, 0x40, 0x88, 0x9a, 0xca, 0xc8, 0xb2, 0x12,
				0x34, 0xda, 0x8f, 0x01, 0x30, 0x40, 0x89, 0x9a,
				0xca, 0xc8, 0xb5, 0x25, 0x42, 0x07, 0x31, 0x7f,
				0x00,
			},
			cidr: &net.IPNet{IP: net.ParseIP("::"), Mask: net.CIDRMask(128, 128)},
			want: "",
		},
		{
			bytes: []byte{
				0x00, 0x00, 0x08, 0x06, 0x01, 0x00, 0x00, 0x00,
				0x00, 0x02, 0x04, 0x10, 0x10, 0x09, 0x0e, 0x07,
				0x07, 0x00, 0x00, 0x04, 0x08, 0x00, 0x00, 0x00,
				0x00, 0x00, 0x00, 0x00, 0x00, 0x13, 0x00, 0x00,
				0x08, 0x06, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02,
				0x04, 0x10, 0x10, 0x09, 0x0e, 0x07, 0x07,
			},
			cidr: &net.IPNet{IP: net.ParseIP("::"), Mask: net.CIDRMask(128, 128)},
			want: "",
		},
		{
			bytes: []byte{
				0x00, 0x00, 0x5e, 0x01, 0x04, 0x00, 0x00, 0x00,
				0x01, 0x83, 0x86, 0x45, 0x95, 0x62, 0x72, 0xd1,
				0x41, 0xfc, 0x1e, 0xca, 0x24, 0x5f, 0x15, 0x85,
				0x2a, 0x4b, 0x63, 0x1b, 0x87, 0xeb, 0x19, 0x68,
				0xa0, 0xff, 0x41, 0x8a, 0xa0, 0xe4, 0x1d, 0x13,
				0x9d, 0x09, 0xb8, 0xf0, 0x00, 0x0f, 0x5f, 0x8b,
				0x1d, 0x75, 0xd0, 0x62, 0x0d, 0x26, 0x3d, 0x4c,
				0x4d, 0x65, 0x64, 0x7a, 0x8d, 0x9a, 0xca, 0xc8,
				0xb4, 0xc7, 0x60, 0x2b, 0x89, 0xe5, 0xc0, 0xb4,
				0x85, 0xef, 0x40, 0x02, 0x74, 0x65, 0x86, 0x4d,
				0x83, 0x35, 0x05, 0xb1, 0x1f, 0x40, 0x89, 0x9a,
				0xca, 0xc8, 0xb2, 0x4d, 0x49, 0x4f, 0x6a, 0x7f,
				0x86, 0x7d, 0xf7, 0xdf, 0x71, 0xeb, 0x7f, 0x00,
				0x00, 0x0c, 0x00, 0x01, 0x00, 0x00, 0x00, 0x01,
				0x00, 0x00, 0x00, 0x00, 0x07, 0x0a, 0x05, 0x41,
				0x62, 0x72, 0x61, 0x6d,
			},
			cidr: &net.IPNet{IP: net.ParseIP("::"), Mask: net.CIDRMask(128, 128)},
			want: "helloworld.Greeter,SayHello,::1,58108,::1,8000,-1,0,Request\n",
		},
		{
			bytes: []byte{
				0x00, 0x00, 0x00, 0x04, 0x00, 0x00, 0x00, 0x00,
				0x00,
			},
			cidr: &net.IPNet{IP: net.ParseIP("::"), Mask: net.CIDRMask(128, 128)},
			want: "",
		},
		{
			bytes: []byte{
				0x00, 0x00, 0x04, 0x08, 0x00, 0x00, 0x00, 0x00,
				0x00, 0x00, 0x00, 0x00, 0x0c, 0x00, 0x00, 0x08,
				0x06, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0x04,
				0x10, 0x10, 0x09, 0x0e, 0x07, 0x07, 0x00, 0x00,
				0x0e, 0x01, 0x04, 0x00, 0x00, 0x00, 0x01, 0x88,
				0x5f, 0x8b, 0x1d, 0x75, 0xd0, 0x62, 0x0d, 0x26,
				0x3d, 0x4c, 0x4d, 0x65, 0x64, 0x00, 0x00, 0x12,
				0x00, 0x00, 0x00, 0x00, 0x00, 0x01, 0x00, 0x00,
				0x00, 0x00, 0x0d, 0x0a, 0x0b, 0x48, 0x65, 0x6c,
				0x6c, 0x6f, 0x20, 0x41, 0x62, 0x72, 0x61, 0x6d,
				0x00, 0x00, 0x18, 0x01, 0x05, 0x00, 0x00, 0x00,
				0x01, 0x40, 0x88, 0x9a, 0xca, 0xc8, 0xb2, 0x12,
				0x34, 0xda, 0x8f, 0x01, 0x30, 0x40, 0x89, 0x9a,
				0xca, 0xc8, 0xb5, 0x25, 0x42, 0x07, 0x31, 0x7f,
				0x00,
			},
			cidr: &net.IPNet{IP: net.ParseIP("::"), Mask: net.CIDRMask(128, 128)},
			want: "",
		},
		{
			bytes: []byte{
				0x00, 0x00, 0x08, 0x06, 0x01, 0x00, 0x00, 0x00,
				0x00, 0x02, 0x04, 0x10, 0x10, 0x09, 0x0e, 0x07,
				0x07, 0x00, 0x00, 0x04, 0x08, 0x00, 0x00, 0x00,
				0x00, 0x00, 0x00, 0x00, 0x00, 0x13, 0x00, 0x00,
				0x08, 0x06, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02,
				0x04, 0x10, 0x10, 0x09, 0x0e, 0x07, 0x07,
			},
			cidr: &net.IPNet{IP: net.ParseIP("::"), Mask: net.CIDRMask(128, 128)},
			want: "",
		},
	}

	for i, test := range tests {
		h2 := http2.HTTP2{}
		err := h2.DecodeFromBytes(test.bytes, nil)
		if err != nil {
			t.Errorf("handlePacket (testcase %d): wrong test case. Test case should be a valid HTTP/2 bytes", i)
		}
		packet := http2.InterceptedPacket{SrcIP: net.IPv6loopback, DstIP: net.IPv6loopback, SrcTCP: 58108, DstTCP: 8000, HTTP2: h2}
		f, err := ioutil.TempFile("", "Test_printEvent*.log")
		if err != nil {
			t.Errorf("handlePacket (testcase %d): %v", i, err)
		}
		defer f.Close()
		defer os.Remove(f.Name())
		elm := logging.NewEventLogManager(10*time.Millisecond, f, test.cidr)

		if ret := handlePacket(elm, packet); ret != test.want {
			t.Errorf("handlePacket (testcase %d): returns incorrect log line", i)
			t.Log(ret)
			t.Log(test.want)
		}
	}
}
