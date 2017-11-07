package resource

import (
	"bullioncoin.githost.io/development/horizon/db2/history"
)

func (this *HistoryAccount) Populate(row history.Account) {
	this.AccountID = row.Address
}
