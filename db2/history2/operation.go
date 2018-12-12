package history2

import (
	"time"

	"gitlab.com/tokend/go/xdr"
)

// Operation - stores details of operation performed
type Operation struct {
	ID               int64
	Type             xdr.OperationType
	OperationDetails OperationDetails
	LedgerCloseTime  time.Time
	Source           string
}
