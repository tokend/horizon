package base

import (
	"bullioncoin.githost.io/development/go/amount"
	"bullioncoin.githost.io/development/go/xdr"
)

type AssetForm struct {
	Name     string `json:"name"`
	UnitSize string `json:"unit_size"`
}

func (a *AssetForm) PopulateFromXdr(assetForm xdr.AssetForm) {
	a.Name = string(assetForm.Name)
	a.UnitSize = amount.String(int64(assetForm.Unit))
}
