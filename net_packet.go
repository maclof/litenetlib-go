package litenetlib

const (
	NetPacketProperty_Unreliable         byte = 0
	NetPacketProperty_Channeled          byte = 1
	NetPacketProperty_Ack                byte = 2
	NetPacketProperty_Ping               byte = 3
	NetPacketProperty_Pong               byte = 4
	NetPacketProperty_ConnectRequest     byte = 5
	NetPacketProperty_ConnectAccept      byte = 6
	NetPacketProperty_Disconnect         byte = 7
	NetPacketProperty_UnconnectedMessage byte = 8
	NetPacketProperty_MtuCheck           byte = 9
	NetPacketProperty_MtuOk              byte = 10
	NetPacketProperty_Broadcast          byte = 11
	NetPacketProperty_Merged             byte = 12
	NetPacketProperty_ShutdownOk         byte = 13
	NetPacketProperty_PeerNotFound       byte = 14
	NetPacketProperty_InvalidProtocol    byte = 15
	NetPacketProperty_NatMessage         byte = 16
	NetPacketProperty_Empty              byte = 17
)

type NetPacket struct {
	bytes []byte
}

func (netPacket *NetPacket) Property() byte {
	return netPacket.bytes[0]
}

func (netPacket *NetPacket) Length() int {
	return len(netPacket.bytes)
}

func NewNetPacket(buf []byte, start int, end int) (*NetPacket, error) {
	return &NetPacket{
		bytes: buf[start:end],
	}, nil
}
