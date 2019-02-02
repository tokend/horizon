package regources

type Role struct {
	Key
	Attributes    RoleAsstr    `json:"attributes"`
	Relationships RoleRelation `json:"relationships"`
}

type RoleAsstr struct {
	Details map[string]interface{} `json:"details"`
}

type RoleRelation struct {
	Rules *RelationCollection `json:"rules"`
}
