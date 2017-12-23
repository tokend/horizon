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
