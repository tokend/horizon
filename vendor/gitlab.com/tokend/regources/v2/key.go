package regources

type ResourceType string

const (
	TypeAccounts      ResourceType = "accounts"
	TypeBalances                   = "balances"
	TypeAssets                     = "assets"
	TypeBalancesState              = "balances-state"
	TypeRoles                      = "roles"
	TypeRules                      = "rules"
	TypeSigners                    = "signers"
	TypeSignerRoles                = "signer-roles"
	TypeSignerRules                = "signer-rules"
)

// Key - identifier of the Resource
type Key struct {
	ID   string       `json:"id"`
	Type ResourceType `json:"type"`
}

//GetKey - returns key of the Resource
func (r *Key) GetKey() Key {
	return *r
}

//GetKeyP - returns key pointer
func (r *Key) GetKeyP() *Key {
	return r
}

// AsRelation - converts key to relation
func (r *Key) AsRelation() *Relation {
	return &Relation{
		Data: r.GetKeyP(),
	}
}
