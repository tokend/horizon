package regources

import "gitlab.com/tokend/go/xdr"

//ManageAccountRole - details of corresponding op
type ManageAccountRole struct {
	Key
	Attributes ManageAccountRoleAttrs `json:"attributes"`
}

// ManageAccountRoleAttrs - details of ManageAccountRuleOp
type ManageAccountRoleAttrs struct {
	Action      xdr.ManageAccountRoleAction `json:"action"`
	RoleID      uint64                      `json:"role_id"`
	CreateAttrs *UpdateAccountRoleAttrs     `json:"create_attrs,omitempty"`
	UpdateAttrs *UpdateAccountRoleAttrs     `json:"update_attrs,omitempty"`
}

// UpdateAccountRoleAttrs - details of new or updated rule
type UpdateAccountRoleAttrs struct {
	RuleIDs []uint64 `json:"rule_ids"`
	Details Details  `json:"details"`
}
