package resource

type Resource interface {
	IsAllowed() (bool, error)
	FindOwner() error
	PopulateAttributes() error
	Response() (interface{}, error)
}

type ResourceType string

const (
	TypeAccounts = "accounts"
)
