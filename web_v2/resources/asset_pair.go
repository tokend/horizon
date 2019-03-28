package resources

import (
	"fmt"
	"gitlab.com/tokend/go/xdr"
	"gitlab.com/tokend/horizon/db2/core2"
	"gitlab.com/tokend/regources/rgenerated"
)

// NewAssetPair - creates new instance of AssetPair from provided one.
func NewAssetPair(record core2.AssetPair) rgenerated.AssetPair {
	return rgenerated.AssetPair{
		Key: rgenerated.Key{
			ID:   fmt.Sprintf("%s:%s", record.Base, record.Quote),
			Type: rgenerated.ASSET_PAIRS,
		},
		Attributes: rgenerated.AssetPairAttributes{
			Price:    rgenerated.Amount(record.CurrentPrice),
			Policies: xdr.AssetPairPolicy(record.Policies),
		},
	}
}
