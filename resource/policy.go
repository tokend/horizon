package resource

import (
	"gitlab.com/tokend/go/xdr"
	"gitlab.com/tokend/horizon/db2/core"
)

type Policy struct {
	Name  string `json:"name"`
	Value int32  `json:"value"`
}

// Signer represents one of an account's signers.
type Policies struct {
	PolicyI  int32    `json:"policy"`
	Policies []Policy `json:"policies"`
}

type ExchangePolicies struct {
	Asset string `json:"asset"`
	Policies
}

// Populate fills out the fields of the signer, using one of an account's
// secondary signers.
func (p *Policies) Populate(row core.Asset) {
	p.PolicyI = row.Policies
	for _, policy := range xdr.AssetPolicyAll {
		if (int32(policy) & p.PolicyI) != 0 {
			p.Policies = append(p.Policies, Policy{
				Value: int32(policy),
				Name:  policy.String(),
			})
		}
	}
}

func (p *Policies) PopulateForAssetPair(row core.AssetPair) {
	p.PolicyI = row.Policies
	for _, policy := range xdr.AssetPairPolicyAll {
		if (int32(policy) & p.PolicyI) != 0 {
			p.Policies = append(p.Policies, Policy{
				Value: int32(policy),
				Name:  policy.String(),
			})
		}
	}
}

func (p *Policies) PopulateForExchange(row core.ExchangePolicies) {
	p.PolicyI = row.Policies
	for _, policy := range xdr.ExchangeAssetPolicyAll {
		if (int32(policy) & p.PolicyI) != 0 {
			p.Policies = append(p.Policies, Policy{
				Value: int32(policy),
				Name:  policy.String(),
			})
		}
	}
}
