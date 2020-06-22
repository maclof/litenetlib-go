package litenetlib

import (
	"log"
	"fmt"
	"net"
	"sync"

	"github.com/enriquebris/goconcurrentqueue"
)

type NetManager struct {
	serverConn *net.UDPConn
	netEventsQueue *goconcurrentqueue.FIFO
	netEventsQueueMutex sync.Mutex
	unsyncedEvents bool
}

func (netManager *NetManager) Start(addr string, port int) error {
	log.Printf("Listening for UDP packets on: %s:%d", addr, port)

	listenAddr, err := net.ResolveUDPAddr("udp4", fmt.Sprintf("%s:%d", addr, port))
	if err != nil {
		return err
	}

	serverConn, err := net.ListenUDP("udp4", listenAddr)
	if err != nil {
		return err
	}
	netManager.serverConn = serverConn

	return nil
}

func (netManager *NetManager) Stop() {
	if netManager.serverConn == nil {
		return
	}
	log.Printf("Stopping listening for UDP packets.")
	netManager.serverConn.Close()
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

func NewNetManager() *NetManager {
	return &NetManager{
		unsyncedEvents: false,
		netEventsQueue: goconcurrentqueue.NewFIFO(),
	}
}


