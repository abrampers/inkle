package http2

import (
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/hpack"
)

var he headers = headers{headers: make(map[string]string)}
var decoder *hpack.Decoder = hpack.NewDecoder(4096, he.readHeaderField)

type headers struct {
	headers map[string]string
}

func Headers(headersframe http2.HeadersFrame) map[string]string {
	he.clean()
	decoder.Write(headersframe.HeaderBlockFragment())
	return he.headers
}

func (h *headers) clean() {
	for key := range h.headers {
		delete(h.headers, key)
	}
}

func (h *headers) readHeaderField(f hpack.HeaderField) {
	h.headers[f.Name] = f.Value
}
