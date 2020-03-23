package intercept

import (
	"log"
	"time"

	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"github.com/google/gopacket/pcap"
)

var (
	p InterceptedPacket
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

	var h2c HTTP2
	parser := gopacket.NewDecodingLayerParser(LayerTypeHTTP2, &h2c)

	decoded := []gopacket.LayerType{}
	for packet := range i.source.Packets() {
		ipLayer := packet.NetworkLayer()
		if ipLayer == nil {
			continue
		}

		ipv4, ipv4Ok := ipLayer.(*layers.IPv4)
		ipv6, ipv6Ok := ipLayer.(*layers.IPv6)
		if !ipv4Ok && !ipv6Ok {
			continue
		}

		tcpLayer := packet.Layer(layers.LayerTypeTCP)
		if tcpLayer == nil {
			continue
		}

		tcp, ok := tcpLayer.(*layers.TCP)
		if !ok {
			continue
		}

		appLayer := packet.ApplicationLayer()
		if appLayer == nil {
			continue
		}

		packetData := appLayer.Payload()
		if err := parser.DecodeLayers(packetData, &decoded); err != nil {
			continue
		}

		if ipv4Ok {
			p = InterceptedPacket{
				IsIPv4: ipv4Ok,
				IPv4:   *ipv4,
				IPv6:   layers.IPv6{},
				TCP:    *tcp,
				HTTP2:  h2c,
			}
		} else {
			p = InterceptedPacket{
				IsIPv4: ipv4Ok,
				IPv4:   layers.IPv4{},
				IPv6:   *ipv6,
				TCP:    *tcp,
				HTTP2:  h2c,
			}
		}

		i.c <- p
	}
}
