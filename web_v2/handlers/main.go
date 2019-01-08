package handlers

import (
	"gitlab.com/tokend/horizon/web_v2/resource"
	"net/http"
)

type Resource interface {
	Prepare(r *http.Request) error
	IsAllowed() (bool, error)
	Fetch(id string) error
	Populate() error
	Response() (interface{}, error)
}

type Collection interface {
	Prepare(r *http.Request) error
	IsAllowed() (bool, error)
	Fetch(resource.PagingParams) error
	Populate() error
	Response() ([]interface{}, error)
}
