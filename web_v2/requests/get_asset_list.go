package requests

import (
	"net/http"
)

const (
	IncludeTypeAssetListOwners = "owner"

	FilterTypeAssetListOwner  = "owner"
	FilterTypeAssetListPolicy = "policy"
)

var IncludeTypeAssetListAll = map[string]struct{}{
	IncludeTypeAssetListOwners: {},
}

var FilterTypeAssetListAll = map[string]struct{}{
	FilterTypeAssetListOwner:  {},
	FilterTypeAssetListPolicy: {},
}

//GetAccountSigners - represents params to be specified for Get Assets handler
type GetAssetList struct {
	*base
	filters struct {
		Policy uint64 `fig:"policy"`
		Owner  string `fig:"owner"`
	}
}

func NewGetAssetList(r *http.Request) (*GetAssetList, error) {
	b, err := newBase(r, baseOpts{
		supportedIncludes: IncludeTypeAssetListAll,
		supportedFilters:  FilterTypeAssetListAll,
	})
	if err != nil {
		return nil, err
	}

	request := GetAssetList{
		base: b,
	}

	err = b.populateFilters(&request.filters)
	if err != nil {
		return nil, err
	}

	return &request, nil
}
