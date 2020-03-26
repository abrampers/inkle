package utils

import (
	"testing"
)

func TestParseGrpcPath(t *testing.T) {
	tests := []struct {
		path        string
		servicename string
		methodname  string
		want        bool
	}{
		{
			path:        "/helloworld.Greeter/SayHello",
			servicename: "helloworld.Greeter",
			methodname:  "SayHello",
			want:        true,
		},
		{
			path:        "/google.pubsub.v2.PublisherService/CreateTopic",
			servicename: "google.pubsub.v2.PublisherService",
			methodname:  "CreateTopic",
			want:        true,
		},
		{
			path:        "/Hello/hello/hello",
			servicename: "",
			methodname:  "",
			want:        false,
		},
		{
			path:        "/hello",
			servicename: "",
			methodname:  "",
			want:        false,
		},
	}

	for _, test := range tests {
		servicename, methodname, err := ParseGrpcPath(test.path)
		if test.want == true && err != nil {
			t.Errorf("ParseGrpcPath(%s) returns err = %v, where shouldn't be no error", test.path, err)
		} else if test.want == false && err == nil {
			t.Errorf("ParseGrpcPath(%s) returns no err, where there should be error", string(test.path))
		} else if servicename != test.servicename {
			t.Errorf("ParseGrpcPath(%s) returns servicename='%s', where it should be '%s'", string(test.path), servicename, test.servicename)
		} else if methodname != test.methodname {
			t.Errorf("ParseGrpcPath(%s) returns methodname='%s', where it should be '%s'", string(test.path), methodname, test.methodname)
		}
	}
}
