package regources

type SignerRoleRelationships struct {
	Owner *Relation           `json:"owner,omitempty"`
	Rules *RelationCollection `json:"rules,omitempty"`
}
