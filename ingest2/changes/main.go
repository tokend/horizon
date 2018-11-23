package changes

import (
	"time"

	"gitlab.com/distributed_lab/logan/v3/errors"

	"gitlab.com/tokend/go/xdr"
)

type LedgerChange struct {
	LedgerSeq       int32
	LedgerCloseTime time.Time
	LedgerChange    xdr.LedgerEntryChange
	Operation       *xdr.Operation
}

var (
	balances           = balanceChanges{}
	contracts          = contractChanges{}
	reviewableRequests = reviewableRequestChanges{}
	sales              = saleChanges{}
)

func (l LedgerChange) Created() error {
	switch l.LedgerChange.Created.Data.Type {
	case xdr.LedgerEntryTypeBalance:
		err := balances.Created(l)
		if err != nil {
			return errors.Wrap(err, "failed to process balance created event")
		}
		return nil
	case xdr.LedgerEntryTypeContract:
		err := contracts.Created(l)
		if err != nil {
			return errors.Wrap(err, "failed to process contract created event")
		}
		return nil
	case xdr.LedgerEntryTypeReviewableRequest:
		err := reviewableRequests.Created(l)
		if err != nil {
			return errors.Wrap(err, "failed to process reviewable request created event")
		}
		return nil
	case xdr.LedgerEntryTypeSale:
		err := sales.Created(l)
		if err != nil {
			return errors.Wrap(err, "failed to process sale created event")
		}
		return nil
	default:
		return errors.Errorf("Unrecognized ledger entry type created: %d", l.LedgerChange.Created.Data.Type)
	}
}

func (l LedgerChange) Updated() error {
	switch l.LedgerChange.Updated.Data.Type {
	case xdr.LedgerEntryTypeContract:
		err := contracts.Updated(l)
		if err != nil {
			return errors.Wrap(err, "failed to process contract updated event")
		}
		return nil
	case xdr.LedgerEntryTypeReviewableRequest:
		err := reviewableRequests.Updated(l)
		if err != nil {
			return errors.Wrap(err, "failed to process reviewable request updated event")
		}
		return nil
	case xdr.LedgerEntryTypeSale:
		err := sales.Updated(l)
		if err != nil {
			return errors.Wrap(err, "failed to process sale updated event")
		}
		return nil
	default:
		return errors.Errorf("Unrecognized ledger entry type updated: %d", l.LedgerChange.Updated.Data.Type)
	}
}

func (l LedgerChange) Deleted() error {
	switch l.LedgerChange.Removed.Type {
	case xdr.LedgerEntryTypeContract:
		err := contracts.Deleted(l)
		if err != nil {
			return errors.Wrap(err, "failed to process contract deleted event")
		}
		return nil
	case xdr.LedgerEntryTypeReviewableRequest:
		err := reviewableRequests.Deleted(l)
		if err != nil {
			return errors.Wrap(err, "failed to process reviewable request deleted event")
		}
		return nil
	default:
		return errors.Errorf("Unrecognized ledger entry type deleted: %d", l.LedgerChange.Removed.Type)
	}
}
