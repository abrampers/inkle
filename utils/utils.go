package utils

import (
	"fmt"
	"net"
	"regexp"

	"github.com/google/gopacket/pcap"
)

// Regular expression to parse gRPC service name and method name.
var re *regexp.Regexp = regexp.MustCompile(`/([a-zA-Z0-9\.]+)/([a-zA-Z0-9\.]+)(.*)`)

func ParseGrpcPath(path string) (servicename string, methodname string, err error) {
	matches := re.FindStringSubmatch(path)
	if len(matches) != 4 || matches[3] != "" {
		return "", "", fmt.Errorf("Failed to match path")
	}
	return matches[1], matches[2], nil
}

func CIDR(dname string) *net.IPNet {
	devices, err := pcap.FindAllDevs()
	if err != nil {
		return &net.IPNet{}
	}

	for _, device := range devices {
		if device.Name == dname {
			address := device.Addresses[0]
			return &net.IPNet{IP: address.IP, Mask: address.Netmask}
		}
	}

	return &net.IPNet{}
}
