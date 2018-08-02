package utils

import (
	"gitlab.com/swarmfund/horizon/db2/core"
	"gitlab.com/tokend/go/xdr"
)

//SmartFeeTable is used to built complete fee overview,
//using different level fees
type SmartFeeTable map[int]map[int64]map[string][]core.FeeEntry

func NewSmartFeeTable(fees []core.FeeEntry) (sft SmartFeeTable) {
	sft = SmartFeeTable{}
	for _, entry := range fees {
		if sft[entry.FeeType] == nil {
			sft[entry.FeeType] = make(map[int64]map[string][]core.FeeEntry)
		}
		if sft[entry.FeeType][entry.Subtype] == nil {
			sft[entry.FeeType][entry.Subtype] = make(map[string][]core.FeeEntry)
		}
		sft[entry.FeeType][entry.Subtype][entry.Asset] = append(sft[entry.FeeType][entry.Subtype][entry.Asset], entry)
	}
	return sft
}

func (sft SmartFeeTable) Update(fees []core.FeeEntry) {
	feeMap := make(map[string][]core.FeeEntry)
	for _, fee := range fees {
		feeMap[fee.Asset] = append(feeMap[fee.Asset], fee)
	}

	for feeType := range sft {
		for feeSubtype := range sft[feeType] {
			for asset, fees := range sft[feeType][feeSubtype] {
				sft[feeType][feeSubtype][asset] = SmartFillFeeGaps(fees, feeMap[asset])
			}
		}
	}
}

func (sft SmartFeeTable) GetValuesByAsset() (byAsset map[string][]core.FeeEntry) {
	byAsset = make(map[string][]core.FeeEntry)
	for feeType := range sft {
		for feeSubtype := range sft[feeType] {
			for asset := range sft[feeType][feeSubtype] {
				byAsset[asset] = sft[feeType][feeSubtype][asset]
			}
		}
	}
	return byAsset
}

func (sft SmartFeeTable) AddZeroFees(assets []core.Asset) {
	for _, ft := range xdr.FeeTypeAll {
		feeType := int(ft)
		if sft[feeType] == nil {
			sft[feeType] = make(map[int64]map[string][]core.FeeEntry)
		}
		subtypes := []int64{0}
		if ft == xdr.FeeTypePaymentFee {
			subtypes = []int64{int64(xdr.PaymentFeeTypeIncoming), int64(xdr.PaymentFeeTypeOutgoing)}
		}
		for _, feeSubtype := range subtypes {
			if sft[feeType][feeSubtype] == nil {
				sft[feeType][feeSubtype] = make(map[string][]core.FeeEntry)
			}
			for _, asset := range assets {
				sft[feeType][feeSubtype][asset.Code] = FillFeeGaps(sft[feeType][feeSubtype][asset.Code],
					core.FeeEntry{
						Asset:   asset.Code,
						Subtype: feeSubtype,
						FeeType: feeType,
					})
			}
		}
	}
}
