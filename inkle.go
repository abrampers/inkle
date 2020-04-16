package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/abrampers/inkle/http2"
	"github.com/abrampers/inkle/logging"
	"github.com/abrampers/inkle/utils"
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

func handlePacket(elm logging.EventLogManager, packet http2.InterceptedPacket) string {
	headers := http2.Headers(packet.HTTP2)
	// Check whether this request is response or not
	if err := validateRequestFrameHeaders(headers); err == nil {
		http2.State.UpdateState(packet.SrcIP.String(), uint16(packet.SrcTCP), packet.DstIP.String(), uint16(packet.DstTCP), headers)
		headers = http2.State.Headers(packet.SrcIP.String(), uint16(packet.SrcTCP), packet.DstIP.String(), uint16(packet.DstTCP))
		servicename, methodname, err := utils.ParseGrpcPath(headers[":path"])
		if err != nil {
			return ""
		}
		return elm.CreateEvent(time.Now(), servicename, methodname, packet.SrcIP.String(), uint16(packet.SrcTCP), packet.DstIP.String(), uint16(packet.DstTCP))
	} else if err := validateResponseFrameHeaders(headers); err == nil {
		http2.State.UpdateState(packet.SrcIP.String(), uint16(packet.SrcTCP), packet.DstIP.String(), uint16(packet.DstTCP), headers)
		headers = http2.State.Headers(packet.SrcIP.String(), uint16(packet.SrcTCP), packet.DstIP.String(), uint16(packet.DstTCP))
		statuscode, ok := headers["grpc-status"]
		if !ok {
			statuscode = "-1"
		}
		return elm.InsertResponse(time.Now(), packet.SrcIP.String(), uint16(packet.SrcTCP), packet.DstIP.String(), uint16(packet.DstTCP), statuscode)
	}
	return ""
}

func outputFile(isstdout bool, filepath string) (*os.File, error) {
	var f *os.File
	if !isstdout {
		f, err = os.OpenFile(filepath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			return nil, err
		}
	} else {
		f = os.Stdout
	}
	return f, nil
}

func main() {
	flag.Parse()
	interceptor := http2.NewPacketInterceptor(device, snaplen, promiscuous, itcpTimeout)
	defer interceptor.Close()
	filepath := filepath.Join(*outputdir, filename)
	f, err := outputFile(*isstdout, filepath)
	if err != nil {
		log.Println("Failed to create output file")
		panic(err)
	}
	elm := logging.NewEventLogManager(*timeout, f)
	defer elm.Stop()

	go elm.CleanupExpiredRequests()

	for packet := range interceptor.Packets() {
		handlePacket(elm, packet)
	}
}
