package ingest2

import (
	"gitlab.com/distributed_lab/logan/v3"
	"gitlab.com/distributed_lab/logan/v3/errors"
	"gitlab.com/tokend/horizon/db2"
	core "gitlab.com/tokend/horizon/db2/core2"
	"gitlab.com/tokend/horizon/db2/history2"
	"time"
)

type ledgerStorage interface {
	Insert(ledger *history2.Ledger) error
}

// LedgerHandler - stores ledger into
type LedgerHandler struct {
	storage ledgerStorage
}


// NewLedgerHandler - creates new instance of ledger handler
func NewLedgerHandler(storage ledgerStorage) *LedgerHandler {
	return &LedgerHandler{
		storage: storage,
	}
}

func (h *LedgerHandler) Name() string {
	return "ledger_saver"
}

// Handle - stores ledger into db. Returns error if failed to store
func (h *LedgerHandler) Handle(header *core.LedgerHeader, txs []core.Transaction) error {
	hHeader := history2.Ledger{
		TotalOrderID: db2.TotalOrderID{
			ID: int64(header.Sequence),
		},
		Sequence:     header.Sequence,
		Hash:         header.LedgerHash,
		PreviousHash: header.PrevHash,
		ClosedAt:     time.Unix(header.CloseTime, 0).UTC(),
		TxCount:      int32(len(txs)),
	}

	err := h.storage.Insert(&hHeader)
	if err != nil {
		return errors.Wrap(err, "failed to insert ledger", logan.F{
			"ledger_seq": header.Sequence,
		})
	}

	return nil
}
