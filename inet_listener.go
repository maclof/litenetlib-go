package litenetlib

import (
	"net"
)

type INetListener interface {
	OnMessageReceived(numBytes int, buf []byte, addr *net.UDPAddr)
}
