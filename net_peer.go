package litenetlib

import (
	"log"
)

type NetPeer struct {
	
}

func (netPeer *NetPeer) ProcessPacket(packet *NetPacket) {
	log.Println("ProcessPacket")
}
