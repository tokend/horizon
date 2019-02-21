package requests

import (
	"net/http"
)

const (
	// IncludeTypeSignerRole - defines if signer roles should be included in the response
	IncludeTypeSignerRole = "role"
	// IncludeTypeSignerRoleRules - defines if rules of signer roles should be included in the response
	IncludeTypeSignerRoleRules = "role.rules"
)

var includeTypeSignerAll = map[string]struct{}{
	IncludeTypeSignerRole:      {},
	IncludeTypeSignerRoleRules: {},
}

//GetAccountSigners - represents params to be specified by user for Get Account Signers handler
type GetAccountSigners struct {
	*base
	Address string
}

//NewGetAccountSigners - returns new instance of GetAccountSigners request
func NewGetAccountSigners(r *http.Request) (*GetAccountSigners, error) {
	b, err := newBase(r, baseOpts{
		supportedIncludes: includeTypeSignerAll,
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
