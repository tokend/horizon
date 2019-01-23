package requests

import "net/http"

const (
	// IncludeTypeOfferBaseAsset - defines if the base asset should be included in the response
	IncludeTypeOfferBaseAsset = "base_asset"
	// IncludeTypeOfferQuoteAsset - defines if the quote asset should be included in the response
	IncludeTypeOfferQuoteAsset = "quote_asset"
)

var includeTypeOfferAll = map[string]struct{}{
	IncludeTypeOfferBaseAsset:  {},
	IncludeTypeOfferQuoteAsset: {},
}

// GetOffer represents params to be specified by user for getOffer handler
type GetOffer struct {
	*base
	ID string
}

// NewGetOffer returns new instance of the GetOffer request
func NewGetOffer(r *http.Request) (*GetOffer, error) {
	b, err := newBase(r, baseOpts{
		supportedIncludes: includeTypeOfferAll,
	})
	if err != nil {
		return nil, err
	}

	id := b.getString("id")

	return &GetOffer{
		base: b,
		ID:   id,
	}, nil
}
