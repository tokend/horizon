package resources

import (
	"strconv"

	"gitlab.com/tokend/horizon/db2/history2"
	"gitlab.com/tokend/regources/rgenerated"
)

//NewOperation - creates new instance of rgenerated.Operation
func NewOperation(op history2.Operation) rgenerated.Operation {
	return rgenerated.Operation{
		Key: NewOperationKey(op.ID),
		Attributes: rgenerated.OperationAttributes{
			AppliedAt: op.LedgerCloseTime,
		},
		Relationships: rgenerated.OperationRelationships{
			Tx:     NewTxKey(op.TxID).AsRelation(),
			Source: NewAccountKey(op.Source).AsRelation(),
		},
	}
}

//NewOperationKey - creates new key for operation
func NewOperationKey(opID int64) rgenerated.Key {
	return rgenerated.Key{
		ID:   strconv.FormatInt(opID, 10),
		Type: rgenerated.OPERATIONS,
	}
}
