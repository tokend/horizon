package smartfeetable

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
	SmartFeeTable map[assetCode]map[feeType]map[feeSubtype][]core.FeeEntry
)

func NewSmartFeeTable(fees []core.FeeEntry) (sft SmartFeeTable) {
	sft = SmartFeeTable{}
	for _, entry := range fees {
		ft := feeType(entry.FeeType)
		asset := assetCode(entry.Asset)
		st := entry.Subtype
		if _, ok := sft[asset]; !ok {
			sft[asset] = make(map[feeType]map[feeSubtype][]core.FeeEntry)
		}
		if _, ok := sft[asset][ft]; !ok {
			sft[asset][ft] = make(map[feeSubtype][]core.FeeEntry)
		}
		sft[asset][ft][st] =
			append(sft[asset][ft][st], entry)
	}
	return sft
}

func (sft SmartFeeTable) Update(fees []core.FeeEntry) {
	update := NewSmartFeeTable(fees)

	for asset, byAsset := range update {
		if _, ok := sft[asset]; !ok {
			sft[asset] = make(map[feeType]map[feeSubtype][]core.FeeEntry)
		}
		for ft, byFeeType := range byAsset {
			if _, ok := sft[asset][ft]; !ok {
				sft[asset][ft] = make(map[feeSubtype][]core.FeeEntry)
			}
			for st, newFees := range byFeeType {
				sft[asset][ft][st] =
					SmartFillFeeGaps(sft[asset][ft][st], newFees)
			}
		}
	}
}

func (sft SmartFeeTable) GetValuesByAsset() (byAsset map[string][]core.FeeEntry) {
	byAsset = make(map[string][]core.FeeEntry)
	for ac := range sft {
		for ft := range sft[ac] {
			for st := range sft[ac][ft] {
				byAsset[string(ac)] = sft[ac][ft][st]
			}
		}
	}
	return byAsset
}

func (sft SmartFeeTable) AddZeroFees(assets []core.Asset) {
	for _, asset := range assets {
		ac := assetCode(asset.Code)
		if _, ok := sft[ac]; !ok {
			sft[ac] = make(map[feeType]map[feeSubtype][]core.FeeEntry)
		}
		for _, ft := range xdr.FeeTypeAll {
			if _, ok := sft[ac][ft]; !ok {
				sft[ac][ft] = make(map[feeSubtype][]core.FeeEntry)
			}
			subtypes := []feeSubtype{0}
			if ft == xdr.FeeTypePaymentFee {
				subtypes = []feeSubtype{
					feeSubtype(xdr.PaymentFeeTypeIncoming),
					feeSubtype(xdr.PaymentFeeTypeOutgoing),
				}
			}
			for _, st := range subtypes {
				sft[ac][ft][st] = FillFeeGaps(
					sft[ac][ft][st],
					core.FeeEntry{
						Asset:   asset.Code,
						Subtype: st,
						FeeType: int(ft),
					})
			}
		}
	}
}
