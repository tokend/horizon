package smartfeetable

import (
	"gitlab.com/swarmfund/horizon/db2/core"
	"gitlab.com/tokend/go/xdr"
)

//SmartFeeTable is used to built complete fee overview,
//using different level fees

type (
	FeeGroup struct {
		AssetCode string
		FeeType   int
		Subtype   int64
	}
	SmartFeeTable map[FeeGroup][]core.FeeEntry
)

func NewSmartFeeTable(fees []core.FeeEntry) (sft SmartFeeTable) {
	sft = SmartFeeTable{}
	for _, entry := range fees {
		key := FeeGroup{
			AssetCode: entry.Asset,
			FeeType:   entry.FeeType,
			Subtype:   entry.Subtype,
		}

		sft[key] = append(sft[key], entry)
	}

	return sft
}

func (sft SmartFeeTable) Update(fees []core.FeeEntry) {
	update := NewSmartFeeTable(fees)

	for k, v := range update {
		sft[k] = SmartFillFeeGaps(sft[k], v)
	}
}

func (sft SmartFeeTable) GetValuesByAsset() (byAsset map[string][]core.FeeEntry) {
	byAsset = make(map[string][]core.FeeEntry)
	for key := range sft {
		byAsset[key.AssetCode] = append(byAsset[key.AssetCode], sft[key]...)
	}

	return byAsset
}

func (sft SmartFeeTable) AddZeroFees(assets []core.Asset) {
	for _, asset := range assets {
		for _, ft := range xdr.FeeTypeAll {
			subtypes := []int64{0}
			if ft == xdr.FeeTypePaymentFee {
				subtypes = []int64{
					int64(xdr.PaymentFeeTypeIncoming),
					int64(xdr.PaymentFeeTypeOutgoing),
				}
			}
			for _, st := range subtypes {
				key := FeeGroup{
					AssetCode: asset.Code,
					FeeType:   int(ft),
					Subtype:   st,
				}

				zeroFee := core.FeeEntry{
					Asset:   asset.Code,
					Subtype: st,
					FeeType: int(ft),
				}

				sft[key] = FillFeeGaps(sft[key], zeroFee)
			}
		}
	}
}
