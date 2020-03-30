package intercept

import (
	"fmt"
	"log"
	"net"
	"time"

	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"github.com/google/gopacket/pcap"
)

var (
	p InterceptedPacket
)

type InterceptedPacket struct {
	SrcIP, DstIP   net.IP
	SrcTCP, DstTCP layers.TCPPort
	HTTP2          HTTP2
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

	for packet := range i.source.Packets() {
		itcpacket, err := extractPacket(packet)
		if err != nil {
			continue
		}

		i.c <- *itcpacket

	}
}

func extractPacket(packet gopacket.Packet) (*InterceptedPacket, error) {
	var h2c HTTP2
	parser := gopacket.NewDecodingLayerParser(LayerTypeHTTP2, &h2c)
	decoded := []gopacket.LayerType{}

	netlayer := packet.NetworkLayer()
	if netlayer == nil {
		return nil, fmt.Errorf("No Network Layer found")
	}

	ipv4, ipv4Ok := netlayer.(*layers.IPv4)
	ipv6, ipv6Ok := netlayer.(*layers.IPv6)
	if !ipv4Ok && !ipv6Ok {
		return nil, fmt.Errorf("Failed to cast Network Layer to IPv4 or IPv6")
	}

	tcplayer := packet.Layer(layers.LayerTypeTCP)
	if tcplayer == nil {
		return nil, fmt.Errorf("No TCP Layer found")
	}

	tcp, ok := tcplayer.(*layers.TCP)
	if !ok {
		return nil, fmt.Errorf("Failed to cast TCP Layer to TCP")
	}

	applayer := packet.ApplicationLayer()
	if applayer == nil {
		return nil, fmt.Errorf("No Application Layer found")
	}

	packetData := applayer.Payload()
	if err := parser.DecodeLayers(packetData, &decoded); err != nil {
		return nil, fmt.Errorf("Failed to parse Application Layer payload to HTTP2")
	}

	if ipv4Ok {
		return &InterceptedPacket{
			SrcIP:  ipv4.SrcIP,
			DstIP:  ipv4.DstIP,
			SrcTCP: tcp.SrcPort,
			DstTCP: tcp.DstPort,
			HTTP2:  h2c,
		}, nil
	} else {
		return &InterceptedPacket{
			SrcIP:  ipv6.SrcIP,
			DstIP:  ipv6.DstIP,
			SrcTCP: tcp.SrcPort,
			DstTCP: tcp.DstPort,
			HTTP2:  h2c,
		}, nil
	}
}
