package api

import (
	"encoding/json"
	"net/http"
)

func jsonError(w http.ResponseWriter, message string, statusCode int) {
	errorResponse := map[string]string{"error": message}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(errorResponse)
}
