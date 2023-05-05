package api

import (
	"github.com/gorilla/mux"
)

func RegisterItemsHandlers(subrouter *mux.Router) {
	subrouter.StrictSlash(false)
	subrouter.HandleFunc("/{key}/", getItemHandler).Methods("GET")
	subrouter.HandleFunc("/", createItemHandler).Methods("POST")
	subrouter.HandleFunc("/{key}/", updateItemHandler).Methods("PUT")
	subrouter.HandleFunc("/{key}/", deleteItemHandler).Methods("DELETE")
}

func RegisterAuthHandlers(subrouter *mux.Router, secret string, key []byte) {
	subrouter.StrictSlash(false)
	subrouter.HandleFunc("/token/", createTokenHandler(secret, key)).Methods("POST")
	subrouter.HandleFunc("/refresh/", refreshTokenHandler(secret, key)).Methods("POST")
}
