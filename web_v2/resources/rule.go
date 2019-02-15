package resources

import (
	"encoding/json"
	"strconv"

	"gitlab.com/tokend/horizon/db2/core2"

	"gitlab.com/tokend/regources/v2"
)

// NewAccountRule - returns a new (mocked) instance of rule
func NewAccountRule(rule core2.AccountRule) regources.AccountRule {
	var details regources.Details
	_ = json.Unmarshal([]byte(rule.Details), &details)
	return regources.AccountRule{
		Key: regources.Key{
			ID:   strconv.FormatUint(rule.ID, 10),
			Type: regources.TypeRules,
		},
		Attributes: regources.AccountRuleAttr{
			Resource: "NOTE: format will be changed",
			Action:   rule.Action,
			IsForbid: rule.IsForbid,
			Details:  details,
		},
	}
}

// NewAccountRuleKey - returns a new (mocked) rule key
func NewAccountRuleKey(ruleID uint64) regources.Key {
	return regources.Key{
		ID:   strconv.FormatUint(ruleID, 10),
		Type: regources.TypeRules,
	}
}

// NewAccountRule - returns a new (mocked) instance of rule
func NewSignerRule(rule core2.SignerRule) regources.SignerRule {
	var details regources.Details
	_ = json.Unmarshal([]byte(rule.Details), &details)
	return regources.SignerRule{
		Key: regources.Key{
			ID:   strconv.FormatUint(rule.ID, 10),
			Type: regources.TypeSignerRules,
		},
		Attributes: regources.SignerRuleAttr{
			Resource:  "NOTE: format will be changed",
			Action:    rule.Action,
			IsForbid:  rule.IsForbid,
			IsDefault: rule.IsDefault,
			OwnerID:   rule.OwnerID,
			Details:   details,
		},
	}
}

// NewAccountRuleKey - returns a new (mocked) rule key
func NewSignerRuleKey(ruleID uint64) regources.Key {
	return regources.Key{
		ID:   strconv.FormatUint(ruleID, 10),
		Type: regources.TypeSignerRules,
	}
}
