package resource

type ResourceType string

type PagingParams struct {
	page  string
	limit string
}

const (
	TypeAccounts = "accounts"
)
