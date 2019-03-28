package resources

import (
	"strconv"

	"gitlab.com/tokend/go/xdr"

	"gitlab.com/tokend/horizon/db2/core2"

	"gitlab.com/tokend/regources/rgenerated"
)

// NewAccountRule - returns a new instance of account rule
func NewAccountRule(rule core2.AccountRule) rgenerated.AccountRule {
	return rgenerated.AccountRule{
		Key: NewAccountRuleKey(rule.ID),
		Attributes: rgenerated.AccountRuleAttributes{
			Resource: rule.Resource,
			Action:   xdr.AccountRuleAction(rule.Action),
			Forbids:  rule.Forbids,
			Details:  rule.Details,
		},
	}
}

func NewAccountRuleKey(id uint64) rgenerated.Key {
	return rgenerated.Key{
		ID:   strconv.FormatUint(id, 10),
		Type: rgenerated.ACCOUNT_RULES,
	}
}

// NewSignerRule - returns a new instance of signer rule
func NewSignerRule(rule core2.SignerRule) rgenerated.SignerRule {
	return rgenerated.SignerRule{
		Key: rgenerated.Key{
			ID:   strconv.FormatUint(rule.ID, 10),
			Type: rgenerated.SIGNER_RULES,
		},
		Attributes: rgenerated.SignerRuleAttributes{
			Resource:  rule.Resource,
			Action:    xdr.SignerRuleAction(rule.Action),
			Forbids:   rule.Forbids,
			IsDefault: rule.IsDefault,
			Details:   rule.Details,
		},
		Relationships: rgenerated.SignerRuleRelationships{
			Owner: NewAccountKey(rule.OwnerID).AsRelation(),
		},
	}
}
