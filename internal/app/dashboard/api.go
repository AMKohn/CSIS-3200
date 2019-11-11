package dashboard

import (
	"csis3200/internal/app/processor"
	"encoding/json"
	"math/rand"
	"net/http"
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

func getResponseTimes(data []map[string]interface{}) []interface{} {
	return []interface{}{
		rand.Intn(50) * 2,
		rand.Intn(50) * 2,
		rand.Intn(50) * 2,
		rand.Intn(50) * 2,
		rand.Intn(50) * 2,
		rand.Intn(50) * 2,
		rand.Intn(50) * 2,
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
