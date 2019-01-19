package requests

import (
	"net/http"
)

const (
	IncludeTypeSignerRoles      = "roles"
	IncludeTypeSignerRolesRules = "roles.rules"
)

var IncludeTypeSignerAll = map[string]struct{}{
	IncludeTypeSignerRoles:      {},
	IncludeTypeSignerRolesRules: {},
}

//GetAccountSigners - represents params to be specified by user for Get Account Signers handler
type GetAccountSigners struct {
	*base
	Address string
}

//NewGetAccountSigners - returns new instance of GetAccountSigners request
func NewGetAccountSigners(r *http.Request) (*GetAccountSigners, error) {
	b, err := newBase(r, baseOpts{
		supportedIncludes: IncludeTypeSignerAll,
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
