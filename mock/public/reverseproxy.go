package public

import "math/rand"

type ReverseProxy struct {
	BaseServer
}

func (r ReverseProxy) GetLogMessage(ts int64) map[string]interface{} {
	var statusCode = 200

	// 10% of requests have non-200 status codes
	if rand.Intn(10) == 1 {
		statusCode = rand.Intn(5) * 100
	}

	return map[string]interface{}{
		"timestamp": ts,
		"type": "reverse_proxy",
		"server_id": r.ID,
		"url": urlPaths[rand.Intn(7)],
		"status_code": statusCode,
		"entry_age": float64(rand.Intn(120)),
		"response_time": float64(rand.Intn(r.SpeedFactor) / 10),
	}
}