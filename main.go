package main

import (
	"flag"
	"fmt"
	"time"

	"github.com/abrampers/inkle/intercept"
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
	defer interceptor.Close()

	elm := &log.NewEventLogManager{timeout: timeout}

	go elm.CleanupEvent()

	for packet := range interceptor.Packets() {
		// if packet.IsIPv4 {
		// 	fmt.Println("IPv4 SrcIP:        ", packet.IPv4.SrcIP)
		// 	fmt.Println("IPv4 DstIP:        ", packet.IPv4.DstIP)
		// } else {
		// 	fmt.Println("IPv6 SrcIP:        ", packet.IPv6.SrcIP)
		// 	fmt.Println("IPv6 DstIP:        ", packet.IPv6.DstIP)
		// }
		// fmt.Println("TCP srcPort:       ", packet.TCP.SrcPort)
		// fmt.Println("TCP dstPort:       ", packet.TCP.DstPort)
		// fmt.Println("HTTP/2:            ", packet.HTTP2.Frame)

		// Check whether this request is response or not
		if request {
			elm.CreateRequest()
		} else { // Response
			elm.InsertResponse()
		}
	}
}
