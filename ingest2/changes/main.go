package changes

import "gitlab.com/tokend/go/xdr"

type LedgerChange struct {
	LedgerSeq int32
	LedgerChange xdr.LedgerEntryChange
}
