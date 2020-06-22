package litenetlib

import (
	"log"
	"sync"
	"errors"

	"github.com/enriquebris/goconcurrentqueue"
)

type NetManager struct {
	config *NetManagerConfig
	listener INetListener
	socket *NetSocket
	netEventsQueue *goconcurrentqueue.FIFO
	netEventsQueueMutex sync.Mutex
	unsyncedEvents bool
}

type NetManagerConfig struct {
	AddrV4 string
	PortV4 int
	AddrV6 string
	PortV6 int
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
	log.Println("PollEvents()")
	for {
		if netManager.netEventsQueue.GetLen() == 0 {
			log.Println("Queue is empty")
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

}

func NewNetManager(config *NetManagerConfig, listener INetListener) *NetManager {
	return &NetManager{
		config: config,
		listener: listener,
		socket: NewNetSocket(listener),
		unsyncedEvents: false,
		netEventsQueue: goconcurrentqueue.NewFIFO(),
	}
}


