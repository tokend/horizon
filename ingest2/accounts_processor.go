package ingest2

import (
	"gitlab.com/distributed_lab/logan/v3"
	"gitlab.com/tokend/go/xdr"
	history "gitlab.com/tokend/horizon/db2/history2"
	"gitlab.com/distributed_lab/logan/v3/errors"
)

type accountsStorageI interface {
	InsertAccounts(accounts []history.Account) error
}

type accountsProcessor struct {
	log *logan.Entry
	accounts map[xdr.AccountId]history.Account
	newAccounts []history.Account
}

func (p *accountsProcessor) Consume(it ledgerChangesIteration) (error) {
	lc := it.LedgerChange
	if lc.Type != xdr.LedgerEntryChangeTypeCreated {
		return nil
	}

	if lc.Created.Data.Type != xdr.LedgerEntryTypeAccount {
		return nil
	}

	newAccount := history.NewAccount(*lc.Created.Data.Account, it.LedgerSeq, it.LedgerGlobOpSeq)
	p.accounts[lc.Created.Data.Account.AccountId] = newAccount
	p.newAccounts = append(p.newAccounts, newAccount)
	return nil
}

func (p *accountsProcessor) Store(storage accountsStorageI) error {
	err := storage.InsertAccounts(p.newAccounts)
	if err != nil {
		return errors.Wrap(err, "failed to insert accounts")
	}

	p.newAccounts = p.newAccounts[0:0]
	return nil
}
