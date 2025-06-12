package api

import (
	"github.com/gorilla/mux"
)

func NewRouter() *mux.Router {
	r := mux.NewRouter()
	r.HandleFunc("{broker}/topics", CreateTopicHandler).Methods("POST")
	r.HandleFunc("{broker}/topics", ListTopicsHandler).Methods("GET")
	r.HandleFunc("{broker}/topics/{name}", DeleteTopicHandler).Methods("DELETE")
	return r
}
