package requests

import (
	"net/http"

	"gitlab.com/distributed_lab/logan/v3/errors"

	"gitlab.com/tokend/horizon/db2"
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
	Filters struct {
		Asset       string `fig:"asset"`
		StatsOpType int32  `fig:"stats_op_type"`
		Account     string `fig:"account"`
		AccountRole uint64 `fig:"account_role"`
	}
	PageParams *db2.OffsetPageParams
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

	pageParams, err := b.getOffsetBasedPageParams()
	if err != nil {
		return nil, errors.Wrap(err, "failed to get page params`")
	}

	request := GetLimitsList{
		base:       b,
		PageParams: pageParams,
	}

	err = b.populateFilters(&request.Filters)
	if err != nil {
		return nil, errors.Wrap(err, "failed to populate filters")
	}

	return &request, nil
}
