package litenetlib

import (
	"log"
	"net"
	"sync"
	"errors"

	"github.com/enriquebris/goconcurrentqueue"
)

type NetManager struct {
	config              *NetManagerConfig
	listener            INetListener
	socket              *NetSocket
	stats               *NetStatistics
	netEventsQueue      *goconcurrentqueue.FIFO
	netEventsQueueMutex sync.Mutex
	unsyncedEvents      bool
	packetsReceived     int
	bytesReceived       int
}

type NetManagerConfig struct {
	AddrV4                     string
	PortV4                     int
	AddrV6                     string
	PortV6                     int
	StatsEnabled               bool
	BroadcastReceiveEnabled    bool
	UnconnectedMessagesEnabled bool
	NatPunchEnabled            bool
}

func (netManager *NetManager) Start() error {
	log.Printf("Starting net manager...")

	if netManager.config.AddrV4 != "" && netManager.config.PortV4 != 0 {
		err := netManager.socket.BindV4(netManager.config.AddrV4, netManager.config.PortV4)
		if err != nil {
			return err
		}
	}

	if netManager.config.AddrV6 != "" && netManager.config.PortV6 != 0 {
		err := netManager.socket.BindV6(netManager.config.AddrV6, netManager.config.PortV6)
		if err != nil {
			return err
		}
	}

	if !netManager.socket.IsListening() {
		return errors.New("No valid socket listeners have been setup.")
	}

	return nil
}

func (netManager *NetManager) Stop() {
	if netManager.socket == nil {
		return
	}
	log.Printf("Stopping net manager...")
	netManager.socket.Close()
}

func (netManager *NetManager) PollEvents() {
	if netManager.unsyncedEvents {
		return
	}
	netManager.netEventsQueueMutex.Lock()
	defer netManager.netEventsQueueMutex.Unlock()
	// log.Println("PollEvents()")
	for {
		if netManager.netEventsQueue.GetLen() == 0 {
			// log.Println("Queue is empty")
			return
		} else {
			event, err := netManager.netEventsQueue.Dequeue()
			if err != nil {
				return
			}

			netManager.processEvent(event)
		}
	}
}

func (netManager *NetManager) processEvent(event interface{}) {
	log.Println("processEvent")
}

func (netManager *NetManager) Stats() *NetStatistics {
	return netManager.stats
}

type socketListener struct {
	netManager *NetManager
}

func (listener *socketListener) OnMessageReceived(numBytes int, buf []byte, addr *net.UDPAddr) {
	log.Println("OnMessageReceived()")

	netManager := listener.netManager

	if netManager.config.StatsEnabled {
		stats := netManager.stats
		stats.PacketsReceived++;
		stats.BytesReceived += int64(numBytes);

		log.Printf("Packets received: %d", stats.PacketsReceived)
		log.Printf("Bytes received: %d", stats.BytesReceived)
	}

	if buf[0] == NetPacketProperty_Empty {
		return
	}

	packet, err := NewNetPacket(buf, 0, numBytes)
	if err != nil {
		return
	}

	log.Printf("Packet length: %d", packet.Length())

	switch packet.Property() {
	case NetPacketProperty_ConnectRequest:
		log.Println("ConnectRequest")
		// if (NetConnectRequestPacket.GetProtocolId(packet) != NetConstants.ProtocolId)
		// SendRawAndRecycle(NetPacketPool.GetWithProperty(PacketProperty.InvalidProtocol), remoteEndPoint);
		break

	case NetPacketProperty_Broadcast:
		if !netManager.config.BroadcastReceiveEnabled {
			return
		}
		log.Println("Broadcast")
		// CreateEvent(NetEvent.EType.Broadcast, remoteEndPoint: remoteEndPoint, readerSource: packet);
		break

	case NetPacketProperty_UnconnectedMessage:
		if !netManager.config.UnconnectedMessagesEnabled {
			return
		}
		log.Println("UnconnectedMessage")
		// CreateEvent(NetEvent.EType.ReceiveUnconnected, remoteEndPoint: remoteEndPoint, readerSource: packet);
		break

	case NetPacketProperty_NatMessage:
		if !netManager.config.NatPunchEnabled {
			return
		}
		log.Println("NatMessage")
		// NatPunchModule.ProcessMessage(remoteEndPoint, packet);
		break
	}
}

func NewNetManager(config *NetManagerConfig, listener INetListener) *NetManager {
	netManager := &NetManager{
		config: config,
		listener: listener,
		unsyncedEvents: false,
		netEventsQueue: goconcurrentqueue.NewFIFO(),
	}

	netManager.socket = NewNetSocket(&socketListener{
		netManager: netManager,
	})

	if config.StatsEnabled {
		netManager.stats = NewNetStatistics()
	}

	return netManager
}


