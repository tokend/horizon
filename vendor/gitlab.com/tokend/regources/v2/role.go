package regources

type AccountRole struct {
	Key
	Attributes    AccountRoleAttrs `json:"attributes"`
	Relationships RoleRelation     `json:"relationships"`
}

type AccountRoleAttrs struct {
	Details Details `json:"details"`
}

type RoleRelation struct {
	Rules *RelationCollection `json:"rules"`
}

type SignerRoleRelation struct {
	Rules *RelationCollection `json:"rules"`
	Owner *Relation           `json:"owner"`
}

type SignerRole struct {
	Key
	Attributes    SignerRoleAttrs    `json:"attributes"`
	Relationships SignerRoleRelation `json:"relationships"`
}

type SignerRoleAttrs struct {
	Details Details `json:"details"`
}
