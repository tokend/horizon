package handlers

import "net/http"

type Handler interface {
	Prepare(w http.ResponseWriter, r *http.Request)
	Serve(w http.ResponseWriter, r *http.Request)
}
