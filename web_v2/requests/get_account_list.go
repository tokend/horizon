package requests

import (
	"net/http"

	"gitlab.com/tokend/horizon/db2"
)

const (
	//
	FilterTypeAccountListAccount = "account"
	FilterTypeAccountListRole    = "role"
)

var filterTypeAccountListAll = map[string]struct{}{
	FilterTypeAccountListAccount: {},
	FilterTypeAccountListRole:    {},
}

// GetAccountList - represents params to be specified by user for Get Account list handler
type GetAccountList struct {
	*base
	Filters struct {
		Account []string `fig:"account"`
		Role    []uint64 `fig:"role"`
	}
	PageParams db2.OffsetPageParams
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
	if err != nil {
		return nil, err
	}

	request := GetAccountList{
		base:       b,
		PageParams: *pageParams,
	}

	err = b.populateFilters(&request.Filters)
	if err != nil {
		return nil, err
	}

	return &request, nil
}
