package main

import (
	P "csis3200/mock/public"
	"encoding/json"
	"net"
	"time"
)

func main() {
	var servers = P.GetServerCluster()

	Conn, _ := net.DialUDP("udp", nil, &net.UDPAddr{IP: []byte{127, 0, 0, 1}, Port: 1337, Zone: ""})

	defer Conn.Close()

	// Send messages every 100 ms
	limiter := time.Tick(100 * time.Millisecond)

	var lastTs = int64(0)

	for {
		<-limiter

		var ts = time.Now().UnixNano() / 1000000

		for _, s := range servers {
			var rate = s.GetRequestRate()

			// Convert RPM to the right number to send every 10 ms
			for i := 0; i < rate / 60 / 10; i++ {
				_ = json.NewEncoder(Conn).Encode(s.GetLogMessage(ts))
			}
		}

		// If it's been more than 1 second, send a health ping too
		if ts - lastTs > 1000 {
			lastTs = ts

			for _, s := range servers {
				_ = json.NewEncoder(Conn).Encode(s.GetSystemMessage(ts))
			}
		}
	}
}
