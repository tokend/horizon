package requests

import (
	"fmt"
	"net/http"

	"gitlab.com/tokend/horizon/db2"
)

// GetSaleList - represents params to be specified by user for getSaleList handler
type GetSaleList struct {
	SalesBase
	PageParams *db2.OffsetPageParams
}

// NewGetSaleList returns new instance of GetSaleList request
func NewGetSaleList(r *http.Request) (*GetSaleList, error) {
	b, err := newBase(r, baseOpts{
		supportedIncludes: includeTypeSaleListAll,
		supportedFilters:  filterTypeSaleListAll,
	})
	if err != nil {
		return nil, err
	}

	pageParams, err := b.getOffsetBasedPageParams()
	if err != nil {
		return nil, err
	}

	request := GetSaleList{
		SalesBase: SalesBase{
			base: b,
		},
		PageParams: pageParams,
	}

	err = b.populateFilters(&request.Filters)
	if err != nil {
		return nil, err
	}

	ids, err := b.getUint64Slice(fmt.Sprintf("filter[%s]", FilterTypeSaleListIDs))
	if err != nil {
		return nil, err
	}
	request.Filters.IDs = ids

	return &request, nil
}
