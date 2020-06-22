
# LiteNetLib-Go
A Golang port of the LiteNetLib Lite reliable UDP library (https://github.com/RevenantX/LiteNetLib)

## Usage samples

### Client
```go

```

### Server
```go
package main

import (
	"log"
	"net"
	"time"

	"github.com/maclof/litenetlib-go"
)

type NetListener struct {}

func (netListener *NetListener) OnMessageReceived(numBytes int, buf []byte, addr *net.UDPAddr) {
	log.Println("OnMessageReceived()")
}

func main() {
	listener := &NetListener{}

	netManager := litenetlib.NewNetManager(&litenetlib.NetManagerConfig{
		AddrV4: "127.0.0.1",
		PortV4: 9050,
		// AddrV6: "::1",
		// PortV6: 9050,
	}, listener)

	err := netManager.Start()
	if err != nil {
		log.Fatalln(err.Error())
	}

	for {
		netManager.PollEvents()

		time.Sleep(time.Millisecond * 5)
	}

	netManager.Stop()
}
```
