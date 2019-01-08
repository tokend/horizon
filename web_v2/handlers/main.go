package handlers

type Resource interface {
	IsAllowed() (bool, error)
	FindOwner() error
	Fetch() error
	Populate() error
	Response() (interface{}, error)
}
