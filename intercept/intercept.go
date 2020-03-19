package main

import (
	"flag"
	"fmt"

	"github.com/abrampers/inkle/intercept"
)

var (
	device = flag.String("device", "lo0", "Network device to sniff on")
	// device       string        = "lo0"
	snaplen = flag.Int("snaplen", 1024, "The maximum size to read for each packet")
	// snapshot_len int32         = 1024
	promiscuous = flag.Bool("prom", false, "Whether to put the interface in promiscuous mode")
	// promiscuous  bool          = false
	timeout = flag.Duration("timeout", 300, "Timeout in millisecond")
	// timeout      time.Duration = 900 * time.Millisecond
)

func main() {
	fmt.Println("device ", device)
	fmt.Println("snaplen ", snaplen)
	fmt.Println("promiscuous ", promiscuous)
	fmt.Println("timeout ", timeout)
	// interceptor := NewPacketInterceptor(
}
