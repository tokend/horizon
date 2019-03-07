package txsub

import (
	"gitlab.com/distributed_lab/logan/v3"
	"gitlab.com/distributed_lab/logan/v3/errors"
)

// ErrorType defines specific type of error
type ErrorType int

const (
	// Timeout error during submission
	Timeout ErrorType = iota
	// RejectedTx occurs when Tx was reject by core
	RejectedTx
	//DuplicateTx occurs when Tx with the same hash was already submitted
	DuplicateTx
)

// An Error represents a transaction submission error.
type Error interface {
	error
	// Type -- returns type of the error
	Type() ErrorType
	// ResultXDR returns base64 xdr encoded xdr.TransactionResult
	// Returns empty string for all non rejected txs
	ResultXDR() string
}

var timeoutError = &txSubError{
	error:     errors.New("timeout"),
	errorType: Timeout,
}

type txSubError struct {
	error
	errorType ErrorType
	resultXDR string
}

func (m *txSubError) Type() ErrorType {
	return m.errorType
}

func (m *txSubError) ResultXDR() string {
	return m.resultXDR
}

// NewRejectedTxError creates error with type RejectedTx
func NewRejectedTxError(resultXDR string) Error {

	return &txSubError{
		error: errors.From(errors.New("core rejected tx"), logan.F{
			"result_xdr": resultXDR,
		}),
		errorType: RejectedTx,
		resultXDR: resultXDR,
	}
}

func NewDuplicateTxError(hash string) Error {
	return &txSubError{
		error: errors.From(errors.New("Tx with the same hash already exists"), logan.F{
			"tx_hash": hash,
		}),
		errorType: DuplicateTx,
	}
}
