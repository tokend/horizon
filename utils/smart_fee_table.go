package utils

import (
	"gitlab.com/swarmfund/horizon/db2/core"
	"gitlab.com/tokend/go/xdr"
)

//SmartFeeTable is used to built complete fee overview,
//using different level fees

type (
	feeType       = int
	feeSubtype    = int64
	assetCode     = string
	SmartFeeTable map[feeType]map[feeSubtype]map[assetCode][]core.FeeEntry
)

func NewSmartFeeTable(fees []core.FeeEntry) (sft SmartFeeTable) {
	sft = SmartFeeTable{}
	for _, entry := range fees {
		if sft[entry.FeeType] == nil {
			sft[entry.FeeType] = make(map[feeSubtype]map[assetCode][]core.FeeEntry)
		}
		if sft[entry.FeeType][entry.Subtype] == nil {
			sft[entry.FeeType][entry.Subtype] = make(map[assetCode][]core.FeeEntry)
		}
		sft[entry.FeeType][entry.Subtype][entry.Asset] = append(sft[entry.FeeType][entry.Subtype][entry.Asset], entry)
	}
	return sft
}

func (sft SmartFeeTable) Update(fees []core.FeeEntry) {
	update := NewSmartFeeTable(fees)

	for feeType, byFeeType := range update {
		if sft[feeType] == nil {
			sft[feeType] = make(map[feeSubtype]map[assetCode][]core.FeeEntry)
		}
		for feeSubtype, byFeeSubType := range byFeeType {
			if sft[feeType][feeSubtype] == nil {
				sft[feeType][feeSubtype] = make(map[assetCode][]core.FeeEntry)
			}
			for asset, newFees := range byFeeSubType {
				sft[feeType][feeSubtype][asset] = SmartFillFeeGaps(sft[feeType][feeSubtype][asset], newFees)
			}
		}
	}
}

func (sft SmartFeeTable) GetValuesByAsset() (byAsset map[string][]core.FeeEntry) {
	byAsset = make(map[assetCode][]core.FeeEntry)
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
			sft[feeType] = make(map[feeSubtype]map[assetCode][]core.FeeEntry)
		}
		subtypes := []int64{0}
		if ft == xdr.FeeTypePaymentFee {
			subtypes = []int64{int64(xdr.PaymentFeeTypeIncoming), int64(xdr.PaymentFeeTypeOutgoing)}
		}
		for _, feeSubtype := range subtypes {
			if sft[feeType][feeSubtype] == nil {
				sft[feeType][feeSubtype] = make(map[assetCode][]core.FeeEntry)
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
