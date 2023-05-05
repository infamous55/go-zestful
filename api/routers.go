package api

import (
	"github.com/gorilla/mux"
)

func AddItemsRoutes(subrouter *mux.Router) {
	subrouter.HandleFunc("/{key}/", getItemHandler).Methods("GET")
	subrouter.HandleFunc("/", createItemHandler).Methods("POST")
	subrouter.HandleFunc("/{key}/", updateItemHandler).Methods("PUT")
	subrouter.HandleFunc("/{key}/", deleteItemHandler).Methods("DELETE")
}

func AddAuthRoutes(subrouter *mux.Router, secret string, key []byte) {
	subrouter.HandleFunc("/token/", createTokenHandler(secret, key)).Methods("POST")
	subrouter.HandleFunc("/refresh/", refreshTokenHandler(secret, key)).Methods("POST")
}
