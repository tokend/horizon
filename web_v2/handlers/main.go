package handlers

type Resource interface {
	IsAllowed() (bool, error)
	Fetch() error
	Populate() error
	Response() (interface{}, error)
}
