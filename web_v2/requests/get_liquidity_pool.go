package requests

import "net/http"

type GetLiquidityPool struct {
	*base
	ID uint64
}

func NewGetLiquidityPool(r *http.Request) (*GetLiquidityPool, error) {
	b, err := newBase(r, baseOpts{})
	if err != nil {
		return nil, err
	}

	id, err := b.getUint64ID()
	if err != nil {
		return nil, err
	}

	return &GetLiquidityPool{
		base: b,
		ID:   id,
	}, nil
}
