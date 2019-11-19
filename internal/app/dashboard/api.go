package dashboard

import (
	"csis3200/internal/app/processor"
	"encoding/json"
	"fmt"
	"math"
	"net/http"
	"sort"
	"strconv"
	"time"
)

/**
 * Calculates a number of general statistics. These are combined into a single function so they
 * can share a single loop, which greatly increases efficiency
 */
func getStats(data []map[string]interface{}) map[string]interface{} {
	var servers = map[string]bool{}

	var totalErrors = 0.0
	var totalCPU = 0.0
	var totalResponseTime = 0

	var counts = map[string]float64 {
		"system": 0,
		"web_request": 0,
		"database": 0,
		"search": 0,
		"reverse_proxy": 0,
	}

	// Sum the statistics we'll need using a combined loop
	for _, i := range data {
		counts[i["type"].(string)]++
		servers[i["server_id"].(string)] = true

		// Average CPU
		if i["type"] == "system" && i["cpu_usage"] != nil {
			totalCPU += i["cpu_usage"].(float64)
		}

		// Average Response Time
		if i["type"] == "web_request" {
			totalResponseTime += int(i["response_time"].(float64))

			// Error rate
			if i["status_code"].(float64) >= 400 {
				totalErrors++
			}
		}
	}

	//Cache Hit Rate
	var cacheHitRate = 0.0
	if counts["web_request"] + counts["reverse_proxy"] > 0 {
		cacheHitRate = math.Round(counts["reverse_proxy"] / (counts["reverse_proxy"] + counts["web_request"]) * 10000) / 100
	}

	// Average response time and error rate
	var errorRate = 0.0
	var averageResponse = 0
	if counts["web_request"] > 0 {
		averageResponse = totalResponseTime / int(counts["web_request"])
		errorRate = math.Round(totalErrors / counts["web_request"] * 10000) / 100
	}

	//Average CPU
	var averageCPU = 0.0
	if counts["system"] > 0.0 {
		averageCPU = math.Round((totalCPU / counts["system"]) * 10) / 10
	}

	// Messages per second
	var perSecond = 0.0
	if len(data) > 0 {
		perSecond = float64(len(data)) / (float64(data[len(data) - 1]["timestamp"].(int64) - data[0]["timestamp"].(int64)) / 1000)
	}

	return map[string]interface{}{
		"webRequests":      counts["web_request"],
		"databaseQueries":  counts["database"],
		"searchQueries":    counts["search"],
		"cacheHitRate":     cacheHitRate,
		"liveServers":      len(servers),
		"cpuUsage":         averageCPU,
		"messageRate":      perSecond,
		"webResponseTime":  averageResponse,
		"overallErrorRate": errorRate,
	}
}

/**
 * Returns the slowest 7 web requests
 */
func getWebRequests(data []map[string]interface{}) []interface{} {
	type request struct {
		Signature string
		TotalRequests int
		TotalResponseTime int
		AverageResponseTime int
	}

	// Use pointers to requests to allow updating the struct, see https://stackoverflow.com/a/32751792/900747
	var requests = map[string]*request{}

	// Sum the statistics for each request path
	for _, e := range data {
		if e["type"] == "web_request" {
			signature := e["request_type"].(string) + " " + e["path"].(string)

			// Create a new request struct if it doesn't exist
			if _, d := requests[signature]; !d {
				requests[signature] = &request{
					signature,
					0,
					0,
					0,
				}
			}

			requests[signature].TotalRequests++
			requests[signature].TotalResponseTime += int(e["response_time"].(float64))
		}
	}

	// Go through the events and calculate the average response time
	for _, e := range requests {
		e.AverageResponseTime = e.TotalResponseTime / e.TotalRequests
	}

	// Copy the pointers to a slice so we can sort them and get the top entries
	var ss []*request
	for _, v := range requests {
		ss = append(ss, v)
	}

	sort.Slice(ss, func(i, j int) bool {
		return ss[i].AverageResponseTime > ss[j].AverageResponseTime
	})

	// Only keep the top seven slowest requests
	if len(ss) > 7 {
		ss = ss[:7]
	}

	// Generate and return the final JSON-compatible output
	var retData []interface{}

	for _, r := range ss {
		retData = append(retData, map[string]interface{}{
			"path": r.Signature,
			"time": strconv.Itoa(r.AverageResponseTime) + " ms",
		})
	}

	return retData
}

/**
 * Gets the top 5 highest load hosts
 */
