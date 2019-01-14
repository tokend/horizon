package resource

type resourceType string

const (
	typeAccounts resourceType = "accounts"
	typeBalances resourceType = "balances"
	typeAssets   resourceType = "assets"
)

// Key - identifier of the resource
type Key struct {
	ID   string       `json:"id"`
	Type resourceType `json:"type"`
}

//GetKey - returns key of the resource
func (r *Key) GetKey() Key {
	return *r
}
