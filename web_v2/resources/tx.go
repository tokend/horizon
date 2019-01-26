package resources

import (
	"strconv"

	"gitlab.com/tokend/regources/v2"
)

//NewTxKey - creates new Tx Key for specified ID
func NewTxKey(txID int64) regources.Key {
	return regources.Key{
		ID:   strconv.FormatInt(txID, 10),
		Type: regources.TypeTxs,
	}
}
