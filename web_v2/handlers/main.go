package handlers

type Resource interface {
	IsAllowed() (bool, error)
	FindOwner() error
	Fetch() error
	PopulateAttributes() error
	Response() (interface{}, error)
}
