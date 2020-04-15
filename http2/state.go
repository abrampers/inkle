package http2

type ipTcpConn struct {
	SrcIP  string
	SrcTCP uint32
	DstIP  string
	DstTCP uint32
}

type HeadersState struct {
	state map[ipTcpConn]map[string]string
}

var State = &HeadersState{}

func (s *HeadersState) Headers(srcip string, srctcp uint32, dstip string, dsttcp uint32) map[string]string {
	if val, ok := s.state[ipTcpConn{srcip, srctcp, dstip, dsttcp}]; ok && val != nil {
		return val
	}
	return map[string]string{}
}

func (s *HeadersState) SetHeaders(srcip string, srctcp uint32, dstip string, dsttcp uint32, key string, value string) {
	conn := ipTcpConn{srcip, srctcp, dstip, dsttcp}
	_, ok := s.state[conn]
	if !ok {
		s.state[conn] = map[string]string{}
	}
	s.state[conn][key] = value
}
