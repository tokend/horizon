package requests

import (
	"net/http"

	"gitlab.com/distributed_lab/kit/pgdb"
	"gitlab.com/distributed_lab/urlval"
)

const (
	// FilterTypeAccountListAccount - defines if we need to filter list by account id
	FilterTypeAccountListAccount = "account"
	// FilterTypeAccountListRole - defines if we need to filter list by role
	FilterTypeAccountListRole = "role"
)

var filterTypeAccountListAll = map[string]struct{}{
	FilterTypeAccountListAccount: {},
	FilterTypeAccountListRole:    {},
}

// GetAccountList - represents params to be specified by user for Get Account list handler
type GetAccountList struct {
	*base
	Filters struct {
		Account []string `filter:"account"`
		Role    []uint64 `filter:"role"`
	}
	PageParams pgdb.OffsetPageParams
}

// NewGetAccountList - returns new instance of GetAccountList request
func NewGetAccountList(r *http.Request) (*GetAccountList, error) {
	b, err := newBase(r, baseOpts{
		supportedFilters: filterTypeAccountListAll,
	})
	if err != nil {
		return nil, err
	}

	pageParams, err := b.getOffsetBasedPageParams()
	var request = GetAccountList{
		base:       b,
		PageParams: *pageParams,
	}
	err = urlval.Decode(r.URL.Query(), &request.Filters)
	if err != nil {
		return nil, err
	}

	return &request, nil
}
