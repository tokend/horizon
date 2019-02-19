package regources

type AccountRule struct {
	Key
	Attributes AccountRuleAttr `json:"attributes"`
}

type AccountRuleAttr struct {
	Resource string  `json:"resource"`
	Action   string  `json:"action"`
	IsForbid bool    `json:"is_forbid"`
	Details  Details `json:"details"`
}

type SignerRule struct {
	Key
	Attributes SignerRuleAttr `json:"attributes"`
}

type SignerRuleAttr struct {
	Resource  string  `json:"resource"`
	Action    string  `json:"action"`
	IsForbid  bool    `json:"is_forbid"`
	IsDefault bool    `json:"is_default"`
	OwnerID   string  `json:"owner_id"`
	Details   Details `json:"details"`
}
