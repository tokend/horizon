package requests

import "net/http"

type GetSaleWhitelist struct {
	*base
	SaleID uint64
}

func NewGetSaleWhiteList(r *http.Request) (*GetSaleWhitelist, error) {
	b, err := newBase(r, baseOpts{})
	if err != nil {
		return nil, err
	}

	id, err := b.getUint64ID()
	if err != nil {
		return nil, err
	}

	return &GetSaleWhitelist{
		base:   b,
		SaleID: id,
	}, nil
}
