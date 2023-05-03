package api

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/infamous55/go-zestful/cache"
)

type contextKey string

const (
	cacheKey contextKey = "cache"
)

func getCache(ctx context.Context) cache.Cache {
	if cache, ok := ctx.Value(cacheKey).(cache.Cache); ok {
		return cache
	}
	return nil
}

func GenerateCacheMiddleware(cache cache.Cache) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(
			func(w http.ResponseWriter, r *http.Request) {
				ctx := context.WithValue(r.Context(), cacheKey, cache)
				next.ServeHTTP(w, r.WithContext(ctx))
			},
		)
	}
}

func jsonError(w http.ResponseWriter, message string, statusCode int) {
	errorResponse := map[string]string{"error": message}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(errorResponse)
}
