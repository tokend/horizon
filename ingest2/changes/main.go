package changes

import (
	"time"

	"gitlab.com/tokend/go/xdr"
)

type LedgerChange struct {
	LedgerSeq       int32
	LedgerCloseTime time.Time
	LedgerChange    xdr.LedgerEntryChange
}
