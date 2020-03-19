package intercept

import (
	"bytes"
	"fmt"

	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"golang.org/x/net/http2"
)

// Create a layer type, should be unique and high, so it doesn't conflict,
// giving it a name and a decoder to use.
var LayerTypeHTTP2 = gopacket.RegisterLayerType(12345, gopacket.LayerTypeMetadata{Name: "HTTP2", Decoder: gopacket.DecodeFunc(decodeHTTP2)})

// Implement my layer
type HTTP2 struct {
	layers.BaseLayer

	frame http2.Frame
}

func (h HTTP2) LayerType() gopacket.LayerType      { return LayerTypeHTTP2 }
func (h *HTTP2) Payload() []byte                   { return nil }
func (h *HTTP2) CanDecode() gopacket.LayerClass    { return LayerTypeHTTP2 }
func (h *HTTP2) NextLayerType() gopacket.LayerType { return gopacket.LayerTypeZero }

// Now implement a decoder... this one strips off the first 4 bytes of the
// packet.
func decodeHTTP2(data []byte, p gopacket.PacketBuilder) error {
	h := &HTTP2{}
	err := h.DecodeFromBytes(data, p)
	if err != nil {
		return err
	}
	p.AddLayer(h)
	p.SetApplicationLayer(h)
	return nil
}

func (h *HTTP2) DecodeFromBytes(data []byte, df gopacket.DecodeFeedback) error {
	err := validateHTTP2(data)
	if err != nil {
		return err
	}

	var framerOutput bytes.Buffer
	r := bytes.NewReader(data)
	framer := http2.NewFramer(&framerOutput, r)

	h.BaseLayer = layers.BaseLayer{Contents: data[:len(data)]}

	frame, err := framer.ReadFrame()
	if err != nil {
		return err
	}
	h.frame = frame

	return nil
}

func validateHTTP2(payload []byte) error {
	frameHeaderLength := uint32(9)
	payloadLength := len(payload)

	payloadIdx := 0
	for payloadIdx < payloadLength {
		if payloadIdx+int(frameHeaderLength) > payloadLength {
			return fmt.Errorf("packet length is not equal with the packet length mentioned in frame header")
		}

		frameLength := (uint32(payload[payloadIdx+0])<<16 | uint32(payload[payloadIdx+1])<<8 | uint32(payload[payloadIdx+2]))
		rBit := payload[payloadIdx+5] >> 7

		if rBit != 0 {
			return fmt.Errorf("R bit is not unset")
		}

		payloadIdx += int(frameLength + frameHeaderLength)
	}

	if payloadIdx != payloadLength {
		return fmt.Errorf("packet length is not equal with the packet length mentioned in frame header")
	}
	return nil
}
