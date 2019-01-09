package resource

type ResourceType string

type PagingParams struct {
	Page  string
	Limit string
}

const (
	TypeAccounts = "accounts"
)

type Response struct {
	Id            string       `json:"id"`
	Type          ResourceType `json:"type"`
	Attributes    interface{}  `json:"attributes, omitempty"`
	Relationships interface{}  `json:"relationships, omitempty"`
}
