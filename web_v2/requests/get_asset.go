package requests

import "net/http"

// GetAsset - represents params to be specified by user for Get Asset handler
type GetAsset struct {
	*base
	Code string
}

// NewGetAsset returns new instance of GetAsset request
func NewGetAsset(r *http.Request) (*GetAsset, error) {
	b, err := newBase(r, baseOpts{
		supportedIncludes: map[string]struct{}{
			"owner": {},
		},
	})
	if err != nil {
		return nil, err
	}

	code := b.getString("code")

	return &GetAsset{
		base: b,
		Code: code,
	}, nil
}

func (r *GetAsset) NeedOwner() bool {
	return r.shouldInclude("owner")
}
