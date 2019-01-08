package handlers

import "gitlab.com/tokend/horizon/web_v2/resource"

type Resource interface {
	IsAllowed() (bool, error)
	Fetch(id string) error
	Populate() error
	Response() (interface{}, error)
}

type Collection interface {
	IsAllowed() (bool, error)
	Fetch(resource.PagingParams) error
	Populate() error
	Response() ([]interface{}, error)
}
