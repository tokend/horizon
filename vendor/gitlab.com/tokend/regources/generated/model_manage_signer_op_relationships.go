package regources

type ManageSignerOpRelationships struct {
	Role   *Relation `json:"role,omitempty"`
	Signer *Relation `json:"signer,omitempty"`
}
