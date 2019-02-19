package regources

import "gitlab.com/tokend/go/xdr"

type AccountRule struct {
	Key
	Attributes AccountRuleAttr `json:"attributes"`
}

type AccountRuleAttr struct {
	Resource xdr.AccountRuleResource `json:"resource"`
	Action   string                  `json:"action"`
	IsForbid bool                    `json:"is_forbid"`
	Details  Details                 `json:"details"`
}

type SignerRule struct {
	Key
	Attributes    SignerRuleAttr     `json:"attributes"`
	Relationships SignerRuleRelation `json:"relationships"`
}

type SignerRuleRelation struct {
	Owner *Relation `json:"owner"`
}

type SignerRuleAttr struct {
	Resource  xdr.SignerRuleResource `json:"resource"`
	Action    string                 `json:"action"`
	IsForbid  bool                   `json:"is_forbid"`
	IsDefault bool                   `json:"is_default"`
	Details   Details                `json:"details"`
}
