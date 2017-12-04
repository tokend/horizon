package resource

import (
	"gitlab.com/swarmfund/go/amount"
	"gitlab.com/swarmfund/horizon/db2/core"
)

type Asset struct {
	Code                 string `json:"code"`
	Owner                string `json:"owner"`
	AvailableForIssuance string `json:"available_for_issuance"`
	Name                 string `json:"name"`
	PreissuedAssetSigner string `json:"preissued_asset_signer"`
	Description          string `json:"description"`
	ExternalResourceLink string `json:"external_resource_link"`
	MaxIssuanceAmount    string `json:"max_issuance_amount"`
	Issued               string `json:"issued"`
	Policies
}

func (a *Asset) Populate(asset *core.Asset) {
	a.Code = asset.Code
	a.Owner = asset.Owner
	a.AvailableForIssuance = amount.StringU(asset.AvailableForIssuance)
	a.Name = asset.Name
	a.PreissuedAssetSigner = asset.PreissuedAssetSigner
	a.Description = asset.Description
	a.ExternalResourceLink = asset.ExternalResourceLink
	a.MaxIssuanceAmount = amount.StringU(asset.MaxIssuanceAmount)
	a.Issued = amount.StringU(asset.Issued)
	a.Policies.Populate(*asset)
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
