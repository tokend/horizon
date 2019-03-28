package internal

import (
	"time"

	"gitlab.com/tokend/go/xdr"
	"gitlab.com/tokend/regources/rgenerated"
)

//FeeFromXdr - converts xdr fee to regources fee
func FeeFromXdr(data xdr.Fee) rgenerated.Fee {
	return rgenerated.Fee{
		Fixed:             rgenerated.Amount(data.Fixed),
		CalculatedPercent: rgenerated.Amount(data.Percent),
	}
}

//TimeFromXdr - converts xdr.Uint64 to unix utc timestamp
func TimeFromXdr(data xdr.Uint64) time.Time {
	return time.Unix(int64(data), 0).UTC()
}