func getHosts(data []map[string]interface{}) []map[string]interface{} {
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

	// Calculate the statistics for all hosts. This uses a single loop instead of one to find
	// the top hosts and one to sum the statistics and is faster than just calculating the
	// statistics for the top hosts
	for _, e := range data {
		hData, exists := hostData[e["server_id"].(string)]

		if !exists {
			hostData[e["server_id"].(string)] = &host{}

			hData = hostData[e["server_id"].(string)]

			hData.Name = e["server_id"].(string)
		}

		hData.EventCount++

		// The events are sorted, so we don't need to confirm that this system event is newer
		if e["type"] == "system" {
			hData.LastSystemEvent = e
		} else if e["type"] == "web_request" {
			hData.TotalRequests++
			hData.TotalResponseTime += int(e["response_time"].(float64))

			// This matches unsuccessful requests
			if e["status_code"].(float64) >= 400 && e["status_code"].(float64) < 600 {
				hData.ErrorRequests++
			}

			if hData.OldestRequestTime == 0 || hData.OldestRequestTime > e["timestamp"].(int64) {
				hData.OldestRequestTime = e["timestamp"].(int64)
			} else if hData.NewestRequestTime == 0 || hData.NewestRequestTime < e["timestamp"].(int64) {
				hData.NewestRequestTime = e["timestamp"].(int64)
			}
		}
	}

	// Copy the pointers to a slice so we can sort them and get the top entries
	var ss []*host
	for _, v := range hostData {
		ss = append(ss, v)
	}

	// Sort the hosts by CPU usage
	sort.Slice(ss, func(i, j int) bool {
		return ss[i].LastSystemEvent["cpu_usage"].(float64) > ss[j].LastSystemEvent["cpu_usage"].(float64)
	})

	// Only keep the top five hosts
	if len(ss) > 5 {
		ss = ss[:5]
	}

	// Generate the final JSON-compatible output and return
	var retData []map[string]interface{}

	for _, h := range ss {
		var minuteRange = float64(h.NewestRequestTime - h.OldestRequestTime) / 1000 / 60

		var throughput, responseTime, errorRate = 0, 0, 0.0

		if minuteRange != 0 {
			throughput = int(float64(h.TotalRequests) / minuteRange)
		}

		if h.TotalRequests != 0 {
			responseTime = h.TotalResponseTime / h.TotalRequests
			errorRate = math.Round((float64(h.ErrorRequests) / float64(h.TotalRequests)) * 10000) / 100
		}

		retData = append(retData, map[string]interface{}{
			"name": h.Name,
			"errorRate": errorRate,
			"memoryUsage": h.LastSystemEvent["memory_usage"].(float64) * 1000,
			"memoryCapacity": h.LastSystemEvent["memory_capacity"].(float64) * 1000,
			"cpuUsage": h.LastSystemEvent["cpu_usage"].(float64),
			"throughput": throughput,
			"responseTime": responseTime,
		})
	}

	return retData
}

/**
 * Returns the average response times in 5 minute periods
 */
func getResponseTimes(data []map[string]interface{}) [7]int {
	var times [7]int
	var counts [7]int

	var startTs = time.Now().Add(time.Duration(-30) * time.Minute).UnixNano() / 1000000

	for _, e := range data {
		// Buckets are (-32.5 m, -27.5 m), (-27.5 m, -22.5 m), etc. This calculates the bucket
		// index based on the time differential
		bucket := int(math.Round((float64(e["timestamp"].(int64) - startTs) / 1000 / 60) / 5))

		if bucket >= 0 && bucket <= 6 &&  e["type"] == "web_request" {
			times[bucket] += int(e["response_time"].(float64))
			counts[bucket]++
		}
	}

	var averages [7]int

	for i := 0; i < 7; i++ {
		if counts[i] > 0 {
			averages[i] = times[i] / counts[i]
		}
	}

	return averages
}

/**
 * Handles the /api web requests. Those all return the JSON dashboard data
 */
func apiHandler(w http.ResponseWriter, r *http.Request) {
	// Set the response content type
	w.Header().Set("Content-Type", "application/json")

	// Get the last 30 minutes of data and pass that to each of the calculation functions
	// This takes a negligible amount of time to run
	data := processor.GetRecentData()

	jsonData := map[string]interface{}{
		"messagesLast30": len(data),
	}

	var statsStart = time.Now()

	jsonData["stats"] = getStats(data)
	statsTime := time.Since(statsStart)

	jsonData["webRequests"] = getWebRequests(data)
	webTime := time.Since(statsStart) - statsTime

	jsonData["hosts"] = getHosts(data)
	hostsTime := time.Since(statsStart) - statsTime - webTime

	jsonData["responseTimes"] = getResponseTimes(data)
	respTime := time.Since(statsStart) - statsTime - webTime - hostsTime

	fmt.Printf(
		"Dashboard API stats compiling took %v for %d messages. %v for main stats, %v for web requests, %v for hosts, %v for response times\n",
		time.Since(statsStart), jsonData["messagesLast30"].(int), statsTime, webTime, hostsTime, respTime)

	// NewEncoder encodes and sends the message
	_ = json.NewEncoder(w).Encode(jsonData)
}
