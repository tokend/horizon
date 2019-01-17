package requests

import "net/http"

//GetAccountSigners - represents params to be specified for Get Assets handler
type GetAssetList struct {
	*base
}

func NewGetAssetList(r *http.Request) (*GetAssetList, error) {
	b, err := newBase(r, map[string]struct{}{
		"owner": {},
	})
	if err != nil {
		return nil, err
	}

	return &GetAssetList{
		base: b,
	}, nil
}

func (r *GetAssetList) NeedOwner() bool {
	return r.shouldInclude("owner")
}
