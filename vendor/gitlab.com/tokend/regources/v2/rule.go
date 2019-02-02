package regources

type Rule struct {
	Key
	Attributes RuleAttr `json:"attributes"`
}

type RuleAttr struct {
	Resource    string                 `json:"resource"`
	Action      string                 `json:"action"`
	IsForbidden string                 `json:"is_forbidden"`
	Details     map[string]interface{} `json:"details"`
}
