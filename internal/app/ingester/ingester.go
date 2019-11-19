package ingester

import (
	"csis3200/internal/app/processor"
	"encoding/json"
	"net"
	"sync"
)

func StartServer(wg *sync.WaitGroup) {
	// Listen for UDP packets on port 1337
	ServerConn, _ := net.ListenUDP("udp", &net.UDPAddr{IP: []byte{0, 0, 0, 0}, Port: 1337, Zone: ""})

	// Make a 1024 byte buffer. The log messages shouldn't exceed this size, but it's possible that they will
	buf := make([]byte, 1024)

	println("Listening for UDP packets on port 1337")

	// Use a goroutine so we don't block while waiting for packets
	go func() {
		defer ServerConn.Close()

		// Tell the WaitGroup we're done after this function finishes
		defer wg.Done()

		// Read the next datagram infinitely
		for {
			n, _, _ := ServerConn.ReadFromUDP(buf)

			var result map[string]interface{}

			// Parse the JSON into result
			if err := json.Unmarshal(buf[0:n], &result); err != nil {
				continue
			}

			processor.HandleMessage(result)
		}
	}()
}
