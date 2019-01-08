package resource

type ResourceType string

type PagingParams struct {
	Page  string
	Limit string
}

const (
	TypeAccounts = "accounts"
)
