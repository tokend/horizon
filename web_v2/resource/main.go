package resource

type Resource interface {
	IsAllowed() (bool, error)
	FindOwner() error
	PopulateModel() error
	MarshalModel() ([]byte, error)
}

const (
	TypeAccounts = "accounts"
)
