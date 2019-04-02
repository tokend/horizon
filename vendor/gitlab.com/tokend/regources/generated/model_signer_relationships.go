package regources

type SignerRelationships struct {
	Account *Relation `json:"account,omitempty"`
	Role    *Relation `json:"role,omitempty"`
}
