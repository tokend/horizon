package regources

import "gitlab.com/tokend/go/xdr"

//ManageAccountRule - details of corresponding op
type ManageAccountRule struct {
	Key
	Attributes ManageAccountRuleAttrs `json:"attributes"`
}

// ManageAccountRuleAttrs - details of ManageAccountRuleOp
type ManageAccountRuleAttrs struct {
	Action      xdr.ManageAccountRuleAction `json:"action"`
	RuleID      uint64                      `json:"rule_id"`
	CreateAttrs *UpdateAccountRuleAttrs     `json:"create_attrs,omitempty"`
	UpdateAttrs *UpdateAccountRuleAttrs     `json:"update_attrs,omitempty"`
}

// UpdateAccountRoleAttrs - details of new or updated rule
type UpdateAccountRuleAttrs struct {
	Resource xdr.AccountRuleResource `json:"resource"`
	Action   string                  `json:"action"`
	IsForbid bool                    `json:"is_forbid"`
	Details  Details                 `json:"details"`
}
