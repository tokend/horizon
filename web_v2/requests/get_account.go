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
		"balances":        {},
		"balances.assets": {},
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
	return a.shouldInclude("balances") || a.NeedBalanceWithAsset()
}

//NeedBalanceWithAsset - returns true if user request to include assets for balance
func (a *GetAccount) NeedBalanceWithAsset() bool {
	return a.shouldInclude("balances.assets")
}
