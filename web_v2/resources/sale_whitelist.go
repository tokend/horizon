package resources

import (
	"fmt"

	regources "gitlab.com/tokend/regources/generated"
)

func NewSaleWhitelistKey(saleID uint64) regources.Key {
	return regources.Key{
		ID:   fmt.Sprintf("%d", saleID),
		Type: regources.SALES,
	}
}
