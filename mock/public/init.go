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
	var servers = GetServerCluster()
	var start = time.Now()

	var retSlice []map[string]interface{}

	var startTs = time.Now().Add(time.Duration(-minutes) * time.Minute).UnixNano() / 1000000
	var endTs = time.Now().UnixNano() / 1000000

	var lastTs = int64(0)

	for ts := startTs; ts < endTs; ts += 100 {
		for _, s := range servers {
			var rate = s.GetRequestRate()

			// Convert RPM to the right number to send every 10 ms
			for i := 0; i < rate / 60 / 10; i++ {
				retSlice = append(retSlice, s.GetLogMessage(ts))
			}
		}

		// If it's been more than 1 second, send a health ping too
		if ts - lastTs > 1000 {
			lastTs = ts

			for _, s := range servers {
				retSlice = append(retSlice, s.GetSystemMessage(ts))
			}
		}
	}

	fmt.Printf("Generating %d messages for %d minutes of mock data took %v\n", len(retSlice), minutes, time.Since(start))

	return retSlice
}