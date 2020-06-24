package requests

import (
	"gitlab.com/distributed_lab/kit/pgdb"
	"gitlab.com/distributed_lab/urlval"
	"net/http"
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

	var request = GetAccountList{base: b}
	err = urlval.Decode(r.URL.Query(), &request) //default sorting is "desc"
	if err != nil {
		return nil, err
	}

	return &request, nil
}
