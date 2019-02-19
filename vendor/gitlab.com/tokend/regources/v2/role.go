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

type SignerRole struct {
	Key
	Attributes    SignerRoleAttrs `json:"attributes"`
	Relationships RoleRelation    `json:"relationships"`
}

type SignerRoleAttrs struct {
	Details Details `json:"details"`
	OwnerID string  `json:"owner_id"`
}
