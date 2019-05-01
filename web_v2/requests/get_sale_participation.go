package requests

import (
	"net/http"

	"gitlab.com/tokend/horizon/db2"
)

type GetSaleParticipation struct {
	*base
	SaleID     uint64
	PageParams *db2.CursorPageParams
}

func NewGetSaleParticipation(r *http.Request) (*GetSaleParticipation, error) {
	b, err := newBase(r, baseOpts{})
	if err != nil {
		return nil, err
	}

	id, err := b.getUint64ID()
	if err != nil {
		return nil, err
	}

	return &GetSaleParticipation{
		base:   b,
		SaleID: id,
	}, nil
}
