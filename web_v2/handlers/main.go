package handlers

type Resource interface {
	IsAllowed() (bool, error)
	FindOwner() error
	PopulateAttributes() error
	Response() (interface{}, error)
}
