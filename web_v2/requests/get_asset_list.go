package requests

import (
	"net/http"
)

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
		supportedIncludes: map[string]struct{}{
			"owner": {},
		},
		supportedFilters: map[string]struct{}{
			"policy": {},
			"owner":  {},
		},
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

func (r *GetAssetList) NeedOwner() bool {
	return r.shouldInclude("owner")
}
