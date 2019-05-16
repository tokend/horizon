package requests

import (
	"net/http"
)

const (
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

// GetAtomicSwapBid - represents params to be specified by user for Get AtomicSwapBid handler
type GetAtomicSwapBid struct {
	*base
	ID uint64
}

// NewGetAtomicSwapBid returns new instance of GetAtomicSwapBid request
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
