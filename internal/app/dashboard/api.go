package dashboard

import (
	"encoding/json"
	"net/http"
)

func apiHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	jsonData := map[string]interface{}{
		"stats": map[string]interface{}{
			"webRequests":      8300000,
			"databaseQueries":  530610,
			"searchQueries":    150300,
			"cacheHitRate":     76,
			"liveServers":      37,
			"cpuUsage":         63,
			"messageRate":      1300,
			"webResponseTime":  85,
			"overallErrorRate": 0.2,
		},
		"webRequests": []interface{}{map[string]interface{}{
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
		}},
		"hosts": []interface{}{map[string]interface{}{
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
		}},
		"responseTimes": []interface{}{100, 78, 34, 46, 106, 99, 164},
	}

	_ = json.NewEncoder(w).Encode(jsonData)
}
