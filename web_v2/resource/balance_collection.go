package resource

import "gitlab.com/tokend/horizon/db2/core"

func NewBalanceCollection(records []core.Balance) *BalanceCollection {
	data := make([]BalanceData, len(records))

	for i, record := range records {
		data[i] = NewBalance(&record).Data
	}

	return &BalanceCollection{Data: data}
}

type BalanceCollection struct {
	Data     []BalanceData `json:"data"`
	Included []interface{} `json:"included,omitempty"`
}

func (b *BalanceCollection) AsRelation() *BalanceCollection {
	data := make([]BalanceData, len(b.Data))

	for _, balance := range b.Data {
		data = append(data, BalanceData{Id: balance.Id, Type: balance.Type})
	}

	return &BalanceCollection{
		Data: data,
	}
}
