package resource

type Resource interface {
	IsAllowed() (bool, error)
	FindOwner() error
	PopulateAttributes() error
}

const (
	TypeAccounts = "accounts"
)
