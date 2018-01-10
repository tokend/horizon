package ingest

import "gitlab.com/swarmfund/go/xdr"

func (is *Session) operationChanges(changes xdr.LedgerEntryChanges) error {
	for i := range changes {
		var err error

		switch changes[i].Type {
		case xdr.LedgerEntryChangeTypeCreated:
			err = is.operationCreatedEntry(changes[i].Created)
		case xdr.LedgerEntryChangeTypeRemoved:
			err = is.operationDeletedEntry(changes[i].Removed)
		case xdr.LedgerEntryChangeTypeUpdated:
			err = is.operationUpdatedEntry(changes[i].Updated)
		default:
			// nothing to do here
		}

		if err != nil {
			is.log.Error("Failed to process operation changes")
			return err
		}
	}

	return nil
}

func (is *Session) operationCreatedEntry(ledgerEntry *xdr.LedgerEntry) error {
	handler, ok := creationHandlers[ledgerEntry.Data.Type]
	if !ok {
		return nil
	}

	return handler(is, ledgerEntry)
}

func (is *Session) operationDeletedEntry(ledgerKey *xdr.LedgerKey) error {
	handler, ok := deletionHandlers[ledgerKey.Type]
	if !ok {
		return nil
	}

	return handler(is, ledgerKey)
}

func (is *Session) operationUpdatedEntry(ledgerEntry *xdr.LedgerEntry) error {
	handler, ok := updateHandlers[ledgerEntry.Data.Type]
	if !ok {
		return nil
	}

	return handler(is, ledgerEntry)
}

var creationHandlers = map[xdr.LedgerEntryType]func(is *Session, ledgerEntry *xdr.LedgerEntry) error{
	xdr.LedgerEntryTypeBalance:           balanceCreated,
	xdr.LedgerEntryTypeReviewableRequest: reviewableRequestCreate,
	xdr.LedgerEntryTypeSale:              saleCreate,
}

var deletionHandlers = map[xdr.LedgerEntryType]func(is *Session, ledgerKey *xdr.LedgerKey) error{}

var updateHandlers = map[xdr.LedgerEntryType]func(is *Session, ledgerKey *xdr.LedgerEntry) error{
	xdr.LedgerEntryTypeBalance:           balanceUpdated,
	xdr.LedgerEntryTypeReviewableRequest: reviewableRequestUpdate,
	xdr.LedgerEntryTypeSale:              saleUpdate,
}
