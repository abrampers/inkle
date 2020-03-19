package intercept

import (
	"fmt"
	"log"
	"time"

	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"github.com/google/gopacket/pcap"
)

var (
	device       string        = "lo0"
	snapshot_len int32         = 1024
	promiscuous  bool          = false
	timeout      time.Duration = 900 * time.Millisecond
)

type InterceptedPacket struct {
	IsIPv4 bool
	IPv4   layers.IPv4
	IPv6   layers.IPv6
	TCP    layers.TCP
	HTTP2  HTTP2
}

type PacketInterceptor struct {
	handle *pcap.Handle
	source *gopacket.PacketSource

	c chan InterceptedPacket
}

func NewPacketInterceptor(device string, snapshotLen int32, isPromiscuous bool, timeout time.Duration) *PacketInterceptor {
	handle, err := pcap.OpenLive(device, snapshotLen, isPromiscuous, timeout)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Successfully opened live sniffing on %s\n", device)

	source := gopacket.NewPacketSource(handle, handle.LinkType())
	return &PacketInterceptor{
		handle: handle,
		source: source,
	}
}

func (i *PacketInterceptor) Close() {
	i.handle.Close()
}

func (i *PacketInterceptor) Packets() chan InterceptedPacket {
	if i.c == nil {
		i.c = make(chan InterceptedPacket, 1000)
		go i.interceptPacket()
	}
	return i.c
}

func (i *PacketInterceptor) interceptPacket() {
	defer close(i.c)
	// Open device
	// handle, err := pcap.OpenLive(device, snapshot_len, promiscuous, timeout)
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// log.Printf("Successfully opened live sniffing on %s\n", device)
	// defer handle.Close()

	var h2c HTTP2
	parser := gopacket.NewDecodingLayerParser(LayerTypeHTTP2, &h2c)

	// Use the handle as a packet source to process all packets
	// source := gopacket.NewPacketSource(handle, handle.LinkType())
	decoded := []gopacket.LayerType{}
	for packet := range i.source.Packets() {
		ipLayer := packet.NetworkLayer()
		if ipLayer == nil {
			// log.Println("No IP")
			continue
		}

		ipv4, ipv4Ok := ipLayer.(*layers.IPv4)
		ipv6, ipv6Ok := ipLayer.(*layers.IPv6)
		if !ipv4Ok && !ipv6Ok {
			// log.Println("Failed to cast packet to IPv4 or IPv6")
			continue
		}

		tcpLayer := packet.Layer(layers.LayerTypeTCP)
		if tcpLayer == nil {
			// log.Println("Not a TCP Packet")
			continue
		}

		tcp, ok := tcpLayer.(*layers.TCP)
		if !ok {
			// log.Println("Failed to cast packet to TCP")
			continue
		}

		appLayer := packet.ApplicationLayer()
		if appLayer == nil {
			// log.Println("No ApplicationLayer payload")
			continue
		}

		packetData := appLayer.Payload()
		if err := parser.DecodeLayers(packetData, &decoded); err != nil {
			// fmt.Printf("Could not decode layers: %v\n", err)
			continue
		}

		// fmt.Println("*****************************************************")
		// if ipv4Ok {
		// 	fmt.Println("IPv4 SrcIP:        ", ipv4.SrcIP)
		// 	fmt.Println("IPv4 DstIP:        ", ipv4.DstIP)
		// } else if ipv6Ok {
		// 	fmt.Println("IPv6 SrcIP:        ", ipv6.SrcIP)
		// 	fmt.Println("IPv6 DstIP:        ", ipv6.DstIP)
		// }
		// fmt.Println("TCP srcPort:       ", tcp.SrcPort)
		// fmt.Println("TCP dstPort:       ", tcp.DstPort)
		// fmt.Println("HTTP/2:            ", h2c.frame)
		// fmt.Println("*****************************************************")
		p := InterceptedPacket{
			IsIPv4: ipv4Ok,
			IPv4:   *ipv4,
			IPv6:   *ipv6,
			TCP:    *tcp,
			HTTP2:  h2c,
		}

		i.c <- p
	}
}
