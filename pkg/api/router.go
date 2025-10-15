package api

import (
	"net/http"
)

func NewRouter() *http.ServeMux {
	r := http.NewServeMux()
	r.HandleFunc("POST /{broker}/topics", CreateTopicHandler)
	r.HandleFunc("GET /{broker}/topics", ListTopicsHandler)
	r.HandleFunc("DELETE /{broker}/topics/{name}", DeleteTopicHandler)
	return r
}
