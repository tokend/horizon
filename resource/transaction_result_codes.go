package resource

import (
	"gitlab.com/swarmfund/go/xdr"
	"gitlab.com/swarmfund/horizon/codes"
)

// TransactionResultCodes represent a summary of result codes returned from
// a single xdr TransactionResult
type TransactionResultCodes struct {
	TransactionCode string   `json:"transaction"`
	OperationCodes  []string `json:"operations,omitempty"`
}

// NewTransactionResultCodes creates tx result codes from xdr result
func NewTransactionResultCodes(txResult xdr.TransactionResult) (*TransactionResultCodes, error) {
	txResultCode, opResultCodes, err := codes.ForTxResult(txResult)
	if err != nil {
		return nil, err
	}

	return &TransactionResultCodes{
		TransactionCode: txResultCode,
		OperationCodes:  opResultCodes,
	}, nil
}
