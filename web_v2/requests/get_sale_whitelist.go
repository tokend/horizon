package requests

import (
	"net/http"

	"gitlab.com/tokend/horizon/db2"
)

type GetSaleWhitelist struct {
	*base
	SaleID     uint64
	PageParams *db2.CursorPageParams
}

func NewGetSaleWhitelist(r *http.Request) (*GetSaleWhitelist, error) {
	b, err := newBase(r, baseOpts{})
	if err != nil {
		return nil, err
	}

	id, err := b.getUint64ID()
	if err != nil {
		return nil, err
	}
	pageParams, err := b.getCursorBasedPageParams()
	if err != nil {
		return nil, err
	}

	return &GetSaleWhitelist{
		base:       b,
		SaleID:     id,
		PageParams: pageParams,
	}, nil
}
