package api

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/gorilla/mux"
)

func getItemHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	cache := getCache(ctx)
	if cache == nil {
		jsonError(w, "cache has not been initialized", http.StatusInternalServerError)
		return
	}

	vars := mux.Vars(r)
	key := vars["key"]
	if key == "" {
		jsonError(w, "invalid key", http.StatusBadRequest)
		return
	}

	value, err := cache.Get(key)
	if err != nil {
		jsonError(w, err.Error(), http.StatusNotFound)
		return
	}

	jsonBytes, err := json.Marshal(value)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonBytes)
}

func updateItemHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	cache := getCache(ctx)
	if cache == nil {
		jsonError(w, "cache has not been initialized", http.StatusInternalServerError)
		return
	}

	vars := mux.Vars(r)
	key := vars["key"]
	if key == "" {
		jsonError(w, "invalid key", http.StatusBadRequest)
		return
	}

	_, err := cache.Get(key)
	if err != nil {
		jsonError(w, err.Error(), http.StatusNotFound)
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		jsonError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var value interface{}
	err = json.Unmarshal(body, &value)
	if err != nil {
		jsonError(w, err.Error(), http.StatusBadRequest)
		return
	}

	cache.Set(key, value)
	w.WriteHeader(http.StatusNoContent)
}

func deleteItemHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	cache := getCache(ctx)
	if cache == nil {
		jsonError(w, "cache has not been initialized", http.StatusInternalServerError)
		return
	}

	vars := mux.Vars(r)
	key := vars["key"]
	if key == "" {
		jsonError(w, "invalid key", http.StatusBadRequest)
		return
	}

	err := cache.Delete(key)
	if err != nil {
		jsonError(w, err.Error(), http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
