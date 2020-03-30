package intercept

import (
	"testing"

	"github.com/google/gopacket"
)

func TestDecodeLayers(t *testing.T) {
	tests := []struct {
		input   []byte
		want    bool
		nstream int
	}{
		{
			input: []byte{
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
			want:    true,
			nstream: 2,
		},
		{
			input: []byte{
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
			want:    true,
			nstream: 2,
		},
		{
			input: []byte{
				0x00, 0x00, 0x00, 0x04, 0x00, 0x00, 0x00, 0x00,
				0x00,
			},
			want:    true,
			nstream: 1,
		},
		{
			input: []byte{
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
			want:    true,
			nstream: 5,
		},
		{
			input: []byte{
				0x00, 0x00, 0x08, 0x06, 0x01, 0x00, 0x00, 0x00,
				0x00, 0x02, 0x04, 0x10, 0x10, 0x09, 0x0e, 0x07,
				0x07, 0x00, 0x00, 0x04, 0x08, 0x00, 0x00, 0x00,
				0x00, 0x00, 0x00, 0x00, 0x00, 0x13, 0x00, 0x00,
				0x08, 0x06, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02,
				0x04, 0x10, 0x10, 0x09, 0x0e, 0x07, 0x07,
			},
			want:    true,
			nstream: 3,
		},
		{
			input: []byte{
				0x47, 0x45, 0x54, 0x20, 0x2f, 0x68, 0x65, 0x61,
				0x6c, 0x74, 0x68, 0x79, 0x20, 0x48, 0x54, 0x54,
				0x50, 0x2f, 0x31, 0x2e, 0x31, 0x0d, 0x0a, 0x48,
				0x6f, 0x73, 0x74, 0x3a, 0x20, 0x31, 0x32, 0x37,
				0x2e, 0x30, 0x2e, 0x30, 0x2e, 0x31, 0x3a, 0x35,
				0x36, 0x35, 0x38, 0x31, 0x0d, 0x0a, 0x55, 0x73,
				0x65, 0x72, 0x2d, 0x41, 0x67, 0x65, 0x6e, 0x74,
				0x3a, 0x20, 0x70, 0x79, 0x74, 0x68, 0x6f, 0x6e,
				0x2d, 0x72, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74,
				0x73, 0x2f, 0x32, 0x2e, 0x32, 0x30, 0x2e, 0x31,
				0x0d, 0x0a, 0x41, 0x63, 0x63, 0x65, 0x70, 0x74,
				0x2d, 0x45, 0x6e, 0x63, 0x6f, 0x64, 0x69, 0x6e,
				0x67, 0x3a, 0x20, 0x67, 0x7a, 0x69, 0x70, 0x2c,
				0x20, 0x64, 0x65, 0x66, 0x6c, 0x61, 0x74, 0x65,
				0x0d, 0x0a, 0x41, 0x63, 0x63, 0x65, 0x70, 0x74,
				0x3a, 0x20, 0x2a, 0x2f, 0x2a, 0x0d, 0x0a, 0x43,
				0x6f, 0x6e, 0x6e, 0x65, 0x63, 0x74, 0x69, 0x6f,
				0x6e, 0x3a, 0x20, 0x6b, 0x65, 0x65, 0x70, 0x2d,
				0x61, 0x6c, 0x69, 0x76, 0x65, 0x0d, 0x0a, 0x63,
				0x6f, 0x6e, 0x74, 0x65, 0x6e, 0x74, 0x2d, 0x74,
				0x79, 0x70, 0x65, 0x3a, 0x20, 0x61, 0x70, 0x70,
				0x6c, 0x69, 0x63, 0x61, 0x74, 0x69, 0x6f, 0x6e,
				0x2f, 0x6a, 0x73, 0x6f, 0x6e, 0x0d, 0x0a, 0x78,
				0x2d, 0x79, 0x63, 0x6d, 0x2d, 0x68, 0x6d, 0x61,
				0x63, 0x3a, 0x20, 0x51, 0x50, 0x48, 0x4f, 0x32,
				0x30, 0x6d, 0x53, 0x78, 0x64, 0x71, 0x51, 0x56,
				0x69, 0x55, 0x76, 0x36, 0x51, 0x47, 0x51, 0x43,
				0x2b, 0x7a, 0x6e, 0x53, 0x39, 0x70, 0x4d, 0x63,
				0x4f, 0x4d, 0x2f, 0x6d, 0x4a, 0x44, 0x43, 0x39,
				0x62, 0x30, 0x2b, 0x70, 0x36, 0x4d, 0x3d, 0x0d,
				0x0a, 0x0d, 0x0a,
			},
			want:    false,
			nstream: 0,
		},
		{
			input: []byte{
				0x48, 0x54, 0x54, 0x50, 0x2f, 0x31, 0x2e, 0x31,
				0x20, 0x32, 0x30, 0x30, 0x20, 0x4f, 0x4b, 0x0d,
				0x0a, 0x43, 0x6f, 0x6e, 0x74, 0x65, 0x6e, 0x74,
				0x2d, 0x4c, 0x65, 0x6e, 0x67, 0x74, 0x68, 0x3a,
				0x20, 0x34, 0x0d, 0x0a, 0x43, 0x6f, 0x6e, 0x74,
				0x65, 0x6e, 0x74, 0x2d, 0x54, 0x79, 0x70, 0x65,
				0x3a, 0x20, 0x61, 0x70, 0x70, 0x6c, 0x69, 0x63,
				0x61, 0x74, 0x69, 0x6f, 0x6e, 0x2f, 0x6a, 0x73,
				0x6f, 0x6e, 0x0d, 0x0a, 0x44, 0x61, 0x74, 0x65,
				0x3a, 0x20, 0x4d, 0x6f, 0x6e, 0x2c, 0x20, 0x32,
				0x33, 0x20, 0x4d, 0x61, 0x72, 0x20, 0x32, 0x30,
				0x32, 0x30, 0x20, 0x30, 0x37, 0x3a, 0x30, 0x32,
				0x3a, 0x31, 0x31, 0x20, 0x47, 0x4d, 0x54, 0x0d,
				0x0a, 0x53, 0x65, 0x72, 0x76, 0x65, 0x72, 0x3a,
				0x20, 0x77, 0x61, 0x69, 0x74, 0x72, 0x65, 0x73,
				0x73, 0x0d, 0x0a, 0x58, 0x2d, 0x59, 0x63, 0x6d,
				0x2d, 0x48, 0x6d, 0x61, 0x63, 0x3a, 0x20, 0x34,
				0x65, 0x34, 0x51, 0x67, 0x32, 0x58, 0x32, 0x4d,
				0x53, 0x45, 0x53, 0x58, 0x4f, 0x5a, 0x6f, 0x72,
				0x55, 0x49, 0x7a, 0x79, 0x55, 0x73, 0x30, 0x66,
				0x56, 0x66, 0x6b, 0x72, 0x56, 0x47, 0x41, 0x51,
				0x48, 0x6d, 0x69, 0x7a, 0x54, 0x42, 0x2f, 0x64,
				0x57, 0x77, 0x3d, 0x0d, 0x0a, 0x0d, 0x0a,
			},
			want:    false,
			nstream: 0,
		},
		{
			input: []byte{
				0x50, 0x4f, 0x53, 0x54, 0x20, 0x2f, 0x72, 0x65,
				0x63, 0x65, 0x69, 0x76, 0x65, 0x5f, 0x6d, 0x65,
				0x73, 0x73, 0x61, 0x67, 0x65, 0x73, 0x20, 0x48,
				0x54, 0x54, 0x50, 0x2f, 0x31, 0x2e, 0x31, 0x0d,
				0x0a, 0x48, 0x6f, 0x73, 0x74, 0x3a, 0x20, 0x31,
				0x32, 0x37, 0x2e, 0x30, 0x2e, 0x30, 0x2e, 0x31,
				0x3a, 0x35, 0x36, 0x35, 0x38, 0x31, 0x0d, 0x0a,
				0x55, 0x73, 0x65, 0x72, 0x2d, 0x41, 0x67, 0x65,
				0x6e, 0x74, 0x3a, 0x20, 0x70, 0x79, 0x74, 0x68,
				0x6f, 0x6e, 0x2d, 0x72, 0x65, 0x71, 0x75, 0x65,
				0x73, 0x74, 0x73, 0x2f, 0x32, 0x2e, 0x32, 0x30,
				0x2e, 0x31, 0x0d, 0x0a, 0x41, 0x63, 0x63, 0x65,
				0x70, 0x74, 0x2d, 0x45, 0x6e, 0x63, 0x6f, 0x64,
				0x69, 0x6e, 0x67, 0x3a, 0x20, 0x67, 0x7a, 0x69,
				0x70, 0x2c, 0x20, 0x64, 0x65, 0x66, 0x6c, 0x61,
				0x74, 0x65, 0x0d, 0x0a, 0x41, 0x63, 0x63, 0x65,
				0x70, 0x74, 0x3a, 0x20, 0x2a, 0x2f, 0x2a, 0x0d,
				0x0a, 0x43, 0x6f, 0x6e, 0x6e, 0x65, 0x63, 0x74,
				0x69, 0x6f, 0x6e, 0x3a, 0x20, 0x6b, 0x65, 0x65,
				0x70, 0x2d, 0x61, 0x6c, 0x69, 0x76, 0x65, 0x0d,
				0x0a, 0x63, 0x6f, 0x6e, 0x74, 0x65, 0x6e, 0x74,
				0x2d, 0x74, 0x79, 0x70, 0x65, 0x3a, 0x20, 0x61,
				0x70, 0x70, 0x6c, 0x69, 0x63, 0x61, 0x74, 0x69,
				0x6f, 0x6e, 0x2f, 0x6a, 0x73, 0x6f, 0x6e, 0x0d,
				0x0a, 0x78, 0x2d, 0x79, 0x63, 0x6d, 0x2d, 0x68,
				0x6d, 0x61, 0x63, 0x3a, 0x20, 0x62, 0x6d, 0x46,
				0x69, 0x78, 0x77, 0x75, 0x52, 0x56, 0x45, 0x78,
				0x69, 0x4e, 0x62, 0x35, 0x54, 0x52, 0x52, 0x71,
				0x4c, 0x67, 0x76, 0x42, 0x6e, 0x2b, 0x59, 0x4a,
				0x76, 0x31, 0x6a, 0x73, 0x61, 0x48, 0x72, 0x76,
				0x73, 0x69, 0x4b, 0x7a, 0x77, 0x33, 0x53, 0x73,
				0x3d, 0x0d, 0x0a, 0x43, 0x6f, 0x6e, 0x74, 0x65,
				0x6e, 0x74, 0x2d, 0x4c, 0x65, 0x6e, 0x67, 0x74,
				0x68, 0x3a, 0x20, 0x34, 0x33, 0x31, 0x33, 0x0d,
				0x0a, 0x0d, 0x0a,
			},
			want:    false,
			nstream: 0,
		},
		{
			input: []byte{
				0x48, 0x54, 0x54, 0x50, 0x2f, 0x31, 0x2e, 0x31,
				0x20, 0x32, 0x30, 0x30, 0x20, 0x4f, 0x4b, 0x0d,
				0x0a, 0x43, 0x6f, 0x6e, 0x74, 0x65, 0x6e, 0x74,
				0x2d, 0x4c, 0x65, 0x6e, 0x67, 0x74, 0x68, 0x3a,
				0x20, 0x31, 0x39, 0x32, 0x36, 0x0d, 0x0a, 0x43,
				0x6f, 0x6e, 0x74, 0x65, 0x6e, 0x74, 0x2d, 0x54,
				0x79, 0x70, 0x65, 0x3a, 0x20, 0x61, 0x70, 0x70,
				0x6c, 0x69, 0x63, 0x61, 0x74, 0x69, 0x6f, 0x6e,
				0x2f, 0x6a, 0x73, 0x6f, 0x6e, 0x0d, 0x0a, 0x44,
				0x61, 0x74, 0x65, 0x3a, 0x20, 0x4d, 0x6f, 0x6e,
				0x2c, 0x20, 0x32, 0x33, 0x20, 0x4d, 0x61, 0x72,
				0x20, 0x32, 0x30, 0x32, 0x30, 0x20, 0x30, 0x37,
				0x3a, 0x31, 0x34, 0x3a, 0x30, 0x38, 0x20, 0x47,
				0x4d, 0x54, 0x0d, 0x0a, 0x53, 0x65, 0x72, 0x76,
				0x65, 0x72, 0x3a, 0x20, 0x77, 0x61, 0x69, 0x74,
				0x72, 0x65, 0x73, 0x73, 0x0d, 0x0a, 0x58, 0x2d,
				0x59, 0x63, 0x6d, 0x2d, 0x48, 0x6d, 0x61, 0x63,
				0x3a, 0x20, 0x76, 0x56, 0x71, 0x72, 0x30, 0x44,
				0x32, 0x70, 0x65, 0x55, 0x4f, 0x35, 0x58, 0x37,
				0x30, 0x6c, 0x59, 0x7a, 0x64, 0x57, 0x55, 0x48,
				0x4a, 0x4d, 0x6d, 0x4a, 0x4a, 0x73, 0x4b, 0x41,
				0x79, 0x46, 0x33, 0x50, 0x42, 0x70, 0x65, 0x6b,
				0x4e, 0x69, 0x57, 0x4e, 0x67, 0x3d, 0x0d, 0x0a,
				0x0d, 0x0a,
			},
			want:    false,
			nstream: 0,
		},
	}

	for i, test := range tests {
		h2 := &HTTP2{}
		err := h2.DecodeFromBytes(test.input, gopacket.NilDecodeFeedback)
		if test.want == true && err != nil {
			t.Errorf("DecodeFromBytes('%s') (testcase %d): returns err = '%v', where there shouldn't be no error", string(test.input), i, err)
		} else if test.want == false && err == nil {
			t.Errorf("DecodeFromBytes('%s') (testcase %d): returns no err, where there should be error", string(test.input), i)
		} else if test.nstream != len(h2.Frames()) {
			t.Errorf("DecodeFromBytes('%s') (testcase %d): produces %d where it should be %d stream(s)", string(test.input), i, len(h2.Frames()), test.nstream)
		}
	}
}
