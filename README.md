
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
	"time"

	"github.com/maclof/litenetlib-go"
)

func main() {
	netManager := litenetlib.NewNetManager()
	err := netManager.Start("127.0.0.1", 9050)
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
