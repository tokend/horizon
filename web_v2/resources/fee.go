package resources

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"

	"gitlab.com/tokend/go/xdr"

	"gitlab.com/tokend/horizon/db2/core2"
	"gitlab.com/tokend/regources/rgenerated"
)

//NewFee - creates new instance of Fee from provided one.
func NewFee(record core2.Fee) rgenerated.FeeRecord {
	hash := CalculateFeeHash(record)
	return rgenerated.FeeRecord{
		Key: rgenerated.Key{
			ID:   hash,
			Type: rgenerated.FEES,
		},
		Attributes: rgenerated.FeeRecordAttributes{
			Fixed:   rgenerated.Amount(record.Fixed),
			Percent: rgenerated.Amount(record.Percent),
			AppliedTo: rgenerated.FeeAppliedTo{
				Asset:           record.Asset,
				Subtype:         record.Subtype,
				FeeType:         record.FeeType,
				FeeTypeExtended: xdr.FeeType(record.FeeType),
				LowerBound:      rgenerated.Amount(record.LowerBound),
				UpperBound:      rgenerated.Amount(record.UpperBound),
			},
		},
	}
}

//NewFeeKey - creates new Key for fee
func NewFeeKey(hash string) rgenerated.Key {
	return rgenerated.Key{
		ID:   hash,
		Type: rgenerated.FEES,
	}
}

func NewCalculatedFeeKey(hash string) rgenerated.Key {
	return rgenerated.Key{
		ID:   hash,
		Type: rgenerated.CALCULATED_FEE,
	}
}

func CalculateFeeHash(fee core2.Fee) string {
	lowerBound := rgenerated.Amount(fee.LowerBound)
	upperBound := rgenerated.Amount(fee.UpperBound)
	data := fmt.Sprintf("type:%d:asset:%s:subtype:%d:lower_bound:%s:upper_bound:%s",
		fee.FeeType,
		fee.Asset,
		fee.Subtype,
		lowerBound.String(),
		upperBound.String(),
	)

	if fee.AccountID != "" {
		data += fmt.Sprintf("account_id:%s", fee.AccountID)
	} else {
		data += fmt.Sprintf("account_role:%d", fee.AccountRole)
	}

	rawHash := sha256.Sum256([]byte(data))

	return hex.EncodeToString(rawHash[:])
}
