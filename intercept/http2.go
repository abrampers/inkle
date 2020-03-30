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

	frames []http2.Frame
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

func (h *HTTP2) Frames() []http2.Frame {
	return h.frames
}

func (h *HTTP2) DecodeFromBytes(data []byte, df gopacket.DecodeFeedback) error {
	var frames []http2.Frame
	frameHeaderLength := uint32(9)
	payloadLength := len(data)

	payloadIdx := 0
	for payloadIdx < payloadLength {
		if payloadIdx+int(frameHeaderLength) > payloadLength {
			return fmt.Errorf("Payload length couldn't contain Frame Headers")
		}

		framePayloadLength := (uint32(data[payloadIdx+0])<<16 | uint32(data[payloadIdx+1])<<8 | uint32(data[payloadIdx+2]))
		frameLength := int(frameHeaderLength + framePayloadLength)

		rBit := data[payloadIdx+5] >> 7

		if rBit != 0 {
			return fmt.Errorf("R bit is not unset")
		}

		if payloadIdx+frameLength > payloadLength {
			return fmt.Errorf("Payload length couldn't contain Payload with the length mentioned in Frame Header")
		}

		r := bytes.NewReader(data[payloadIdx : payloadIdx+frameLength])
		framer := http2.NewFramer(nil, r)

		frame, err := framer.ReadFrame()
		if err != nil {
			return err
		}
		frames = append(frames, frame)

		payloadIdx += int(frameLength)
	}

	if payloadIdx != payloadLength {
		return fmt.Errorf("Payload length is not equal with the Frame length mentioned in Frame Header")
	}

	h.BaseLayer = layers.BaseLayer{Contents: data[:len(data)]}
	h.frames = frames
	return nil
}
