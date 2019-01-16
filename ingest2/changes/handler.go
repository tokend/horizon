package changes

import (
	"time"

	"gitlab.com/distributed_lab/logan/v3"
	"gitlab.com/distributed_lab/logan/v3/errors"
	"gitlab.com/tokend/go/xdr"
	core "gitlab.com/tokend/horizon/db2/core2"
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

// Handler - handles ledger changes to populate changes for entries
type Handler struct {
	Create map[xdr.LedgerEntryType]creatable
	Update map[xdr.LedgerEntryType]updatable
	Remove map[xdr.LedgerEntryType]removable
}

//NewHandler - returns new instance of handler
func NewHandler(account accountStorage, balance balanceStorage, request reviewableRequestStorage,
	sale saleStorage) *Handler {

	reviewRequestHandlerInst := newReviewableRequestHandler(request)
	saleHandlerInst := newSaleHandler(sale)

	return &Handler{
		Create: map[xdr.LedgerEntryType]creatable{
			xdr.LedgerEntryTypeAccount:           newAccountHandler(account),
			xdr.LedgerEntryTypeBalance:           newBalanceHandler(account, balance),
			xdr.LedgerEntryTypeReviewableRequest: reviewRequestHandlerInst,
			xdr.LedgerEntryTypeSale:              saleHandlerInst,
		},
		Update: map[xdr.LedgerEntryType]updatable{
			xdr.LedgerEntryTypeReviewableRequest: reviewRequestHandlerInst,
			xdr.LedgerEntryTypeSale:              saleHandlerInst,
		},
		Remove: map[xdr.LedgerEntryType]removable{
			xdr.LedgerEntryTypeReviewableRequest: reviewRequestHandlerInst,
			xdr.LedgerEntryTypeSale:              saleHandlerInst,
		},
	}
}

// Handle - processes all the ledger changes for specified ledger
func (h *Handler) Handle(header *core.LedgerHeader, txs []core.Transaction) error {
	for txI := range txs {
		tx := txs[txI]
		ops := tx.ResultMeta.MustOperations()
		for opI := range ops {
			op := ops[opI]
			for changeI := range op.Changes {
				change := op.Changes[changeI]
				err := h.handle(ledgerChange{
					LedgerSeq:       header.Sequence,
					LedgerCloseTime: time.Unix(header.CloseTime, 0).UTC(),
					LedgerChange:    change,
					Operation:       &tx.Envelope.Tx.Operations[opI],
					OperationResult: tx.Result.Result.Result.MustResults()[opI].Tr,
				})

				if err != nil {
					return errors.Wrap(err, "failed to process ledger change", logan.F{
						"ledger_seq": header.Sequence,
						"tx_i":       txI,
						"op_i":       opI,
						"change_i":   changeI,
						"change":     change,
					})
				}
			}
		}
	}

	return nil
}

// handle - tries to find corresponding ledger change handler and handle it.
func (h *Handler) handle(lc ledgerChange) error {
	switch lc.LedgerChange.Type {
	case xdr.LedgerEntryChangeTypeCreated:
		return h.created(lc)
	case xdr.LedgerEntryChangeTypeUpdated:
		return h.updated(lc)
	case xdr.LedgerEntryChangeTypeRemoved:
		return h.removed(lc)
	case xdr.LedgerEntryChangeTypeState:
		return nil
	default:
		return errors.From(errors.New("Unrecognized ledger entry change type"), logan.F{
			"change_type": lc.LedgerChange.Type,
		})
	}
}

func (h *Handler) created(lc ledgerChange) error {
	handler, ok := h.Create[lc.LedgerChange.Created.Data.Type]
	if !ok {
		return nil
	}

	return handler.Created(lc)
}

func (h *Handler) updated(lc ledgerChange) error {
	handler, ok := h.Update[lc.LedgerChange.Updated.Data.Type]
	if !ok {
		return nil
	}

	return handler.Updated(lc)
}

func (h *Handler) removed(lc ledgerChange) error {
	handler, ok := h.Remove[lc.LedgerChange.Removed.Type]
	if !ok {
		return nil
	}

	return handler.Removed(lc)
}

//Name - name of the handler
func (h *Handler) Name() string {
	return "ledger_changes_handler"
}
