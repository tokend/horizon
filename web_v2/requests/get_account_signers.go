package requests

import (
	"gitlab.com/distributed_lab/kit/pgdb"
	"gitlab.com/distributed_lab/urlval"
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
	Address    string
	PageParams *pgdb.OffsetPageParams
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

	var pageParams pgdb.OffsetPageParams
	//pageParams, err := b.getOffsetBasedPageParams()
	err=urlval.Decode(r.URL.Query(),&pageParams)
	if err != nil {
		return nil, err
	}

	return &GetAccountSigners{
		base:       b,
		Address:    address,
		PageParams: &pageParams,
	}, nil
}
