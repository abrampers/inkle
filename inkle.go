package main

import (
	"flag"
	"fmt"
	"os"
	"strings"
	"time"

	ihttp2 "github.com/abrampers/inkle/http2"
	"github.com/abrampers/inkle/logging"
	"github.com/abrampers/inkle/utils"
	"golang.org/x/net/http2"
)

var (
	isstdout  = flag.Bool("stdout", false, "Write logs to stdout")
	outputdir = flag.String("output", ".", "Output directory of the logs. Ignored if -stdout flag set.")
	timeout   = flag.Duration("timeout", 800*time.Millisecond, "Request timeout in nanosecond")
	err       error
)

const (
	device      string        = "lo0"
	snaplen     int32         = 131072
	promiscuous bool          = false
	itcpTimeout time.Duration = 1000 * time.Millisecond
	filename    string        = "inkle.log"
)

func isGRPC(headers map[string]string) bool {
	for k, _ := range headers {
		if strings.Contains(k, "grpc-") {
			return true
		}
	}
	return false
}

func validateRequestFrameHeaders(headers map[string]string) error {
	method, ok := headers[":method"]
	if !ok {
		return fmt.Errorf("No :method header in frame")
	}
	if method != "POST" {
		return fmt.Errorf(":method is not supported")
	}
	scheme, ok := headers[":scheme"]
	if !ok {
		return fmt.Errorf("No :scheme header in frame")
	}
	if scheme != "http" {
		return fmt.Errorf(":scheme is not supported")
	}
	return nil
}

func validateResponseFrameHeaders(headers map[string]string) error {
	status, ok := headers[":status"]
	if !ok {
		return fmt.Errorf("No :status header in frame")
	}
	if status != "200" {
		return fmt.Errorf("Incorrect status header")
	}
	return nil
}

func requestFrame(h2 ihttp2.HTTP2) (map[string]string, error) {
	for _, frame := range h2.Frames() {
		if frame.Header().Type == http2.FrameHeaders {
			headersframe := frame.(*http2.HeadersFrame)
			headers := ihttp2.Headers(*headersframe)
			_, containsmethod := headers[":method"]
			_, containsscheme := headers[":scheme"]
			_, containspath := headers[":path"]
			_, containsauthority := headers[":authority"]
			_, containsgrpctimeout := headers["grpc-timeout"]
			_, containscontenttype := headers["content-type"]

			if isGRPC(headers) && containsmethod && containsscheme && containspath && containsauthority && containsgrpctimeout && containscontenttype {
				return headers, nil
			}
		}
	}
	return map[string]string{}, fmt.Errorf("No request frame")
}

func responseFrame(h2 ihttp2.HTTP2) (map[string]string, error) {
	for _, frame := range h2.Frames() {
		if frame.Header().Type == http2.FrameHeaders {
			headersframe := frame.(*http2.HeadersFrame)
			headers := ihttp2.Headers(*headersframe)
			_, containsgrpcstatus := headers["grpc-status"]

			if isGRPC(headers) && containsgrpcstatus {
				return headers, nil
			}
		}
	}
	return map[string]string{}, fmt.Errorf("No request frame")
}

func main() {
	flag.Parse()

	interceptor := NewPacketInterceptor(device, snaplen, promiscuous, itcpTimeout)
	defer interceptor.Close()

	var f *os.File
	if !*isstdout {
		fulldir := fmt.Sprintf("%s/%s", *outputdir, filename)
		f, err = os.OpenFile(fulldir, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			panic(err)
		}
	} else {
		f = os.Stdout
	}

	elm := logging.NewEventLogManager(*timeout, f)
	defer elm.Stop()

	go elm.CleanupExpiredRequests()

	for packet := range interceptor.Packets() {
		// Check whether this request is response or not
		if requestheaders, err := requestFrame(packet.HTTP2); err == nil {
			servicename, methodname, err := utils.ParseGrpcPath(requestheaders[":path"])
			if err != nil {
				continue
			}
			elm.CreateEvent(time.Now(), servicename, methodname, packet.SrcIP.String(), uint16(packet.SrcTCP), packet.DstIP.String(), uint16(packet.DstTCP))
		} else if responseheaders, err := responseFrame(packet.HTTP2); err == nil {
			statuscode, ok := responseheaders["grpc-status"]
			if !ok {
				statuscode = "-1"
			}
			elm.InsertResponse(time.Now(), packet.SrcIP.String(), uint16(packet.SrcTCP), packet.DstIP.String(), uint16(packet.DstTCP), statuscode)
		}
	}
}
