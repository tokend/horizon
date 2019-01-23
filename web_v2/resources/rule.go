package resources

import "gitlab.com/tokend/regources/v2"

func NewRule() regources.Rule {
	return regources.Rule{
		Key: regources.Key{
			ID:   "mocked_rule_id",
			Type: regources.TypeRules,
		},
		Attributes: regources.RuleAttr{
			Resource: "NOTE: format will be changed",
			Action:   "view",
			Details: map[string]interface{}{
				"name": "Name of the mocked Rule",
			},
		},
	}
}

func NewRuleKey() regources.Key {
	return regources.Key{
		ID:   "mocked_rule_id",
		Type: regources.TypeRules,
	}
}
