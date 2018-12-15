package changes

import (
	"time"

	"gitlab.com/tokend/go/xdr"
)

type ledgerChange struct {
	LedgerSeq       int32
	LedgerCloseTime time.Time
	LedgerChange    xdr.LedgerEntryChange
	Operation       *xdr.Operation
}

func unixToTime(t int64) time.Time {
	return time.Unix(t, 0).UTC()
}
