package resources

import "gitlab.com/tokend/regources/v2"

// NewRule - returns a new (mocked) instance of rule
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

// NewRuleKey - returns a new (mocked) rule key
func NewRuleKey() regources.Key {
	return regources.Key{
		ID:   "mocked_rule_id",
		Type: regources.TypeRules,
	}
}
