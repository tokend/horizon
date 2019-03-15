package requests

import (
	"net/http"
)

const (
	// IncludeTypeAssetOwner - defines if the asset owner should be included in the response
	IncludeTypeBidBaseBalance = "base_balance"
	IncludeTypeBidOwner       = "owner"
	IncludeTypeBidBaseAsset   = "base_asset"
	IncludeTypeBidQuoteAssets = "quote_assets"
)

var includeTypeBidAll = map[string]struct{}{
	IncludeTypeBidBaseBalance: {},
	IncludeTypeBidOwner:       {},
	IncludeTypeBidBaseAsset:   {},
	IncludeTypeBidQuoteAssets: {},
}

// GetAsset - represents params to be specified by user for Get Asset handler
type GetAtomicSwapBid struct {
	*base
	ID uint64
}

// NewGetAsset returns new instance of GetAsset request
func NewGetAtomicSwapBid(r *http.Request) (*GetAtomicSwapBid, error) {
	b, err := newBase(r, baseOpts{
		supportedIncludes: includeTypeBidAll,
	})
	if err != nil {
		return nil, err
	}

	// bid relations has not asset relation, we use balance relation
	if _, ok := b.include[IncludeTypeBidBaseAsset]; ok {
		b.include[IncludeTypeBidBaseBalance] = struct{}{}
	}

	id, err := b.getUint64ID()
	if err != nil {
		return nil, err
	}

	return &GetAtomicSwapBid{
		base: b,
		ID:   id,
	}, nil
}
