package metrics

import (
	"net"
)

type newlineDelimPacketSender struct {
	conn        *net.UDPConn
	buffer      []byte
	lastNewline int
	packetSize  int
}

func (s *newlineDelimPacketSender) Write(b []byte) (int, error) {
	s.maybeFlush(len(b))
	s.buffer = append(s.buffer, b...)
	return len(b), nil
}

func (s *newlineDelimPacketSender) WriteNewline() {
	s.maybeFlush(1)
	s.lastNewline = len(s.buffer)
	s.buffer = append(s.buffer, byte('\n'))
}

func (s *newlineDelimPacketSender) maybeFlush(n int) {
	if len(s.buffer)+n < s.packetSize {
		return
	}

	s.conn.WriteToUDP(s.buffer[:s.lastNewline], nil)
	copy(s.buffer[0:], s.buffer[s.lastNewline:])
	s.buffer = s.buffer[:len(s.buffer)-s.lastNewline]
}
