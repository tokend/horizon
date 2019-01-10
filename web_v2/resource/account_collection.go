package resource

import "gitlab.com/tokend/horizon/db2/core"

func NewAccountCollection(records []core.Account) *AccountCollection {
	data := make([]AccountData, len(records))

	for i, record := range records {
		data[i] = NewAccount(&record).Data
	}

	return &AccountCollection{Data: data}
}

type AccountCollection struct {
	Data     []AccountData `json:"data"`
	Included []interface{} `json:"included,omitempty"`
}

func (a *AccountCollection) IncludeBalances(balances []BalanceData) {
	for _, balance := range balances {
		a.Included = append(a.Included, balance)
	}
}
