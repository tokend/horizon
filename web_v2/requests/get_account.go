package requests

import (
	"net/http"
)

const (
	IncludeTypeAccountBalances        = "balances"
	IncludeTypeAccountBalancesAsset   = "balances.asset"
	IncludeTypeAccountBalancesState   = "balances.state"
	IncludeTypeAccountAccountReferrer = "referrer"
	IncludeTypeAccountRole            = "role"
	IncludeTypeAccountRoleRules       = "role.rules"
)

var IncludeTypeAccountAll = map[string]struct{}{
	IncludeTypeAccountBalances:        {},
	IncludeTypeAccountBalancesAsset:   {},
	IncludeTypeAccountBalancesState:   {},
	IncludeTypeAccountAccountReferrer: {},
	IncludeTypeAccountRole:            {},
	IncludeTypeAccountRoleRules:       {},
}

//GetAccount - represents params to be specified by user for Get Account handler
type GetAccount struct {
	*base
	Address string
}

//NewGetAccount - returns new instance of GetAccount request
func NewGetAccount(r *http.Request) (*GetAccount, error) {
	b, err := newBase(r, baseOpts{
		supportedIncludes: IncludeTypeAccountAll,
	})
	if err != nil {
		return nil, err
	}
	address, err := newAccountAddress(b, "id")
	if err != nil {
		return nil, err
	}

	return &GetAccount{
		base:    b,
		Address: address,
	}, nil
}
