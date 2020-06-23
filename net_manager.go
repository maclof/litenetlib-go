package litenetlib

import (
	"log"
	"net"
	"sync"
	"errors"

	"github.com/emirpasic/gods/maps/hashmap"
	"github.com/enriquebris/goconcurrentqueue"
)

type NetManager struct {
	config               *NetManagerConfig
	listener             INetListener
	socket               *NetSocket
	stats                *NetStatistics
	peers                *hashmap.Map
	peersMutex           sync.Mutex
	connectRequests      *hashmap.Map
	connectRequestsMutex sync.Mutex
	netEventsQueue       *goconcurrentqueue.FIFO
	netEventsQueueMutex  sync.Mutex
	unsyncedEvents       bool
	packetsReceived      int
	bytesReceived        int
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
			eventInterface, err := netManager.netEventsQueue.Dequeue()
			if err != nil {
				return
			}

			netManager.processEvent(eventInterface.(*NetEvent))
		}
	}
}

func (netManager *NetManager) Stats() *NetStatistics {
	return netManager.stats
}

func (netManager *NetManager) processEvent(event *NetEvent) {
	log.Println("processEvent")
}

func (netManager *NetManager) dataReceived(numBytes int, buf []byte, addr *net.UDPAddr) {
	if buf[0] == NetPacketProperty_Empty {
		return
	}

	packet, err := NewNetPacket(buf, 0, numBytes)
	if err != nil {
		return
	}

	if netManager.config.StatsEnabled {
		stats := netManager.stats
		stats.PacketsReceived++;
		stats.BytesReceived += int64(packet.Length());

		log.Printf("Packets received: %d", stats.PacketsReceived)
		log.Printf("Bytes received: %d", stats.BytesReceived)
	}

	log.Printf("Packet length: %d", packet.Length())

	switch packet.Property() {
	case NetPacketProperty_ConnectRequest:
		// todo:
		// if (NetConnectRequestPacket.GetProtocolId(packet) != NetConstants.ProtocolId)
		// SendRawAndRecycle(NetPacketPool.GetWithProperty(PacketProperty.InvalidProtocol), remoteEndPoint);
		break

	case NetPacketProperty_Broadcast:
		if !netManager.config.BroadcastReceiveEnabled {
			return
		}
		netManager.newNetEvent(&NetEvent{
			eventType: NetEventType_Broadcast,
			addr: addr,
			packet: packet,
		})
		break

	case NetPacketProperty_UnconnectedMessage:
		if !netManager.config.UnconnectedMessagesEnabled {
			return
		}
		netManager.newNetEvent(&NetEvent{
			eventType: NetEventType_ReceiveUnconnected,
			addr: addr,
			packet: packet,
		})
		break

	case NetPacketProperty_NatMessage:
		if !netManager.config.NatPunchEnabled {
			return
		}
		// todo:
		// NatPunchModule.ProcessMessage(remoteEndPoint, packet);
		break
	}

	netManager.peersMutex.Lock()
	peerInterface, peerFound := netManager.peers.Get(addr.String())
	netManager.peersMutex.Unlock()

	var peer *NetPeer
	if peerFound {
		peer = peerInterface.(*NetPeer)
	}

	switch packet.Property() {
	case NetPacketProperty_ConnectRequest:
		netManager.processConnectRequest(peer, addr, packet);
		break

	case NetPacketProperty_PeerNotFound:
		log.Println("PeerNotFound")
		break

	case NetPacketProperty_InvalidProtocol:
		log.Println("InvalidProtocol")
		break

	case NetPacketProperty_Disconnect:
		log.Println("Disconnect")
		break

	case NetPacketProperty_ConnectAccept:
		log.Println("ConnectAccept")
		break

	default:
		if peerFound {
			peer.ProcessPacket(packet)
		} else {
			// todo:
			// SendRawAndRecycle(NetPacketPool.GetWithProperty(PacketProperty.PeerNotFound), remoteEndPoint);
		}
		break
	}
}

func (netManager *NetManager) newNetEvent(event *NetEvent) {
	netManager.netEventsQueueMutex.Lock()
	netManager.netEventsQueue.Enqueue(event)
	netManager.netEventsQueueMutex.Unlock()
}

func (netManager *NetManager) processConnectRequest(peer *NetPeer, addr *net.UDPAddr, packet *NetPacket) {
	if peer != nil {

	}

	netManager.connectRequestsMutex.Lock()
	connectRequestInterface, connectRequestFound := netManager.connectRequests.Get(addr.String())
	netManager.connectRequestsMutex.Unlock()

	var connectRequest *ConnectRequest
	if connectRequestFound {
		connectRequest = connectRequestInterface.(*ConnectRequest)
		connectRequest.UpdateRequest(packet)
	} else {
		netManager.connectRequests.Put(addr.String(), &ConnectRequest{
			addr: addr,
			packet: packet,
		})
	}
	
	netManager.newNetEvent(&NetEvent{
		eventType: NetEventType_ConnectionRequest,
		connectRequest: connectRequest,
	})
}

type netManagerSocketListener struct {
	netManager *NetManager
}

func (listener *netManagerSocketListener) OnMessageReceived(numBytes int, buf []byte, addr *net.UDPAddr) {
	listener.netManager.dataReceived(numBytes, buf, addr)
}

func NewNetManager(config *NetManagerConfig, listener INetListener) *NetManager {
	netManager := &NetManager{
		config: config,
		listener: listener,
		peers: hashmap.New(),
		connectRequests: hashmap.New(),
		netEventsQueue: goconcurrentqueue.NewFIFO(),
		unsyncedEvents: false,
	}

	netManager.socket = NewNetSocket(&netManagerSocketListener{
		netManager: netManager,
	})

	if config.StatsEnabled {
		netManager.stats = NewNetStatistics()
	}

	return netManager
}
