package resources

import "gitlab.com/tokend/regources/v2"

//NewSaleKey - creates new Key for asset
func NewSaleKey(id int64) regources.Key {
	return regources.NewKeyInt64(id, regources.TypeSales)
}
