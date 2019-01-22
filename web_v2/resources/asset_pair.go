package resources

import (
	"fmt"
	"gitlab.com/tokend/go/xdr"
	"gitlab.com/tokend/horizon/db2/core2"
	"gitlab.com/tokend/regources/v2"
)

// NewAssetPair - creates new instance of AssetPair from provided one.
func NewAssetPair(record core2.AssetPair) regources.AssetPair {
	return regources.AssetPair{
		Key: regources.Key{
			ID:   fmt.Sprintf("%s:%s", record.BaseAsset, record.QuoteAsset),
			Type: regources.TypeAssetPairs,
		},
		Attributes: regources.AssetPairAttrs{
			Price:    regources.Amount(record.CurrentPrice),
			Policies: xdr.AssetPairPolicy(record.Policies),
		},
	}
}
