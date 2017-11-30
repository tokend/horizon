package ingest

import (
	"gitlab.com/swarmfund/go/xdr"
	"errors"
)

func assetPairUpdated(is *Session, ledgerEntry *xdr.LedgerEntry) error {
	apEntry := ledgerEntry.Data.AssetPair
	if apEntry == nil {
		return errors.New("expected assetPair not to be nil")
	}

	is.Cursor.PriceHistoryProvider().Put(string(apEntry.Base), string(apEntry.Quote), int64(apEntry.CurrentPrice))
	return nil
}
