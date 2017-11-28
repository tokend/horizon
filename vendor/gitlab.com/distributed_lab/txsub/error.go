package txsub

import (
	"fmt"
	"github.com/go-errors/errors"
)

// ErrorType defines specific type of error
type ErrorType int

const (
	// Timeout error during submission
	Timeout ErrorType = iota
	// RejectedTx occurs when Tx was reject by core
	RejectedTx
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
		error:     fmt.Errorf("core rejected tx: %s", resultXDR),
		errorType: RejectedTx,
		resultXDR: resultXDR,
	}
}
