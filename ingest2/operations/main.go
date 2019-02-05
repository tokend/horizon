// Package operations providers handler which converts xdr operations into details suitable for client side applications
// participants of the operations and effects occurred for their balances
package operations

import (
	"gitlab.com/tokend/go/xdr"
	"gitlab.com/tokend/horizon/db2/history2"
)

//go:generate mockery -case underscore -name IDProvider -inpkg -testonly
type IDProvider interface {
	// MustAccountID returns int value which corresponds to xdr.AccountId
	MustAccountID(raw xdr.AccountId) uint64
	// MustBalanceID returns int value which corresponds to xdr.BalanceId
	MustBalanceID(raw xdr.BalanceId) uint64
}

//go:generate mockery -case underscore -name balanceProvider -inpkg -testonly
type balanceProvider interface {
	// MustBalance returns history balance struct for specific balance id
	MustBalance(balanceID xdr.BalanceId) history2.Balance
}

// handler used to get info and changes from success operation
type handler interface {
	// CreatorDetails returns db suitable operation details,
	// returns error if operation has not existing action (union switch)
	Details(op rawOperation, opRes xdr.OperationResultTr) (history2.OperationDetails, error)
	// ParticipantsEffects returns slice of participant effects of each participants
	// that was affected by operation, can include effects (changes) on participants balances
	ParticipantsEffects(opBody xdr.OperationBody, opRes xdr.OperationResultTr,
		source history2.ParticipantEffect, ledgerChanges []xdr.LedgerEntryChange,
	) ([]history2.ParticipantEffect, error)
}

// rawOperation - inner struct to pass source with operation body
// as one parameter in CreatorDetails method
type rawOperation struct {
	Source xdr.AccountId
	Body   xdr.OperationBody
}
