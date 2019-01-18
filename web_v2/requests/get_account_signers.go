package requests

import (
	"net/http"
)

//GetAccountSigners - represents params to be specified by user for Get Account Signers handler
type GetAccountSigners struct {
	*base
	Address string
}

//NewGetAccountSigners - returns new instance of GetAccountSigners request
func NewGetAccountSigners(r *http.Request) (*GetAccountSigners, error) {
	b, err := newBase(r, baseOpts{
		supportedIncludes: map[string]struct{}{
			"roles":       {},
			"roles.rules": {},
		},
	})
	if err != nil {
		return nil, err
	}
	address, err := newAccountAddress(b, "id")
	if err != nil {
		return nil, err
	}

	return &GetAccountSigners{
		base:    b,
		Address: address,
	}, nil
}

//NeedRoles - returns true if user requested to include roles or any of roles relationships
func (a *GetAccountSigners) NeedRoles() bool {
	return a.shouldInclude("roles") || a.NeedRules()
}

//NeedRules - returns true if user request to include rules for roles
func (a *GetAccountSigners) NeedRules() bool {
	return a.shouldInclude("roles.rules")
}
