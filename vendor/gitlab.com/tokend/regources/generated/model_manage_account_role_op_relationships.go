package regources

type ManageAccountRoleOpRelationships struct {
	Role  *Relation           `json:"role,omitempty"`
	Rules *RelationCollection `json:"rules,omitempty"`
}
