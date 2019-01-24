package regources

import "gitlab.com/tokend/go/xdr"

//ManageAssetPair - details of corresponding op
type ManageAssetPair struct {
	Key
	Attributes ManageAssetPairAttrs `json:"attributes"`
}

//ManageAssetPairAttrs - details of corresponding op
type ManageAssetPairAttrs struct {
	BaseAsset               string              `json:"base_asset"`
	QuoteAsset              string              `json:"quote_asset"`
	PhysicalPrice           Amount              `json:"physical_price"`
	PhysicalPriceCorrection Amount              `json:"physical_price_correction"`
	MaxPriceStep            Amount              `json:"max_price_step"`
	Policies                xdr.AssetPairPolicy `json:"policies"`
}
