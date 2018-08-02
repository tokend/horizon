package utils

import (
	"gitlab.com/swarmfund/horizon/db2/core"
	"gitlab.com/tokend/go/xdr"
)

//SmartFeeTable is used to built complete fee overview,
//using different level fees

type (
	feeType       = xdr.FeeType
	feeSubtype    = int64
	assetCode     = xdr.AssetCode
	SmartFeeTable map[feeType]map[feeSubtype]map[assetCode][]core.FeeEntry
)

func NewSmartFeeTable(fees []core.FeeEntry) (sft SmartFeeTable) {
	sft = SmartFeeTable{}
	for _, entry := range fees {
		ft := feeType(entry.FeeType)
		asset := assetCode(entry.Asset)
		if _, ok := sft[ft]; !ok {
			sft[ft] = make(map[feeSubtype]map[assetCode][]core.FeeEntry)
		}
		if _, ok := sft[ft][entry.Subtype]; !ok {
			sft[ft][entry.Subtype] = make(map[assetCode][]core.FeeEntry)
		}
		sft[ft][entry.Subtype][asset] =
			append(sft[ft][entry.Subtype][asset], entry)
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
				sft[feeType][feeSubtype][asset] =
					SmartFillFeeGaps(sft[feeType][feeSubtype][asset], newFees)
			}
		}
	}
}

func (sft SmartFeeTable) GetValuesByAsset() (byAsset map[string][]core.FeeEntry) {
	byAsset = make(map[string][]core.FeeEntry)
	for feeType := range sft {
		for feeSubtype := range sft[feeType] {
			for asset := range sft[feeType][feeSubtype] {
				byAsset[string(asset)] = sft[feeType][feeSubtype][asset]
			}
		}
	}
	return byAsset
}

func (sft SmartFeeTable) AddZeroFees(assets []core.Asset) {
	for _, ft := range xdr.FeeTypeAll {
		if sft[ft] == nil {
			sft[ft] = make(map[feeSubtype]map[assetCode][]core.FeeEntry)
		}
		subtypes := []feeSubtype{0}
		if ft == xdr.FeeTypePaymentFee {
			subtypes = []feeSubtype{
				feeSubtype(xdr.PaymentFeeTypeIncoming),
				feeSubtype(xdr.PaymentFeeTypeOutgoing),
			}
		}
		for _, feeSubtype := range subtypes {
			if sft[ft][feeSubtype] == nil {
				sft[ft][feeSubtype] = make(map[assetCode][]core.FeeEntry)
			}
			for _, asset := range assets {
				ac := assetCode(asset.Code)
				sft[ft][feeSubtype][ac] = FillFeeGaps(
					sft[ft][feeSubtype][ac],
					core.FeeEntry{
						Asset:   asset.Code,
						Subtype: feeSubtype,
						FeeType: int(ft),
					})
			}
		}
	}
}
