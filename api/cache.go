package api

import (
	"encoding/json"
	"net/http"
)

func getCacheInfoHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	cache := getCache(ctx)
	if cache == nil {
		jsonError(w, "cache has not been initialized", http.StatusInternalServerError)
		return
	}

	info, err := cache.Info()
	if err != nil {
		jsonError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response := map[string]interface{}{"info": info}
	jsonBytes, err := json.Marshal(response)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonBytes)
}

func purgeCacheHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	cache := getCache(ctx)
	if cache == nil {
		jsonError(w, "cache has not been initialized", http.StatusInternalServerError)
		return
	}

	err := cache.Purge()
	if err != nil {
		jsonError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
