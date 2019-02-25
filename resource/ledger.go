package resource

import (
	"time"

	"gitlab.com/tokend/go/amount"
	"gitlab.com/tokend/horizon/db2/history"
	"golang.org/x/net/context"
)

// Ledger represents a single closed ledger
type Ledger struct {
	ID               string    `json:"id"`
	PT               string    `json:"paging_token"`
	Hash             string    `json:"hash"`
	PrevHash         string    `json:"prev_hash,omitempty"`
	Sequence         int32     `json:"sequence"`
	TransactionCount int32     `json:"transaction_count"`
	OperationCount   int32     `json:"operation_count"`
	ClosedAt         time.Time `json:"closed_at"`
	TotalCoins       string    `json:"total_coins"`
	FeePool          string    `json:"fee_pool"`
	BaseFee          int32     `json:"base_fee"`
	BaseReserve      string    `json:"base_reserve"`
	MaxTxSetSize     int32     `json:"max_tx_set_size"`
}

func (this *Ledger) Populate(ctx context.Context, row history.Ledger) {
	this.ID = row.LedgerHash
	this.PT = row.PagingToken()
	this.Hash = row.LedgerHash
	this.PrevHash = row.PreviousLedgerHash.String
	this.Sequence = row.Sequence
	this.TransactionCount = row.TransactionCount
	this.OperationCount = row.OperationCount
	this.ClosedAt = row.ClosedAt
	this.TotalCoins = amount.String(row.TotalCoins)
	this.FeePool = amount.String(row.FeePool)
	this.BaseFee = row.BaseFee
	this.BaseReserve = amount.String(int64(row.BaseReserve))
	this.MaxTxSetSize = row.MaxTxSetSize
	return
}

func (this Ledger) PagingToken() string {
	return this.PT
}
