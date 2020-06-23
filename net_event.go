package litenetlib

import (
	"net"
)

const (
	NetEventType_Connect                  byte = 0
    NetEventType_Disconnect               byte = 1
    NetEventType_Receive                  byte = 2
    NetEventType_ReceiveUnconnected       byte = 3
    NetEventType_Error                    byte = 4
    NetEventType_ConnectionLatencyUpdated byte = 5
    NetEventType_Broadcast                byte = 6
    NetEventType_ConnectionRequest        byte = 7
    NetEventType_MessageDelivered         byte = 8
)

type NetEvent struct {
	eventType      byte
	addr           *net.UDPAddr
	packet         *NetPacket
	connectRequest *ConnectRequest
}
