package api

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v5"
	"github.com/infamous55/go-zestful/cache"
)

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

func GenerateAuthMiddleware(key []byte) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(
			func(w http.ResponseWriter, r *http.Request) {
				authHeader := r.Header.Get("Authorization")
				if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
					jsonError(w, "invalid authorization header", http.StatusUnauthorized)
					return
				}

				tokenString := strings.TrimPrefix(authHeader, "Bearer ")
				_, err := jwt.ParseWithClaims(tokenString, &jwt.RegisteredClaims{}, func(token *jwt.Token) (interface{}, error) {
					if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
						return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
					}
					return key, nil
				})
				if err != nil {
					jsonError(w, err.Error(), http.StatusUnauthorized)
					return
				}

				next.ServeHTTP(w, r)
			},
		)
	}
}
