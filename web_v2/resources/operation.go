package resources

import (
	"strconv"

	"gitlab.com/tokend/horizon/db2/history2"
	"gitlab.com/tokend/regources/v2"
)

//NewOperation - creates new instance of regources.Operation
func NewOperation(op history2.Operation) regources.Operation {
	return regources.Operation{
		Key: NewOperationKey(op.ID),
		Attributes: regources.OperationAttr{
			Details:   op.Details,
			AppliedAt: op.LedgerCloseTime,
		},
		Relationships: regources.OperationRelation{
			Tx:     NewTxKey(op.TxID).AsRelation(),
			Source: NewAccountKey(op.Source).AsRelation(),
		},
	}
}

//NewOperationKey - creates new key for operation
func NewOperationKey(opID int64) regources.Key {
	return regources.Key{
		ID:   strconv.FormatInt(opID, 10),
		Type: regources.TypeOperations,
	}
}
