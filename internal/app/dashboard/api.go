package dashboard

import (
	"encoding/json"
	"net/http"
)

func apiHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	json.NewEncoder(w).Encode(map[string]interface{}{
		"status": 200,
	})
}