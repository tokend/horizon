// Package resource contains the type definitions for all of horizons
// response resources.
package resource

import (
	"time"

	"gitlab.com/swarmfund/horizon/db2/history"
	"gitlab.com/swarmfund/horizon/render/hal"
	"gitlab.com/swarmfund/horizon/resource/operations"
	"golang.org/x/net/context"
)

// AccountThresholds represents an accounts "thresholds", the numerical values
// needed to satisfy the authorization of a given operation.
type AccountThresholds struct {
	LowThreshold  byte `json:"low_threshold"`
	MedThreshold  byte `json:"med_threshold"`
	HighThreshold byte `json:"high_threshold"`
}

type BalancePublic struct {
	ID        string `json:"id,omitempty"`
	BalanceID string `json:"balance_id"`
	AccountID string `json:"account_id"`
	Asset     string `json:"asset"`
}

// Balance represents an account's holdings for a single currency type
type Balance struct {
	BalancePublic
	Balance          string `json:"balance,omitempty"`
	Locked           string `json:"locked,omitempty"`
	RequireReview    bool   `json:"require_review"`
	IncentivePerCoin string `json:"incentive_per_coin"`
}

// HistoryAccount is a simple resource, used for the account collection actions.
// It provides only the "TotalOrderID" of the account and its account id.
type HistoryAccount struct {
	ID        string `json:"id,omitempty"`
	PT        string `json:"paging_token,omitempty"`
	AccountID string `json:"account_id"`
}

// Ledger represents a single closed ledger
type Ledger struct {
	Links struct {
		Self         hal.Link `json:"self"`
		Transactions hal.Link `json:"transactions"`
		Operations   hal.Link `json:"operations"`
		Payments     hal.Link `json:"payments"`
	} `json:"_links"`
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

// Transaction represents a single, successful transaction
type Transaction struct {
	Links struct {
		Self       hal.Link `json:"self"`
		Account    hal.Link `json:"account"`
		Ledger     hal.Link `json:"ledger"`
		Operations hal.Link `json:"operations"`
		Precedes   hal.Link `json:"precedes"`
		Succeeds   hal.Link `json:"succeeds"`
	} `json:"_links"`
	ID              string    `json:"id"`
	PT              string    `json:"paging_token"`
	Hash            string    `json:"hash"`
	Ledger          int32     `json:"ledger"`
	LedgerCloseTime time.Time `json:"created_at"`
	Account         string    `json:"source_account"`
	FeePaid         int32     `json:"fee_paid"`
	OperationCount  int32     `json:"operation_count"`
	EnvelopeXdr     string    `json:"envelope_xdr"`
	ResultXdr       string    `json:"result_xdr"`
	ResultMetaXdr   string    `json:"result_meta_xdr"`
	FeeMetaXdr      string    `json:"fee_meta_xdr"`
	MemoType        string    `json:"memo_type"`
	Memo            string    `json:"memo,omitempty"`
	Signatures      []string  `json:"signatures"`
	ValidAfter      string    `json:"valid_after,omitempty"`
	ValidBefore     string    `json:"valid_before,omitempty"`
}

// TransactionSuccess represents the result of a successful transaction
// submission.
type TransactionSuccess struct {
	Links struct {
		Transaction hal.Link `json:"transaction"`
	} `json:"_links"`
	Hash   string `json:"hash"`
	Ledger int32  `json:"ledger"`
	Env    string `json:"envelope_xdr"`
	Result string `json:"result_xdr"`
	Meta   string `json:"result_meta_xdr"`
}

// NewOperation returns a resource of the appropriate sub-type for the provided
// operation record.
func NewOperation(
	ctx context.Context,
	row history.Operation,
	participants []*history.Participant,
) (result hal.Pageable, err error) {
	return operations.New(ctx, row, participants, false)
}

func NewPublicOperation(
	ctx context.Context, row history.Operation, participants []*history.Participant,
) (result hal.Pageable, err error) {
	return operations.New(ctx, row, participants, true)
}
