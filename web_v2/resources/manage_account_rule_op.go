package resources

import (
	"gitlab.com/tokend/go/xdr"
	"gitlab.com/tokend/horizon/db2/history2"
	"gitlab.com/tokend/regources/v2"
)

func newManageAccountRule(id int64, details history2.ManageAccountRuleDetails,
) *regources.ManageAccountRule {
	result := &regources.ManageAccountRule{
		Key: regources.NewKeyInt64(id, regources.TypeManageAccountRule),
		Attributes: regources.ManageAccountRuleAttrs{
			Action: details.Action,
			RuleID: details.RuleID,
		},
	}

	switch details.Action {
	case xdr.ManageAccountRuleActionCreate:
		createDetails := regources.UpdateAccountRuleAttrs(*details.CreateDetails)
		result.Attributes.CreateAttrs = &createDetails
	case xdr.ManageAccountRuleActionUpdate:
		updateDetails := regources.UpdateAccountRuleAttrs(*details.UpdateDetails)
		result.Attributes.UpdateAttrs = &updateDetails
	}

	return result
}
