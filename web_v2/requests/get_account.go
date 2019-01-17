package requests

import (
	"net/http"
)

//GetAccount - represents params to be specified by user for Get Account handler
type GetAccount struct {
	*base
	Address string
}

//NewGetAccount - returns new instance of GetAccount request
func NewGetAccount(r *http.Request) (*GetAccount, error) {
	b, err := newBase(r, map[string]struct{}{
		"balances":       {},
		"balances.asset": {},
		"balances.state": {},
		"referrer":       {},
		"state":          {},
		"role":           {},
		"role.rules":     {},
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

//NeedBalance - returns true if user requested to include balances or any of balance relationships
func (a *GetAccount) NeedBalance() bool {
	return a.shouldInclude("balances") || a.NeedBalanceWithAsset() || a.NeedBalanceState()
}

//NeedBalanceWithAsset - returns true if user request to include assets for balance
func (a *GetAccount) NeedBalanceWithAsset() bool {
	return a.shouldInclude("balances.asset")
}

//NeedBalanceState - returns true if user requested to include balance state for balance
func (a *GetAccount) NeedBalanceState() bool {
	return a.shouldInclude("balances.state")
}

//NeedReferrer - returns true if user requested to include referrer
func (a *GetAccount) NeedReferrer() bool {
	return a.shouldInclude("referrer")
}

//NeedAccountState - returns true if user requested to include account state
func (a *GetAccount) NeedAccountState() bool {
	return a.shouldInclude("state")
}

//NeedRole - returns true if user requested to include role
func (a *GetAccount) NeedRole() bool {
	return a.shouldInclude("role") || a.NeedRules()
}

//NeedRules - returns true if user requested to include rule
func (a *GetAccount) NeedRules() bool {
	return a.shouldInclude("role.rules")
}
