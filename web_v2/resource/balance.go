package resource

import (
	"gitlab.com/tokend/go/amount"
	"gitlab.com/tokend/horizon/db2/core"
)

func NewBalance(record *core.Balance) *Balance {
	return &Balance{
		Data: BalanceData{
			Id:   record.BalanceID,
			Type: TypeBalances,
			Attributes: &BalanceAttributes{
				Locked:    amount.String(record.Locked),
				Available: amount.String(record.Amount),
			},
		},
	}
}

type Balance struct {
	Data     BalanceData   `json:"data"`
	Included []interface{} `json:"included,omitempty"`
}

type BalanceData struct {
	Id         string             `json:"id"`
	Type       string             `json:"type"`
	Attributes *BalanceAttributes `json:"attributes,omitempty"`
}

type BalanceAttributes struct {
	Available string `json:"available"`
	Locked    string `json:"locked"`
}
