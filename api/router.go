package api

import (
	"github.com/gorilla/mux"
)

func New() *mux.Router {
	r := mux.NewRouter()
	r.HandleFunc("/items/{key}", getItemHandler).Methods("GET")
	r.HandleFunc("/items/{key}", updateItemHandler).Methods("PUT")
	r.HandleFunc("/items/{key}", deleteItemHandler).Methods("DELETE")
	return r
}

var Router = New()
