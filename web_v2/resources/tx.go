package resources

import (
	"strconv"

	regources "gitlab.com/tokend/regources/v2/generated"
)

//NewTxKey - creates new Tx Key for specified ID
func NewTxKey(txID int64) regources.Key {
	return regources.Key{
		ID:   strconv.FormatInt(txID, 10),
		Type: regources.TRANSACTIONS,
	}
}
