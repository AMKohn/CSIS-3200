package public

import "math/rand"

type DBServer struct {
	BaseServer
}

var objectTypes = [4]string{"user", "page", "order", "product"}

var actionTypes = [4]string{"create", "read", "update", "delete"}

func (d DBServer) GetLogMessage(ts int64) map[string]interface{} {
	var actionType = "GET"

	// 20% are randomly picked
	if rand.Intn(5) == 1 {
		actionType = actionTypes[rand.Intn(4)]
	}

	return map[string]interface{}{
		"timestamp": ts,
		"type": "database",
		"server_id": d.ID,
		"action_type": actionType,
		"object_type": objectTypes[rand.Intn(4)],
		"object_id": float64(rand.Intn(10000)),
		"retrieval_time": float64(rand.Intn(d.SpeedFactor) / 10),
	}
}