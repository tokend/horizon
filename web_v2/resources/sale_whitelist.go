package resources

import (
	"fmt"

	regources "gitlab.com/tokend/regources/generated"
)

func NewSaleWhitelistKey(saleID uint64) regources.Key {
	return regources.Key{
		ID:   fmt.Sprintf("%d", saleID),
		Type: regources.SALE_WHITELIST,
	}
}

func NewSaleWhitelist(saleID uint64, address string) regources.SaleWhitelist {
	return regources.SaleWhitelist{
		Key: NewSaleWhitelistKey(saleID),
		Relationships: regources.SaleWhitelistRelationships{
			Participant: NewAccountKey(address).AsRelation(),
		},
	}
}
