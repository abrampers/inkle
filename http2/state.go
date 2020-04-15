package http2

type ipTcpConn struct {
	SrcIP  string
	SrcTCP uint16
	DstIP  string
	DstTCP uint16
}

type HeadersState struct {
	state map[ipTcpConn]map[string]string
}

var State = &HeadersState{state: map[ipTcpConn]map[string]string{}}

func (s *HeadersState) Headers(srcip string, srctcp uint16, dstip string, dsttcp uint16) map[string]string {
	if val, ok := s.state[ipTcpConn{srcip, srctcp, dstip, dsttcp}]; ok && val != nil {
		return val
	}
	return map[string]string{}
}

func (s *HeadersState) SetHeaders(srcip string, srctcp uint16, dstip string, dsttcp uint16, key string, value string) {
	conn := ipTcpConn{srcip, srctcp, dstip, dsttcp}
	_, ok := s.state[conn]
	if !ok {
		s.state[conn] = map[string]string{}
	}
	s.state[conn][key] = value
}

func (s *HeadersState) UpdateState(srcip string, srctcp uint16, dstip string, dsttcp uint16, headers map[string]string) {
	for k, v := range headers {
		s.SetHeaders(srcip, srctcp, dstip, dsttcp, k, v)
	}
}
