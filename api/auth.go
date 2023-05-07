package api

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type createTokenBody struct {
	Secret string `json:"secret"`
}

func createTokenHandler(secret string, key []byte) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		requestBody, err := io.ReadAll(r.Body)
		if err != nil {
			jsonError(w, err.Error(), http.StatusInternalServerError)
			return
		}

		var parsedBody createTokenBody
		err = json.Unmarshal(requestBody, &parsedBody)
		if err != nil {
			jsonError(w, err.Error(), http.StatusBadRequest)
			return
		}

		if parsedBody.Secret != secret {
			jsonError(w, "invalid secret", http.StatusUnauthorized)
			return
		}

		claims := &jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(20 * time.Minute)),
		}
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
		signedToken, err := token.SignedString(key)
		if err != nil {
			jsonError(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// cookie := http.Cookie{
		// 	Name:     "token",
		// 	Value:    signedToken,
		// 	MaxAge:   int(20 * time.Minute),
		// 	HttpOnly: true,
		// 	Secure:   true,
		// 	SameSite: http.SameSiteStrictMode,
		// }
		// http.SetCookie(w, &cookie)

		response := map[string]string{"token": signedToken}
		w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	}
}

func refreshTokenHandler(secret string, key []byte) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
			jsonError(w, "invalid authorization header", http.StatusUnauthorized)
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		token, err := jwt.ParseWithClaims(tokenString, &jwt.RegisteredClaims{}, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return key, nil
		})
		if err != nil {
			jsonError(w, err.Error(), http.StatusUnauthorized)
			return
		}

		if claims, ok := token.Claims.(*jwt.RegisteredClaims); ok {
			if time.Until(claims.ExpiresAt.Time) > 2*time.Minute {
				jsonError(w, "refresh attempt too early", http.StatusBadRequest)
				return
			}

			newClaims := &jwt.RegisteredClaims{
				ExpiresAt: jwt.NewNumericDate(time.Now().Add(20 * time.Minute)),
			}
			newToken := jwt.NewWithClaims(jwt.SigningMethodHS256, newClaims)
			newSignedToken, err := newToken.SignedString(key)
			if err != nil {
				jsonError(w, err.Error(), http.StatusInternalServerError)
				return
			}

			response := map[string]string{"token": newSignedToken}
			w.WriteHeader(http.StatusOK)
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(response)
		} else {
			jsonError(w, "invalid token", http.StatusUnauthorized)
		}
	}
}
