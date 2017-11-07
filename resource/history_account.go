package resource

import (
	"gitlab.com/distributed_lab/tokend/horizon/db2/history"
)

func (this *HistoryAccount) Populate(row history.Account) {
	this.AccountID = row.Address
}
