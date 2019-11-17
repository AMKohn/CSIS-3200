package public

import (
	"fmt"
	"math/rand"
	"strconv"
	"time"
)

func getWebRequest(ts int64) map[string]interface{} {
	var statusCode = 200

	// Only 10% of requests have non-200 status codes
	if rand.Intn(10) == 1 {
		statusCode = rand.Intn(5) * 100
	}

	return map[string]interface{}{
		"timestamp": ts,
		"type": "web_request",
		"server_id": "app-" + strconv.Itoa(rand.Intn(60) + 75),
		"request_type": "GET",
		"path": "/dashboard-" + strconv.Itoa(rand.Intn(10)),
		"status_code": float64(statusCode),
		"response_time": float64(rand.Intn(300) + 20),
		"ip_address": "10.0.0.1",
		"user_agent": "",
	}
}

// Gets the requests to pre-populate the processor messagesDb
func GetInitRequests(minutes int) []map[string]interface{} {
	var start = time.Now()

	var retSlice []map[string]interface{}

	var startTs = time.Now().Add(time.Duration(-minutes) * time.Minute).UnixNano() / 1000000
	var endTs = time.Now().UnixNano() / 1000000

	// Messages are normally sent 8 every 10 ms, this simulates that for 30 minutes of data
	for ts := startTs; ts < endTs; ts += 10 {
		retSlice = append(retSlice,
			getWebRequest(ts),
			getWebRequest(ts),
			getWebRequest(ts),
			getWebRequest(ts),
			getWebRequest(ts),
			getWebRequest(ts),
			getWebRequest(ts),
			getWebRequest(ts))
	}

	fmt.Printf("Generating %d minutes of mock messages took %v\n", minutes, time.Since(start))

	return retSlice
}