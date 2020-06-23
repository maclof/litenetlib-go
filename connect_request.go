package litenetlib

import (
	"net"
)

type ConnectRequest struct {
	addr   *net.UDPAddr
	packet *NetPacket
}

func (connectRequest *ConnectRequest) UpdateRequest(packet *NetPacket) {
	connectRequest.packet = packet
}
