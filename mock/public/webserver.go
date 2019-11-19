package public

import "math/rand"

type WebServer struct {
	BaseServer
}

var requestTypes = [4]string{"GET", "PUT", "POST", "DELETE"}

func (w WebServer) GetLogMessage(ts int64) map[string]interface{} {
	var statusCode = 200
	var requestType = "GET"

	// Only 10% of requests have non-200 status codes
	if rand.Intn(10) == 1 {
		statusCode = rand.Intn(5) * 100
	}

	// 20% of requests have a non-GET type
	if rand.Intn(5) == 1 {
		requestType = requestTypes[rand.Intn(4)]
	}

	return map[string]interface{}{
		"timestamp": ts,
		"type": "web_request",
		"server_id": w.ID,
		"request_type": requestType,
		"path": urlPaths[rand.Intn(7)],
		"status_code": float64(statusCode),
		"response_time": float64(rand.Intn(w.SpeedFactor / 2) + (w.SpeedFactor / 10)) + float64((ts / 1000 / 60) % 10) * 10,
		"ip_address": "10.0.0.1",
		"user_agent": "",
	}
}