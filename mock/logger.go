package main

import (
	"fmt"
	"math/rand"
	"net"
	"time"
)

func getRandomWebRequest() string {
	var statusCode = 200

	// Only 10% of requests have non-200 status codes
	if rand.Intn(10) == 1 {
		statusCode = rand.Intn(5) * 100
	}

	return fmt.Sprint(`{
		"type": "web_request",
		"server_id": "app-`, rand.Intn(60) + 75, `",
		"request_type": "GET",
		"path": "/dashboard-`, rand.Intn(10), `",
		"status_code": `, statusCode, `,
		"response_time": `, rand.Intn(300) + 20, `,
		"ip_address": "10.0.0.1",
		"user_agent": ""
	}`)
}

func main() {
	Conn, _ := net.DialUDP("udp", nil, &net.UDPAddr{IP: []byte{127, 0, 0, 1}, Port: 1337, Zone: ""})

	defer Conn.Close()

	// Send 8 messages every 10 ms. This translates to 8 * 100 * 60 = 48,000 RPM
	limiter := time.Tick(10 * time.Millisecond)

	for {
		<-limiter

		for i := 1; i < 8; i++ {
			_, _ = Conn.Write([]byte(getRandomWebRequest()))
		}
	}
}
