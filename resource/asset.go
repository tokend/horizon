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
	LogoID 				 string `json:"logo_id"`
}

func (a *Asset) Populate(asset *core.Asset) {
	a.Code = asset.Code
	a.Owner = asset.Owner
	a.AvailableForIssuance = amount.StringU(asset.AvailableForIssuance)

	a.PreissuedAssetSigner = asset.PreissuedAssetSigner

	a.MaxIssuanceAmount = amount.StringU(asset.MaxIssuanceAmount)
	a.Issued = amount.StringU(asset.Issued)
	a.Policies.Populate(*asset)
	details := asset.GetDetails()
	a.Name = details.Name
	a.LogoID = details.LogoID
	a.Description = details.Description
	a.ExternalResourceLink = details.ExternalResourceLink
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
