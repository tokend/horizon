package changes

import (
	"time"

	"gitlab.com/tokend/go/xdr"
)

//LedgerChange is struct for storing single LedgerEntryChange
// along with ledger details and operation that triggered change
type LedgerChange struct {
	LedgerSeq       int32
	LedgerCloseTime time.Time
	LedgerChange    xdr.LedgerEntryChange
	Operation       *xdr.Operation
}
