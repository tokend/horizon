package resource

import (
	"gitlab.com/tokend/go/amount"
	core "gitlab.com/tokend/horizon/db2/core2"
)

// BalanceResponse - JSON:API response for Balance resource
type BalanceResponse struct {
	Data     Balance       `json:"data"`
	Included []interface{} `json:"included,omitempty"`
}

func NewBalanceResponse(record *core.Balance) *BalanceResponse {
	return &BalanceResponse{
		Data: NewBalance(record),
	}
}

// Balance - resource object representing BalanceEntry
type Balance struct {
	Key
	Attributes    *BalanceAttributes    `json:"attributes,omitempty"`
	Relationships *BalanceRelationships `json:"relationships"`
}

func NewBalance(record *core.Balance) Balance {
	return Balance{
		Key: Key{
			ID:   record.BalanceAddress,
			Type: typeBalances,
		},
		Attributes: &BalanceAttributes{
			Locked:    amount.String(record.Locked),
			Available: amount.String(record.Amount),
		},
		Relationships: &BalanceRelationships{
			Asset: &Key{
				ID:   record.Code,
				Type: typeAssets,
			},
		},
	}
}

// BalanceAttributes - represents information about Balance
type BalanceAttributes struct {
	Available string `json:"available"`
	Locked    string `json:"locked"`
}

//BalanceRelationships -represents reference from Account to other resource objects
type BalanceRelationships struct {
	Asset *Key `json:"asset,omitempty"`
}
