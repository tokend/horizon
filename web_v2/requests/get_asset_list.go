package requests

import (
	"net/http"

	"gitlab.com/tokend/horizon/db2"
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
)

var includeTypeAssetListAll = map[string]struct{}{
	IncludeTypeAssetListOwners: {},
}

var filterTypeAssetListAll = map[string]struct{}{
	FilterTypeAssetListOwner:  {},
	FilterTypeAssetListPolicy: {},
	FilterTypeAssetListState:  {},
}

//GetAssetList - represents params to be specified for Get Assets handler
type GetAssetList struct {
	*base
	Filters struct {
		Policy uint64 `fig:"policy"`
		Owner  string `fig:"owner"`
		State  uint32 `fig:"state"`
	}
	PageParams *db2.OffsetPageParams
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

	pageParams, err := b.getOffsetBasedPageParams()
	if err != nil {
		return nil, err
	}

	request := GetAssetList{
		base:       b,
		PageParams: pageParams,
	}

	err = b.populateFilters(&request.Filters)
	if err != nil {
		return nil, err
	}

	return &request, nil
}
