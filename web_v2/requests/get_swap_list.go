package requests

import (
	"gitlab.com/tokend/horizon/db2"
	"net/http"
)

const (
	FilterTypeSwapListSource             = "source"
	FilterTypeSwapListDestination        = "destination"
	FilterTypeSwapListSourceBalance      = "source_balance"
	FilterTypeSwapListDestinationBalance = "destination_balance"
	FilterTypeSwapListAsset              = "asset"
	FilterTypeSwapListState              = "state"
)

var filterTypeSwapListAll = map[string]struct{}{
	FilterTypeSwapListSource:             {},
	FilterTypeSwapListDestination:        {},
	FilterTypeSwapListSourceBalance:      {},
	FilterTypeSwapListDestinationBalance: {},
	FilterTypeSwapListAsset:              {},
	FilterTypeSwapListState:              {},
}

type GetSwapList struct {
	*base
	Filters struct {
		Source             string `json:"source"`
		SourceBalance      string `json:"source_balance"`
		Destination        string `json:"destination"`
		DestinationBalance string `json:"destination_balance"`
		Asset              string `json:"asset"`

		State int32 `json:"state"`
	}
	PageParams *db2.CursorPageParams
}

func NewGetSwapList(r *http.Request) (*GetSwapList, error) {
	b, err := newBase(r, baseOpts{
		supportedFilters:  filterTypeSwapListAll,
		supportedIncludes: includeTypeSwapAll,
	})
	if err != nil {
		return nil, err
	}

	pageParams, err := b.getCursorBasedPageParams()
	if err != nil {
		return nil, err
	}

	request := GetSwapList{
		base:       b,
		PageParams: pageParams,
	}

	err = b.populateFilters(&request.Filters)
	if err != nil {
		return nil, err
	}

	return &request, nil
}
