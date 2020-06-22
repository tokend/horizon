package requests

import (
	"gitlab.com/distributed_lab/kit/pgdb"
	"gitlab.com/distributed_lab/urlval"
	"net/http"
)

const (
	// IncludeTypeAssetListOwners - defines if the asset owners should be included in the response
	IncludeTypeAssetListOwners = "owner"

	// FilterTypeAssetListOwner - defines if we need to filter the list by owner
	FilterTypeAssetListOwner = "owner"
	// FilterTypeAssetListPolicy - defines if we need to filter the list by policy
	FilterTypeAssetListPolicy = "policy"
	//FilterTypeAssetListState - defines if we need to filter the list by asset state
	FilterTypeAssetListState = "state"
	//FilterTypeAssetListCodes - defines if we need to filter the list by asset codes
	FilterTypeAssetListCodes = "codes"
)

var includeTypeAssetListAll = map[string]struct{}{
	IncludeTypeAssetListOwners: {},
}

var filterTypeAssetListAll = map[string]struct{}{
	FilterTypeAssetListOwner:  {},
	FilterTypeAssetListPolicy: {},
	FilterTypeAssetListState:  {},
	FilterTypeAssetListCodes:  {},
}

//GetAssetList - represents params to be specified for Get Assets handler
type GetAssetList struct {
	*base
	Filters struct {
		Policy []uint64   `filter:"policy"`
		Owner  []string   `filter:"owner"`
		State  []uint32   `filter:"state"`
		Codes  []string `filter:"codes"`
	}
	PageParams *pgdb.OffsetPageParams
}

// NewGetAssetList returns the new instance of GetAssetList request
func NewGetAssetList(r *http.Request) (*GetAssetList, error) {
	b, err := newBase(r, baseOpts{
		supportedIncludes: includeTypeAssetListAll,
		supportedFilters:  filterTypeAssetListAll,
	})
	if err != nil {
		return nil, err
	}


	var pageParams pgdb.OffsetPageParams
	request := GetAssetList{
		base:       b,
		PageParams: &pageParams,
		Filters: struct {
			Policy []uint64   `filter:"policy"`
			Owner  []string   `filter:"owner"`
			State  []uint32   `filter:"state"`
			Codes  []string `filter:"codes"`
		}{
			Policy: []uint64{0},
			Owner: []string{""},
			State: []uint32{0},
			Codes: []string{""},
			},
	}
	err=urlval.Decode(r.URL.Query(),&request.Filters)
	err=urlval.Decode(r.URL.Query(),request.PageParams)


	return &request, nil
}
