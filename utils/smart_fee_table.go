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
		st := entry.Subtype
		if _, ok := sft[ft]; !ok {
			sft[ft] = make(map[feeSubtype]map[assetCode][]core.FeeEntry)
		}
		if _, ok := sft[ft][st]; !ok {
			sft[ft][st] = make(map[assetCode][]core.FeeEntry)
		}
		sft[ft][st][asset] =
			append(sft[ft][st][asset], entry)
	}
	return sft
}

func (sft SmartFeeTable) Update(fees []core.FeeEntry) {
	update := NewSmartFeeTable(fees)

	for ft, byFeeType := range update {
		if _, ok := sft[ft]; !ok {
			sft[ft] = make(map[feeSubtype]map[assetCode][]core.FeeEntry)
		}
		for st, byFeeSubType := range byFeeType {
			if _, ok := sft[ft][st]; !ok {
				sft[ft][st] = make(map[assetCode][]core.FeeEntry)
			}
			for asset, newFees := range byFeeSubType {
				sft[ft][st][asset] =
					SmartFillFeeGaps(sft[ft][st][asset], newFees)
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
		if _, ok := sft[ft]; !ok {
			sft[ft] = make(map[feeSubtype]map[assetCode][]core.FeeEntry)
		}
		subtypes := []feeSubtype{0}
		if ft == xdr.FeeTypePaymentFee {
			subtypes = []feeSubtype{
				feeSubtype(xdr.PaymentFeeTypeIncoming),
				feeSubtype(xdr.PaymentFeeTypeOutgoing),
			}
		}
		for _, st := range subtypes {
			if _, ok := sft[ft][st]; !ok {
				sft[ft][st] = make(map[assetCode][]core.FeeEntry)
			}
			for _, asset := range assets {
				ac := assetCode(asset.Code)
				sft[ft][st][ac] = FillFeeGaps(
					sft[ft][st][ac],
					core.FeeEntry{
						Asset:   asset.Code,
						Subtype: st,
						FeeType: int(ft),
					})
			}
		}
	}
}
