package dashboard

import (
	"csis3200/internal/app/processor"
	"encoding/json"
	"math"
	"net/http"
	"sort"
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
	} else {
		cacheHitRate = reverseProxy / webRequests
	}

	return map[string]interface{}{
		"webRequests":      webRequests,
		"databaseQueries":  databaseQueries,
		"searchQueries":    searchQueries,
		"cacheHitRate":     cacheHitRate,
		"liveServers":      calculateLiveServers(data),
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
	type host struct {
		Name   string
		EventCount int
		TotalRequests int
		ErrorRequests int
		TotalResponseTime int
		OldestRequestTime int64
		NewestRequestTime int64
		LastSystemEvent map[string]interface{}
	}

	// Use pointers to hosts to allow updating the struct, see https://stackoverflow.com/a/32751792/900747
	var hostData = map[string]*host{}

	// Count the number of events for each host and put them into the map
	for _, e := range data {
		if _, d := hostData[e["server_id"].(string)]; !d {
			hostData[e["server_id"].(string)] = &host{
				e["server_id"].(string),
				0,
				0,
				0,
				0,
				0,
				0,
				map[string]interface{}{ // Default
					"type": "system",
					"timestamp": int64(0),
					"server_id": e["server_id"].(string),
					"cpu_usage": 0.0,
					"memory_usage": 0.0,
					"memory_capacity": 0.0,
				},
			}
		}

		hostData[e["server_id"].(string)].EventCount++
	}

	// Copy the pointers to a slice so we can sort them and get the top entries
	var ss []*host
	for _, v := range hostData {
		ss = append(ss, v)
	}

	sort.Slice(ss, func(i, j int) bool {
		return ss[i].EventCount > ss[j].EventCount
	})

	// Only keep the top five hosts
	if len(ss) > 5 {
		ss = ss[:5]
	}

	// And put them back into the hostData map so we can update them efficiently
	hostData = map[string]*host{}

	for _, host := range ss {
		hostData[host.Name] = host
	}
	
	// Go through the events and calculate the relevant statistics
	for _, e := range data {
		hData, exists := hostData[e["server_id"].(string)]

		// Skip the event if it's not from one of our top hosts
		if !exists {
			continue
		}

		switch e["type"] {
			case "web_request":
				hData.TotalRequests++
				hData.TotalResponseTime += int(e["response_time"].(float64))

				// This matches unsuccessful requests
				if e["status_code"].(float64) >= 400 && e["status_code"].(float64) < 600 {
					hData.ErrorRequests++
				}

				if hData.OldestRequestTime > e["timestamp"].(int64) {
					hData.OldestRequestTime = e["timestamp"].(int64)
				} else if hData.NewestRequestTime < e["timestamp"].(int64) {
					hData.NewestRequestTime = e["timestamp"].(int64)
				}

				break

			case "system":
				if hData.LastSystemEvent["timestamp"].(int64) <= e["timestamp"].(int64) {
					hData.LastSystemEvent = e
				}

				break
		}
	}

	var retData []interface{}

	for _, h := range hostData {
		var minuteRange = int((h.NewestRequestTime - h.OldestRequestTime) / 1000 / 60)

		retData = append(retData, map[string]interface{}{
			"name": h.Name,
			"errorRate": math.Round((float64(h.ErrorRequests) / float64(h.TotalRequests)) * 10000) / 100,
			"memoryUsage": h.LastSystemEvent["memory_usage"].(float64),
			"memoryCapacity": h.LastSystemEvent["memory_capacity"].(float64),
			"cpuUsage": h.LastSystemEvent["cpu_usage"].(float64),
			"throughput": h.TotalRequests / minuteRange,
			"responseTime": h.TotalResponseTime / h.TotalRequests,
		})
	}

	return retData
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
		"messagesLast30":len(data),
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
	if requests == 0{
		return 0
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
	var servers = map[string]bool{}

	for _, e := range data {
		servers[e["server_id"].(string)] = true
	}

	return len(servers)
}
