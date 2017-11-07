package resource

import (
	"bullioncoin.githost.io/development/go/amount"
	"gitlab.com/distributed_lab/tokend/horizon/db2/core"
	"gitlab.com/distributed_lab/tokend/horizon/resource/base"
)

type Asset struct {
	Code       string           `json:"code"`
	Token      string           `json:"token"`
	AssetForms []base.AssetForm `json:"asset_forms"`
	Policies
}

func (a *Asset) Populate(asset *core.Asset) {
	a.Code = asset.Code
	a.Policies.Populate(*asset)
	a.Token = asset.Token
	a.AssetForms = make([]base.AssetForm, len(asset.AssetForms))
	for i := range asset.AssetForms {
		a.AssetForms[i].PopulateFromXdr(asset.AssetForms[i])
	}
}

type AssetPair struct {
	BaseAsset               string `json:"base"`
	QuoteAsset              string `json:"quote"`
	CurrentPrice            string `json:"current_price"`
	PhysicalPrice           string `json:"physical_price"`
	PhysicalPriceCorrection string `json:"physical_price_correction"`
	MaxPriceStep            string `json:"max_price_step"`
	Policies
}

func (a *AssetPair) Populate(asset *core.AssetPair) {
	a.BaseAsset = asset.BaseAsset
	a.QuoteAsset = asset.QuoteAsset
	a.CurrentPrice = amount.String(asset.CurrentPrice)
	a.PhysicalPrice = amount.String(asset.PhysicalPrice)
	a.PhysicalPriceCorrection = amount.String(asset.PhysicalPriceCorrection)
	a.MaxPriceStep = amount.String(asset.MaxPriceStep)
	a.Policies.PopulateForAssetPair(*asset)
}
