package resource

import (
	"gitlab.com/tokend/go/amount"
	"gitlab.com/swarmfund/horizon/db2/core"
	"gitlab.com/swarmfund/horizon/db2/history"
)

func (b *BalancePublic) Populate(balance history.Balance) {
	b.BalanceID = balance.BalanceID
	b.AccountID = balance.AccountID
	b.Asset = balance.Asset
}

func (b *Balance) Populate(balance core.Balance) {
	b.BalanceID = balance.BalanceID
	b.AccountID = balance.AccountID
	b.Balance = amount.String(balance.Amount)
	b.Locked = amount.String(balance.Locked)
	b.Asset = balance.Asset
}

func (balance BalancePublic) PagingToken() string {
	return balance.ID
}
