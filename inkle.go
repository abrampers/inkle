package main

import (
	"flag"
	"fmt"
	"strings"
	"time"

	"github.com/abrampers/inkle/intercept"
	"github.com/abrampers/inkle/logging"
	"github.com/abrampers/inkle/utils"
	"golang.org/x/net/http2"
)

var (
	timeout  = flag.Duration("timeout", 2000*time.Millisecond, "Request timeout in nanosecond")
	isstdout = flag.Bool("stdout", false, "Write logs to stdout")
)

const (
	device      string        = "lo0"
	snaplen     int32         = 131072
	promiscuous bool          = false
	itcpTimeout time.Duration = 1000 * time.Millisecond
)

func main() {
	flag.Parse()

	interceptor := intercept.NewPacketInterceptor(device, snaplen, promiscuous, itcpTimeout)
	elm := logging.NewEventLogManager(*timeout, *isstdout)
	defer elm.Stop()
	defer interceptor.Close()

	go elm.CleanupExpiredRequests()

	for packet := range interceptor.Packets() {
		// Check whether this request is response or not
		if requestheaders, err := requestFrame(packet.HTTP2); err != nil {
			servicename, methodname, err := utils.ParseGrpcPath(requestheaders[":path"])
			if err != nil {
				continue
			}
			elm.CreateEvent(time.Now(), servicename, methodname, packet.SrcIP.String(), uint16(packet.SrcTCP), packet.DstIP.String(), uint16(packet.DstTCP))
		} else if responseheaders, err := responseFrame(packet.HTTP2); err != nil {
			statuscode, ok := responseheaders["grpc-status"]
			if !ok {
				statuscode = "-1"
			}
			elm.InsertResponse(time.Now(), packet.SrcIP.String(), uint16(packet.SrcTCP), packet.DstIP.String(), uint16(packet.DstTCP), statuscode)
		}
	}
}

func isGRPC(headers map[string]string) bool {
	for k, _ := range headers {
		if strings.Contains(k, "grpc-") {
			return true
		}
	}
	return false
}

func requestFrame(h2 intercept.HTTP2) (map[string]string, error) {
	for _, frame := range h2.Frames() {
		if frame.Header().Type == http2.FrameHeaders {
			headersframe := frame.(*http2.HeadersFrame)
			headers := intercept.Headers(*headersframe)
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

func responseFrame(h2 intercept.HTTP2) (map[string]string, error) {
	for _, frame := range h2.Frames() {
		if frame.Header().Type == http2.FrameHeaders {
			headersframe := frame.(*http2.HeadersFrame)
			headers := intercept.Headers(*headersframe)
			_, containsgrpcstatus := headers["grpc-status"]

			if isGRPC(headers) && containsgrpcstatus {
				return headers, nil
			}
		}
	}
	return map[string]string{}, fmt.Errorf("No request frame")
}
