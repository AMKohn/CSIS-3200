package ingester

import (
	"csis3200/internal/app/processor"
	"encoding/json"
	"net"
	"sync"
)

func StartServer(wg *sync.WaitGroup) {
	ServerConn, _ := net.ListenUDP("udp", &net.UDPAddr{IP: []byte{0, 0, 0, 0}, Port: 1337, Zone: ""})

	buf := make([]byte, 1024)

	println("Listening for UDP packets on port 1337")

	// Use a goroutine so we don't block while waiting for packets
	go func() {
		defer ServerConn.Close()

		// Tell the WaitGroup we're done after this function finishes
		defer wg.Done()

		for {
			n, _, _ := ServerConn.ReadFromUDP(buf)
			//println("Received ", string(buf[0:n]), " from ", addr)// addr is the second return value in ReadFromUDP


			var result map[string]interface{}
			json.Unmarshal(buf[0:n], &result)

			processor.HandleMessage(result)
		}
	}()
}
