package requests

import (
	"net/http"

	"gitlab.com/distributed_lab/kit/pgdb"
)

type GetLiquidityPoolList struct {
	LiquidityPoolsBase
	PageParams *pgdb.OffsetPageParams
}

func NewGetLiquidityPoolList(r *http.Request) (*GetLiquidityPoolList, error) {
	b, err := newBase(r, baseOpts{
		supportedFilters:  filterTypeLiquidityPoolListAll,
		supportedIncludes: includeTypeLiquidityPoolListAll,
	})
	if err != nil {
		return nil, err
	}

	pageParams, err := b.getOffsetBasedPageParams()
	if err != nil {
		return nil, err
	}

	request := GetLiquidityPoolList{
		LiquidityPoolsBase: LiquidityPoolsBase{
			base: b,
		},
		PageParams: pageParams,
	}

	err = b.populateFilters(&request.Filters)
	if err != nil {
		return nil, err
	}

	return &request, nil
}
