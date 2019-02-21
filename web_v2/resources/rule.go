package resources

import (
	"strconv"

	"gitlab.com/tokend/horizon/db2/core2"

	"gitlab.com/tokend/regources/v2"
)

// NewAccountRule - returns a new instance of account rule
func NewAccountRule(rule core2.AccountRule) regources.AccountRule {
	return regources.AccountRule{
		Key: NewAccountRuleKey(rule.ID),
		Attributes: regources.AccountRuleAttr{
			Resource: rule.Resource,
			Action:   rule.Action,
			Forbids:  rule.Forbids,
			Details:  rule.Details,
		},
	}
}

func NewAccountRuleKey(id uint64) regources.Key {
	return regources.Key{
		ID:   strconv.FormatUint(id, 10),
		Type: regources.TypeAccountRules,
	}
}

// NewSignerRule - returns a new instance of signer rule
func NewSignerRule(rule core2.SignerRule) regources.SignerRule {
	return regources.SignerRule{
		Key: regources.Key{
			ID:   strconv.FormatUint(rule.ID, 10),
			Type: regources.TypeSignerRules,
		},
		Attributes: regources.SignerRuleAttr{
			Resource:  rule.Resource,
			Action:    rule.Action,
			Forbids:   rule.Forbids,
			IsDefault: rule.IsDefault,
			Details:   rule.Details,
		},
		Relationships: regources.SignerRuleRelation{
			Owner: NewAccountKey(rule.OwnerID).AsRelation(),
		},
	}
}
