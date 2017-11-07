package base

import (
	"bullioncoin.githost.io/development/go/amount"
)

type ForfeitItem struct {
	AssetForm
	UnitsNum int64 `json:"units_num"`
}

func (assetForm *ForfeitItem) Populate(name string, unitSize, unitsNum int64) {
	assetForm.Name = name
	assetForm.UnitSize = amount.String(unitSize)
	assetForm.UnitsNum = unitsNum
}

type ForfeitRequest struct {
	TotalPercentFee string        `json:"total_percent_fee"`
	TotalFixedFee   string        `json:"total_fixed_fee"`
	Forms           []ForfeitItem `json:"forms"`
}
