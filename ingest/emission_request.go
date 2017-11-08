package ingest

import (
	"gitlab.com/tokend/go/xdr"
	"github.com/pkg/errors"
)

func emissionRequestCreated(is *Session, ledgerEntry *xdr.LedgerEntry) error {
	emissionRequest := ledgerEntry.Data.CoinsEmissionRequest
	if emissionRequest == nil {
		return errors.New("Expected emission request not to be nil")
	}

	return is.Ingestion.InsertEmissionRequest(is.Cursor.Ledger(), emissionRequest)
}

func emissionRequestDeleted(is *Session, ledgerKey *xdr.LedgerKey) error {
	emissionRequest := ledgerKey.CoinsEmissionRequest
	if emissionRequest == nil {
		return errors.New("Expected emission request key not to be nil")
	}

	return is.Ingestion.DeleteEmissionRequest(
		uint64(emissionRequest.RequestId),
	)
}
