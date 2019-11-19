package main

import (
	P "csis3200/mock/public"
	"encoding/json"
	"net"
	"time"
)

func main() {
	// Get a new mock server cluster
	var servers = P.GetServerCluster()

	// Make a UDP connection to the server on localhost
	Conn, _ := net.DialUDP("udp", nil, &net.UDPAddr{IP: []byte{127, 0, 0, 1}, Port: 1337, Zone: ""})

	defer Conn.Close()

	// Send messages every 100 ms
	var limiter = time.Tick(100 * time.Millisecond)

	var lastTs = int64(0)

	for {
		<-limiter

		var ts = time.Now().UnixNano() / 1000000

		for _, s := range servers {
			var rate = s.GetRequestRate()

			// Convert RPM to the right number of messages to send every 100 ms
			for i := 0; i < rate / 60 / 10; i++ {
				_ = json.NewEncoder(Conn).Encode(s.GetLogMessage(ts))
			}
		}

		// If it's been more than 1 second, send a health ping for all servers too
		if ts - lastTs > 1000 {
			lastTs = ts

			for _, s := range servers {
				_ = json.NewEncoder(Conn).Encode(s.GetSystemMessage(ts))
			}
		}
	}
}
