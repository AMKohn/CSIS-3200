package dashboard

import (
	"csis3200/internal/app/processor"
	"encoding/json"
	"net/http"
	"time"
)

func getStats(data []map[string]interface{}) map[string]interface{} {
	var webRequests = 0
	var databaseQueries = 0
	var searchQueries = 0
	var reverseProxy = 0
	var cacheHitRate = 0

	for _, i := range data {
		if i["type"] == "web_request"{
			webRequests++
		}
		if i["type"] == "database"{
			databaseQueries++
		}
		if i["type"] == "search"{
			searchQueries++
		}
		if i["type"] == "reverse_proxy"{
			reverseProxy++
		}
	}

	if webRequests == 0{
		cacheHitRate = 0
	}
	cacheHitRate = reverseProxy / webRequests

	return map[string]interface{}{
		"webRequests":      webRequests,
		"databaseQueries":  databaseQueries,
		"searchQueries":    searchQueries,
		"cacheHitRate":     cacheHitRate,
		"liveServers":      37,
		"cpuUsage":         averageCPU(data),
		"messageRate":      msgPerSec(data),
		"webResponseTime":  averageResponseTimes(data),
		"overallErrorRate": averageErrorRate(data),
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
		inRange := false

		// If we were passed a valid start time, use that
		if startTime > 0 && endTime > 0 {
			inRange = e["timestamp"].(int64) > startTime && e["timestamp"].(int64) <= endTime
		}

		if e["type"] == "web_request" && inRange {
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

func averageResponseTimes(data []map[string]interface{}) int{
	var time = 0
	var requests = 0
	for _, i := range data {
		if i["type"] == "webRequest"{
			time += int(i["responseTime"].(float64))
			requests++
		}
	}
	return time / requests
}


func averageCPU(data []map[string]interface{}) float64{
	var totalCPU = 0.0
	var totalRequests = 0.0
	for _, i := range data{
		totalCPU += i["cpuUsage"].(float64)
		totalRequests ++
	}
	if totalRequests > 0.0 {
		return totalCPU/totalRequests * 100
	}
	return 0.0
}

func msgPerSec(data []map[string]interface{}) float64{

	return float64(len(data)) / 1800
}

func averageErrorRate(data []map[string]interface{}) float64{
	var totalError = 0.0
	var totalMessage = 0.0

	for _, i := range data {
		if i["type"] == "webRequest"{
			if i["status_code"].(float64) >= 400 {
				totalError++
			}
			totalMessage++
		}
	}

	if totalMessage == 0{
		return 0.0
	}

	return totalError / totalMessage
}


func calculateLiveServers(data []map[string]interface{}) int{
	var servers map[string]bool

	

	return len(servers)
}
