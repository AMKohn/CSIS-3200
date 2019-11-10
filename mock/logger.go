package main

import (
	"net"
)

func main() {
	Conn, _ := net.DialUDP("udp", nil, &net.UDPAddr{IP: []byte{127, 0, 0, 1}, Port: 1337, Zone: ""})

	defer Conn.Close()

	for i := 1; i < 101000; i++ {
		Conn.Write([]byte(`{ "type": "web_request", "server_id": "string", "request_type": "GET/POST/...", "url": "", "status_code": 400, "response_time": 1000, "ip_address": "10.0.0.1", "user_agent": "" }`))
	}
}
