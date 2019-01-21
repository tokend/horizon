package internal

import (
	"time"

	"gitlab.com/tokend/go/xdr"
	"gitlab.com/tokend/regources/v2"
)

//FeeFromXdr - converts xdr fee to regources fee
func FeeFromXdr(data xdr.Fee) regources.Fee {
	return regources.Fee{
		Fixed:             regources.Amount(data.Fixed),
		CalculatedPercent: regources.Amount(data.Percent),
	}
}

//ToTime - converts xdr.Uint64 to unix utc timestamp
func ToTime(data xdr.Uint64) time.Time {
	return time.Unix(int64(data), 0).UTC()
}
