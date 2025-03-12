package api

import (
	"github.com/gorilla/mux"
)

func NewRouter() *mux.Router {
	r := mux.NewRouter()
	r.HandleFunc("/topics", CreateTopicHandler).Methods("POST")
	r.HandleFunc("/topics", ListTopicsHandler).Methods("GET")
	r.HandleFunc("/topics/{name}", DeleteTopicHandler).Methods("DELETE")
	return r
}
