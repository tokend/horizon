package regources

type CreateAccountOpRelationships struct {
	Account *Relation `json:"account,omitempty"`
	Role    *Relation `json:"role,omitempty"`
}
