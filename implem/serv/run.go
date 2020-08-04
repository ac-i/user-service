package serv

import (
	"log"
	"time"
)

// main start a gRPC server and waits for connection
func RunServerMAIN() {
	// chTheEnd channel to signal TheEnd of the function /otherwais waits
	chTheEnd := make(chan struct{})

	// fire the gRPC server with default config in a goroutine
	go func() {
		err := runServerGRCP()
		if err != nil {
			log.Fatalf("ERR runServerGRCP: %s", err)
		}
	}()

	// log order (can be removed)
	time.Sleep(time.Millisecond)

	// fire the runServerRESTGW - REST server with default config in a goroutine
	go func() {
		err := runServerRESTGW()
		if err != nil {
			log.Fatalf("ERR runServerRESTGW: %s", err)
		}
	}()

	// waits for signal on chTheEnd like: chTheEnd <- struct{}{}
	<-chTheEnd
}
