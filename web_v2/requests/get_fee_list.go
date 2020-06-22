package requests

import (
	"gitlab.com/distributed_lab/kit/pgdb"
	"gitlab.com/distributed_lab/urlval"
	"net/http"

	regources "gitlab.com/tokend/regources/generated"
)

const (
	IncludeTypeFeeListAsset = "asset"

	FilterTypeFeeListAsset       = "asset"
	FilterTypeFeeListSubtype     = "subtype"
	FilterTypeFeeListAccount     = "account"
	FilterTypeFeeListAccountRole = "account_role"
	FilterTypeFeeListFeeType     = "fee_type"
	FilterTypeFeeListLowerBound  = "lower_bound"
	FilterTypeFeeListUpperBound  = "upper_bound"
)

var includeTypeFeeListAll = map[string]struct{}{
	IncludeTypeFeeListAsset: {},
}

var filterTypeFeeListAll = map[string]struct{}{
	FilterTypeFeeListAsset:       {},
	FilterTypeFeeListSubtype:     {},
	FilterTypeFeeListAccount:     {},
	FilterTypeFeeListAccountRole: {},
	FilterTypeFeeListLowerBound:  {},
	FilterTypeFeeListUpperBound:  {},
	FilterTypeFeeListFeeType:     {},
}

//GetFeeList - represents params to be specified for Get Fees handler
type GetFeeList struct {
	*base
	Filters GetFeeListFilters
	PageParams *pgdb.OffsetPageParams
}

type GetFeeListFilters struct {
	Asset       []string           `filter:"asset"`
	Subtype     []int64            `filter:"subtype"`
	FeeType     []int32            `filter:"fee_type"`
	Account     []string           `filter:"account"`
	AccountRole []uint64           `filter:"account_role"`
	LowerBound  regources.Amount `filter:"lower_bound"`
	UpperBound  regources.Amount `filter:"upper_bound"`
	}

// NewGetFeeList returns the new instance of GetFeeList request
func NewGetFeeList(r *http.Request) (*GetFeeList, error) {
	b, err := newBase(r, baseOpts{
		supportedIncludes: includeTypeFeeListAll,
		supportedFilters:  filterTypeFeeListAll,
	})
	if err != nil {
		return nil, err
	}

	var pageParams pgdb.OffsetPageParams
	err=urlval.Decode(r.URL.Query(), &pageParams)

	request := GetFeeList{
		base:       b,
		PageParams: &pageParams,
	}


	request.Filters = GetFeeListFilters {
		[]string{""},
		[]int64{0},
		[]int32{0},
		[]string{""},
		[]uint64{0},
		0,
		0,
	}
	err=urlval.Decode(r.URL.Query(), &request.Filters)

	return &request, nil
}
