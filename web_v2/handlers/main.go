package handlers

import "net/http"

type Handler interface {
	Prepare(w http.ResponseWriter, r *http.Request)
	Serve(w http.ResponseWriter, r *http.Request)
}

type Resource interface {
	IsAllowed() (bool, error)
	FindOwner() error
	PopulateAttributes() error
	Response() (interface{}, error)
}
