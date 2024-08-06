package metrics

import (
	"net"
	"net/netip"
)

type newlineDelimPacketSender struct {
	conn        *net.UDPConn
	buffer      []byte
	lastNewline int
	packetSize  int
}

func newNewlineDelimPacketSender(addr netip.AddrPort) (*newlineDelimPacketSender, error) {
	conn, err := net.DialUDP(
		"udp",
		nil, // laddr
		net.UDPAddrFromAddrPort(addr),
	)
	if err != nil {
		return nil, err
	}

	const packetSize = 1400

	return &newlineDelimPacketSender{
		conn:       conn,
		buffer:     make([]byte, 0, packetSize*2),
		packetSize: packetSize,
	}, nil
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

func (s *newlineDelimPacketSender) Flush() {
	s.conn.WriteMsgUDP(s.buffer[:s.lastNewline], nil /*oob*/, nil /*addr*/)
	copy(s.buffer[0:], s.buffer[s.lastNewline:])
	s.buffer = s.buffer[:len(s.buffer)-s.lastNewline]
}

func (s *newlineDelimPacketSender) maybeFlush(n int) {
	if len(s.buffer)+n < s.packetSize {
		return
	}
	s.Flush()
}
