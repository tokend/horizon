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

type AccountRuleResponse struct {
	Data     AccountRule `json:"data"`
}

type AccountRulesResponse struct {
	Links    *Links        `json:"links"`
	Data     []AccountRule `json:"data"`
}
