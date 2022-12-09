package requests

import "net/http"

const (
	IncludeTypeLiquidityPoolAssets = "assets"
)

var includeTypeLiquidityPoolAll = map[string]struct{}{
	IncludeTypeLiquidityPoolAssets: {},
}

type GetLiquidityPool struct {
	*base
	ID       uint64
	Includes struct {
		Assets bool `include:"assets"`
	}
}

func NewGetLiquidityPool(r *http.Request) (*GetLiquidityPool, error) {
	b, err := newBase(r, baseOpts{
		supportedIncludes: includeTypeLiquidityPoolAll,
	})
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
