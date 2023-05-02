package api

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
)

func getItemHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	key := vars["key"]

	w.Header().Set("Content-Type", "application/json")

	if key == "" {
		errorResponse := map[string]string{"error": "invalid key"}
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(errorResponse)
		return
	}

	ctx := r.Context()
	cache := getCache(ctx)
	if cache == nil {
		http.Error(w, "cache has not been initialized", http.StatusInternalServerError)
		return
	}

	value, err := cache.Get(key)
	if err != nil {
		errorResponse := map[string]string{"error": err.Error()}
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(errorResponse)
		return
	}

	jsonBytes, err := json.Marshal(value)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Write(jsonBytes)
}

func deleteItemHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	key := vars["key"]

	if key == "" {
		errorResponse := map[string]string{"error": "invalid key"}
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(errorResponse)
		return
	}

	ctx := r.Context()
	cache := getCache(ctx)
	if cache == nil {
		http.Error(w, "cache has not been initialized", http.StatusInternalServerError)
		return
	}

	err := cache.Delete(key)
	if err != nil {
		errorResponse := map[string]string{"error": err.Error()}
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(errorResponse)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
