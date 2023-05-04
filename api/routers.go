package api

import (
	"github.com/gorilla/mux"
)

func NewItemsRouter() *mux.Router {
	router := mux.NewRouter()
	router.HandleFunc("/{key}", getItemHandler).Methods("GET")
	router.HandleFunc("/", createItemHandler).Methods("POST")
	router.HandleFunc("/{key}", updateItemHandler).Methods("PUT")
	router.HandleFunc("/{key}", deleteItemHandler).Methods("DELETE")
	return router
}

func NewAuthRouter(secret string, key []byte) *mux.Router {
	router := mux.NewRouter()
	router.HandleFunc("/token", createTokenHandler(secret, key)).Methods("POST")
	router.HandleFunc("/refresh", refreshTokenHandler(secret, key)).Methods("POST")
	return router
}
