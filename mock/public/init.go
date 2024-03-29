package public

import (
	"fmt"
	"time"
)

/**
 * Gets the requests to pre-populate the processor messagesDb
 */
func GetInitRequests(minutes int) []map[string]interface{} {
	// Get a server cluster
	var servers = GetServerCluster()
	var start = time.Now()

	var retSlice []map[string]interface{}

	// Determine the range we're running over
	var startTs = time.Now().Add(time.Duration(-minutes) * time.Minute).UnixNano() / 1000000
	var endTs = time.Now().UnixNano() / 1000000

	var lastTs = int64(0)

	// Simulate the message sending that the logger would do
	for ts := startTs; ts < endTs; ts += 100 {
		for _, s := range servers {
			var rate = s.GetRequestRate()

			// Convert RPM to the right number to send every 100 ms
			for i := 0; i < rate / 60 / 10; i++ {
				retSlice = append(retSlice, s.GetLogMessage(ts))
			}
		}

		// If it's been more than 1 second, send a health ping for all servers too
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