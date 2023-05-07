package api

import (
	"encoding/json"
	"io"
	"net/http"
	"time"

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

	response := map[string]interface{}{"value": value}
	jsonBytes, err := json.Marshal(response)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonBytes)
}

type createItemBody struct {
	Key        string      `json:"key"`
	TimeToLive *string     `json:"ttl,omitempty"`
	Value      interface{} `json:"value"`
}

func createItemHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	cache := getCache(ctx)
	if cache == nil {
		jsonError(w, "cache has not been initialized", http.StatusInternalServerError)
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		jsonError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var newItem createItemBody
	err = json.Unmarshal(body, &newItem)
	if err != nil {
		jsonError(w, err.Error(), http.StatusBadRequest)
		return
	}

	if newItem.Key == "" {
		jsonError(w, "invalid key", http.StatusBadRequest)
		return
	}

	if newItem.Value == nil {
		jsonError(w, "invalid value", http.StatusBadRequest)
		return
	}

	existingItem, _ := cache.Get(newItem.Key)
	if existingItem != nil {
		jsonError(w, "item already exists", http.StatusConflict)
		return
	}

	var ttl time.Duration
	if newItem.TimeToLive != nil {
		ttl, err = time.ParseDuration(*newItem.TimeToLive)
		if err != nil {
			jsonError(w, "invalid time-to-live", http.StatusBadRequest)
			return
		}
		cache.Set(newItem.Key, newItem.Value, ttl)
	} else {
		cache.Set(newItem.Key, newItem.Value)
	}

	w.WriteHeader(http.StatusNoContent)
}

type updateItemBody struct {
	TimeToLive *string     `json:"ttl,omitempty"`
	Value      interface{} `json:"value"`
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

	var updatedItem updateItemBody
	err = json.Unmarshal(body, &updatedItem)
	if err != nil {
		jsonError(w, err.Error(), http.StatusBadRequest)
		return
	}

	if updatedItem.Value == nil {
		jsonError(w, "invalid value", http.StatusBadRequest)
		return
	}

	var ttl time.Duration
	if updatedItem.TimeToLive != nil {
		ttl, err = time.ParseDuration(*updatedItem.TimeToLive)
		if err != nil {
			jsonError(w, "invalid time-to-live", http.StatusBadRequest)
			return
		}
		cache.Set(key, updatedItem.Value, ttl)
	} else {
		cache.Set(key, updatedItem.Value)
	}

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
