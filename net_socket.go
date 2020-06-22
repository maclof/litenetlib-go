package litenetlib

import (
	"log"
	"fmt"
	"net"
)

type NetSocket struct {
	listener INetListener
	isRunningV4 bool
	udpConnV4 *net.UDPConn
	isRunningV6 bool
	udpConnV6 *net.UDPConn
}

func (netSocket *NetSocket) BindV4(addr string, port int) error {
	log.Printf("Binding V4 UDP socket: %s:%d", addr, port)

	listenAddr, err := net.ResolveUDPAddr("udp4", fmt.Sprintf("%s:%d", addr, port))
	if err != nil {
		return err
	}

	udpConnV4, err := net.ListenUDP("udp4", listenAddr)
	if err != nil {
		return err
	}
	netSocket.udpConnV4 = udpConnV4

	go netSocket.ReceiveV4Packets()

	return nil
}

func (netSocket *NetSocket) BindV6(addr string, port int) error {
	log.Printf("Binding V6 UDP socket: [%s]:%d", addr, port)

	listenAddr, err := net.ResolveUDPAddr("udp6", fmt.Sprintf("[%s]:%d", addr, port))
	if err != nil {
		return err
	}

	udpConnV6, err := net.ListenUDP("udp6", listenAddr)
	if err != nil {
		return err
	}
	netSocket.udpConnV6 = udpConnV6

	go netSocket.ReceiveV6Packets()

	return nil
}

func (netSocket *NetSocket) ReceiveV4Packets() {
	var buf []byte
	for {
		if !netSocket.isRunningV4 {
			break
		}

		buf = make([]byte, MAX_PACKET_SIZE)
		numBytes, addr, err := netSocket.udpConnV4.ReadFromUDP(buf)
		if err != nil {
			continue
		}

		debugMsg := fmt.Sprintf("Received V4 UDP packet: len(%d) - ", len(buf))
		for i := 0; i < len(buf); i++ {
			debugMsg += fmt.Sprintf("%d ", buf[i])
		}
		log.Println(debugMsg)

		netSocket.listener.OnMessageReceived(numBytes, buf, addr);
	}
}

func (netSocket *NetSocket) ReceiveV6Packets() {
	var buf []byte
	for {
		if !netSocket.isRunningV6 {
			break
		}

		buf = make([]byte, MAX_PACKET_SIZE)
		numBytes, addr, err := netSocket.udpConnV6.ReadFromUDP(buf)
		if err != nil {
			continue
		}

		debugMsg := fmt.Sprintf("Received V6 UDP packet: len(%d) - ", len(buf))
		for i := 0; i < len(buf); i++ {
			debugMsg += fmt.Sprintf("%d ", buf[i])
		}
		log.Println(debugMsg)

		netSocket.listener.OnMessageReceived(numBytes, buf, addr);
	}
}

func (netSocket *NetSocket) Close() {
	netSocket.CloseV4()
	netSocket.CloseV6()
}

func (netSocket *NetSocket) CloseV4() {
	netSocket.isRunningV4 = false
	if netSocket.udpConnV4 == nil {
		return
	}
	netSocket.udpConnV4.Close()
	netSocket.udpConnV4 = nil
}

func (netSocket *NetSocket) CloseV6() {
	netSocket.isRunningV6 = false
	if netSocket.udpConnV6 == nil {
		return
	}
	netSocket.udpConnV6.Close()
	netSocket.udpConnV6 = nil
}

func (netSocket *NetSocket) IsListening() bool {
	return netSocket.udpConnV4 != nil || netSocket.udpConnV6 != nil
}

func NewNetSocket(listener INetListener) *NetSocket {
	return &NetSocket{
		listener: listener,
	}
}
