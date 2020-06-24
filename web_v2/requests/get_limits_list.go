package requests

import (
	"gitlab.com/distributed_lab/kit/pgdb"
	"gitlab.com/distributed_lab/urlval"
	"net/http"
)

const (
	IncludeTypeLimitsListAsset = "asset"

	FilterTypeLimitsListAsset       = "asset"
	FilterTypeLimitsListAccount     = "account"
	FilterTypeLimitsListAccountRole = "account_role"
	FilterTypeLimitsListStatsOpType = "stats_op_type"
)

var includeTypeLimitsListAll = map[string]struct{}{
	IncludeTypeLimitsListAsset: {},
}

var filterTypeLimitsListAll = map[string]struct{}{
	FilterTypeLimitsListAsset:       {},
	FilterTypeLimitsListAccount:     {},
	FilterTypeLimitsListAccountRole: {},
	FilterTypeLimitsListStatsOpType: {},
}

//GetLimitsList - represents params to be specified for Get Fees handler
type GetLimitsList struct {
	*base
	Filters    GetLimitsListFilters
	PageParams *pgdb.OffsetPageParams
}
type GetLimitsListFilters struct {
	Asset       *string `filter:"asset"`
	StatsOpType *int32  `filter:"stats_op_type"`
	Account     *string `filter:"account"`
	AccountRole *uint64 `filter:"account_role"`
}

// NewGetLimitsList returns the new instance of GetLimitsList request
func NewGetLimitsList(r *http.Request) (*GetLimitsList, error) {
	b, err := newBase(r, baseOpts{
		supportedIncludes: includeTypeLimitsListAll,
		supportedFilters:  filterTypeLimitsListAll,
	})
	if err != nil {
		return nil, err
	}

	var pageParams pgdb.OffsetPageParams
	err = urlval.Decode(r.URL.Query(), &pageParams)

	request := GetLimitsList{
		base:       b,
		PageParams: &pageParams,
	}

	err = urlval.Decode(r.URL.Query(), &request.Filters)

	return &request, nil
}
