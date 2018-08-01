package utils

import "gitlab.com/swarmfund/horizon/db2/core"
//SmartFeeTable is used to built complete fee overview,
//using different level fees
type SmartFeeTable map[int]map[int64]map[string][]core.FeeEntry

func (sft SmartFeeTable) Populate(fees ...core.FeeEntry){
	for _, entry := range fees{
		if sft[entry.FeeType] == nil{
			sft[entry.FeeType] = make(map[int64]map[string][]core.FeeEntry)
		}
		if sft[entry.FeeType][entry.Subtype]== nil{
			sft[entry.FeeType][entry.Subtype] = make(map[string][]core.FeeEntry)
		}
		sft[entry.FeeType][entry.Subtype][entry.Asset] = append(sft[entry.FeeType][entry.Subtype][entry.Asset], entry)
	}
}

func (sft SmartFeeTable) Update(fees []core.FeeEntry){
	feeMap := make(map[string][]core.FeeEntry)
	for _,fee := range fees{
		feeMap[fee.Asset] = append(feeMap[fee.Asset], fee)
	}

	for feeType := range sft{
		for feeSubtype := range sft[feeType]{
			for asset, fees := range sft[feeType][feeSubtype]{
				sft[feeType][feeSubtype][asset] = SmartFillFeeGaps(fees, feeMap[asset])
			}
		}

	}
}

func (sft SmartFeeTable) GetValuesByAsset() (byAsset map[string][]core.FeeEntry){
	byAsset = make(map[string][]core.FeeEntry)
	for feeType := range sft{
		for feeSubtype := range sft[feeType]{
			for asset := range sft[feeType][feeSubtype]{
				byAsset[asset] = sft[feeType][feeSubtype][asset]
			}
		}
	}
	return byAsset
}
//TODO actually use it
func (sft SmartFeeTable) AddZeroFees(assets []core.Asset){
	for feeType := range sft{
		for feeSubtype := range sft[feeType]{
			for _, asset := range assets{
				sft[feeType][feeSubtype][asset.Code] = FillFeeGaps(sft[feeType][feeSubtype][asset.Code], core.FeeEntry{})
			}
		}
	}
}
