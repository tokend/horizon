package changes

import (
	"gitlab.com/distributed_lab/logan/v3"
	"gitlab.com/distributed_lab/logan/v3/errors"
	"gitlab.com/tokend/go/xdr"
)

type creatable interface {
	Created(change ledgerChange) error
}
type updatable interface {
	Updated(change ledgerChange) error
}
type removable interface {
	Removed(change ledgerChange) error
}

// Handler - handles ledger changes
type Handler struct {
	Create map[xdr.LedgerEntryType]creatable
	Update map[xdr.LedgerEntryType]updatable
	Remove map[xdr.LedgerEntryType]removable
}

func NewHandler(account accountStorage, balance balanceStorage, contract contractStorage, request reviewableRequestStorage, sale saleStorage) *Handler {
	contractHandler := newContractHandler(contract, request)
	reviewableRequestHandler := newReviewableRequestHandler(request)
	saleHandler := newSaleHandler(sale)

	return &Handler{
		Create: map[xdr.LedgerEntryType]creatable{
			xdr.LedgerEntryTypeAccount:           newAccountHandler(account),
			xdr.LedgerEntryTypeBalance:           newBalanceHandler(account, balance),
			xdr.LedgerEntryTypeContract:          contractHandler,
			xdr.LedgerEntryTypeReviewableRequest: reviewableRequestHandler,
			xdr.LedgerEntryTypeSale:              saleHandler,
		},
		Update: map[xdr.LedgerEntryType]updatable{
			xdr.LedgerEntryTypeContract:          contractHandler,
			xdr.LedgerEntryTypeReviewableRequest: reviewableRequestHandler,
			xdr.LedgerEntryTypeSale:              saleHandler,
		},
		Remove: map[xdr.LedgerEntryType]removable{
			xdr.LedgerEntryTypeContract:          contractHandler,
			xdr.LedgerEntryTypeReviewableRequest: reviewableRequestHandler,
		},
	}
}

// Handle - tries to find corresponding ledger change handler and handle it.
func (c *Handler) Handle(lc ledgerChange) error {
	switch lc.LedgerChange.Type {
	case xdr.LedgerEntryChangeTypeCreated:
		return c.Created(lc)
	case xdr.LedgerEntryChangeTypeUpdated:
		return c.Updated(lc)
	case xdr.LedgerEntryChangeTypeRemoved:
		return c.Removed(lc)
	case xdr.LedgerEntryChangeTypeState:
		return nil
	default:
		return errors.From(errors.New("Unrecognized ledger entry change type"), logan.F{
			"change_type": lc.LedgerChange.Type,
		})
	}
}

func (c *Handler) Created(lc ledgerChange) error {
	handler, ok := c.Create[lc.LedgerChange.Created.Data.Type]
	if !ok {
		return nil
	}

	return handler.Created(lc)
}

func (c *Handler) Updated(lc ledgerChange) error {
	handler, ok := c.Update[lc.LedgerChange.Updated.Data.Type]
	if !ok {
		return nil
	}

	return handler.Updated(lc)
}

func (c *Handler) Removed(lc ledgerChange) error {
	handler, ok := c.Remove[lc.LedgerChange.Removed.Type]
	if !ok {
		return nil
	}

	return handler.Removed(lc)
}
