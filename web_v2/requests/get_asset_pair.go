package requests

import (
	"github.com/go-ozzo/ozzo-validation"
	"gitlab.com/distributed_lab/logan/v3/errors"
	"net/http"
	"strings"
)

const (
	// IncludeTypeAssetPairBaseAsset - defines if the base asset should be included in the response
	IncludeTypeAssetPairBaseAsset = "base_asset"
	// IncludeTypeAssetPairQuoteAsset - defines if the quote asset should be included in the response
	IncludeTypeAssetPairQuoteAsset = "quote_asset"
)

var includeTypeAssetPairAll = map[string]struct{}{
	IncludeTypeAssetPairBaseAsset:  {},
	IncludeTypeAssetPairQuoteAsset: {},
}

// GetAssetPair - represents params to be specified by user for getAssetPair handler
type GetAssetPair struct {
	*base

	BaseAsset  string
	QuoteAsset string
}

// NewGetAssetPair returns new instance of GetAssetPair request
func NewGetAssetPair(r *http.Request) (*GetAssetPair, error) {
	b, err := newBase(r, baseOpts{
		supportedIncludes: includeTypeAssetPairAll,
	})
	if err != nil {
		return nil, err
	}

	id := b.getString("id")

	baseAsset, quoteAsset, err := getAssetCodesFromId(id)
	if err != nil {
		return nil, err
	}

	return &GetAssetPair{
		base:       b,
		BaseAsset:  baseAsset,
		QuoteAsset: quoteAsset,
	}, nil
}

// assetCodesFromId receives the compound ID in `base:quote` format and
// returns separately both codes or validation error if id is invalid
func getAssetCodesFromId(id string) (baseAsset, quoteAsset string, err error) {
	codes := strings.Split(id, ":")
	if len(codes) != 2 || codes[0] == "" || codes[1] == "" {
		err = validation.Errors{
			"id": errors.New("should be in `base:quote` format"),
		}
		return
	}

	if codes[0] == codes[1] {
		err = validation.Errors{
			"id": errors.New("can't contain equal base:quote values"),
		}
		return
	}

	baseAsset = codes[0]
	quoteAsset = codes[1]

	return
}
