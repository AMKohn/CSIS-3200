package dashboard

import (
	"csis3200/internal/app/processor"
	"encoding/json"
	"net/http"
	"time"
)

func getStats(data []map[string]interface{}) map[string]interface{} {
	return map[string]interface{}{
		"webRequests":      8300000,
		"databaseQueries":  530610,
		"searchQueries":    150300,
		"cacheHitRate":     76,
		"liveServers":      37,
		"cpuUsage":         63,
		"messageRate":      1300,
		"webResponseTime":  85,
		"overallErrorRate": 0.2,
	}
}

func getWebRequests(data []map[string]interface{}) []interface{} {
	return []interface{}{map[string]interface{}{
		"path": "get /",
		"time": "312 ms",
	}, map[string]interface{}{
		"path": "post /contact",
		"time": "268 ms",
	}, map[string]interface{}{
		"path": "get /",
		"time": "312 ms",
	}, map[string]interface{}{
		"path": "post /contact",
		"time": "268 ms",
	}, map[string]interface{}{
		"path": "get /",
		"time": "312 ms",
	}, map[string]interface{}{
		"path": "post /contact",
		"time": "268 ms",
	}, map[string]interface{}{
		"path": "get /",
		"time": "312 ms",
	}}
}

func getHosts(data []map[string]interface{}) []interface{} {
	return []interface{}{map[string]interface{}{
		"name":           "app-119",
		"errorRate":      0.02,
		"memoryUsage":    3600,
		"memoryCapacity": 4000,
		"cpuUsage":       89,
		"throughput":     43400,
		"responseTime":   312,
	}, map[string]interface{}{
		"name":           "app-109",
		"errorRate":      0.02,
		"memoryUsage":    3600,
		"memoryCapacity": 4000,
		"cpuUsage":       89,
		"throughput":     43400,
		"responseTime":   312,
	}, map[string]interface{}{
		"name":           "app-109",
		"errorRate":      0.02,
		"memoryUsage":    3600,
		"memoryCapacity": 4000,
		"cpuUsage":       89,
		"throughput":     43400,
		"responseTime":   312,
	}, map[string]interface{}{
		"name":           "app-109",
		"errorRate":      0.02,
		"memoryUsage":    3600,
		"memoryCapacity": 4000,
		"cpuUsage":       89,
		"throughput":     43400,
		"responseTime":   312,
	}, map[string]interface{}{
		"name":           "app-109",
		"errorRate":      0.02,
		"memoryUsage":    3600,
		"memoryCapacity": 4000,
		"cpuUsage":       89,
		"throughput":     43400,
		"responseTime":   312,
	}}
}

func getResponseTimeRange(data []map[string]interface{}, startTime int64, endTime int64) int {
	var totalTime int = 0
	var numRequests int = 0

	for _, e := range data {
		if e["type"] == "web_request" && e["timestamp"].(int64) > startTime && e["timestamp"].(int64) <= endTime {
			totalTime += int(e["response_time"].(float64))
			numRequests++
		}
	}

	if numRequests > 0 {
		return totalTime / numRequests
	}

	return 0
}

func getResponseTimes(data []map[string]interface{}) []interface{} {
	return []interface{}{
		getResponseTimeRange(
			data,
			time.Now().Add(time.Duration(-325)*time.Minute/10).UnixNano()/1000000,
			time.Now().Add(time.Duration(-275)*time.Minute/10).UnixNano()/1000000),
		getResponseTimeRange(
			data,
			time.Now().Add(time.Duration(-275)*time.Minute/10).UnixNano()/1000000,
			time.Now().Add(time.Duration(-225)*time.Minute/10).UnixNano()/1000000),
		getResponseTimeRange(
			data,
			time.Now().Add(time.Duration(-225)*time.Minute/10).UnixNano()/1000000,
			time.Now().Add(time.Duration(-175)*time.Minute/10).UnixNano()/1000000),
		getResponseTimeRange(
			data,
			time.Now().Add(time.Duration(-175)*time.Minute/10).UnixNano()/1000000,
			time.Now().Add(time.Duration(-125)*time.Minute/10).UnixNano()/1000000),
		getResponseTimeRange(
			data,
			time.Now().Add(time.Duration(-125)*time.Minute/10).UnixNano()/1000000,
			time.Now().Add(time.Duration(-75)*time.Minute/10).UnixNano()/1000000),
		getResponseTimeRange(
			data,
			time.Now().Add(time.Duration(-75)*time.Minute/10).UnixNano()/1000000,
			time.Now().Add(time.Duration(-25)*time.Minute/10).UnixNano()/1000000),
		getResponseTimeRange(
			data,
			time.Now().Add(time.Duration(-25)*time.Minute/10).UnixNano()/1000000,
			time.Now().UnixNano()/1000000),
	}
}

func apiHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	data := processor.GetRecentData()

	jsonData := map[string]interface{}{
		"stats":         getStats(data),
		"webRequests":   getWebRequests(data),
		"hosts":         getHosts(data),
		"responseTimes": getResponseTimes(data),
	}

	_ = json.NewEncoder(w).Encode(jsonData)
}

func AvgResponseTimes(data []map[string]interface{}) int{
	var webRequests int
	for _, i := range data {
		if i["webRequests"] == 
	}

	var average int
	return average
}