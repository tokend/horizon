package requests

import (
	"gitlab.com/distributed_lab/kit/pgdb"
	"gitlab.com/distributed_lab/urlval"
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
		Source             []string `filter:"source"`
		SourceBalance      []string `filter:"source_balance"`
		Destination        []string `filter:"destination"`
		DestinationBalance []string `filter:"destination_balance"`
		Asset              []string `filter:"asset"`

		State []int32 `filter:"state"`
		//Source             string `json:"source"`
		//SourceBalance      string `json:"source_balance"`
		//Destination        string `json:"destination"`
		//DestinationBalance string `json:"destination_balance"`
		//Asset              string `json:"asset"`
		//
		//State int32 `json:"state"`
	}
	PageParams *pgdb.CursorPageParams
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

	request.Filters.State=[]int32{0}
	err=urlval.Decode(r.URL.Query(), &request.Filters)

	return &request, nil
}
