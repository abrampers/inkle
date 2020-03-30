package main

import (
	"flag"
	"time"

	"github.com/abrampers/inkle/intercept"
	"github.com/abrampers/inkle/logging"
	"github.com/abrampers/inkle/utils"
)

var (
	timeout = flag.Duration("timeout", 2000*time.Millisecond, "Request timeout in nanosecond")
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
	elm := logging.NewEventLogManager(*timeout)
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
			statuscode := responseheaders["grpc-status"]
			elm.InsertResponse(time.Now(), packet.SrcIP.String(), uint16(packet.SrcTCP), packet.DstIP.String(), uint16(packet.DstTCP), statuscode)
		}
	}
}

func requestFrame(frame intercept.HTTP2) (map[string]string, error) {
	return map[string]string{}, nil
}

func responseFrame(frame intercept.HTTP2) (map[string]string, error) {
	return map[string]string{}, nil
}
