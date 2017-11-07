package resource

import (
	"bullioncoin.githost.io/development/go/xdr"
)

type AccountPolicyType struct {
	Name  string `json:"name"`
	Value int32  `json:"value"`
}

type AccountPolicies struct {
	AccountPoliciesTypeI int32               `json:"account_policies_type_i"`
	AccountPoliciesTypes []AccountPolicyType `json:"account_policies_types"`
}

func (ap *AccountPolicies) Populate(accountPoliciesType int32) {
	ap.AccountPoliciesTypeI = accountPoliciesType
	for _, accountPolicyType := range xdr.AccountPoliciesAll {
		if (int32(accountPolicyType) & ap.AccountPoliciesTypeI) != 0 {
			ap.AccountPoliciesTypes = append(ap.AccountPoliciesTypes, AccountPolicyType{
				Value: int32(accountPolicyType),
				Name:  accountPolicyType.String(),
			})
		}
	}
}
