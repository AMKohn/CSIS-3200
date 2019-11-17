package public

import "math/rand"

type SearchServer struct {
	BaseServer
}

var searchQueries = [10]string{"funny", "goofy", "weird", "cute", "", "little", "baby", "gray", "grey", "tabby"}

var searchIndexes = [4]string{"user", "web", "order", "product"}

func (s SearchServer) GetLogMessage(ts int64) map[string]interface{} {
	return map[string]interface{}{
		"timestamp": ts,
		"type": "search",
		"server_id": s.ID,
		"search_query": searchQueries[rand.Intn(10)] + " cat pictures",
		"search_index": searchIndexes[rand.Intn(4)],
		"query_time": float64(rand.Intn(s.SpeedFactor) / 10),
		"result_count": float64(rand.Intn(10000)),
	}
}