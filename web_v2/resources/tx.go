package resources

import (
	"strconv"

	"gitlab.com/tokend/regources/rgenerated"
)

//NewTxKey - creates new Tx Key for specified ID
func NewTxKey(txID int64) rgenerated.Key {
	return rgenerated.Key{
		ID:   strconv.FormatInt(txID, 10),
		Type: rgenerated.TRANSACTIONS,
	}
}
