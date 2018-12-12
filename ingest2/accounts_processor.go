package ingest2

import (
	"gitlab.com/distributed_lab/logan/v3"
	"gitlab.com/tokend/go/xdr"
	history "gitlab.com/tokend/horizon/db2/history2"
	"gitlab.com/distributed_lab/logan/v3/errors"
)

type accountsStorage interface {
	InsertAccount(account history.Account) error
}

type accountsProcessor struct {
	log *logan.Entry
	storage accountsStorage
	accounts map[xdr.AccountId]history.Account
}

func (p *accountsProcessor) Consume(it ledgerChangesIteration) (error) {
	lc := it.LedgerChange
	account := lc.MustCreated().Data.MustAccount()
	return p.ConsumeAccountDetails(account.AccountId, int32(account.AccountType), it.LedgerSeq, it.LedgerGlobOpSeq)
}

func (p *accountsProcessor) ConsumeAccountDetails(accountID xdr.AccountId, accountType int32, ledgerSeq, ledgerGlobalOpSeq int32) error {
	newAccount := history.NewAccount(accountID.Address(), accountType, ledgerSeq, ledgerGlobalOpSeq)
	p.accounts[accountID] = newAccount
	err := p.storage.InsertAccount(newAccount)
	if err != nil {
		return errors.Wrap(err, "failed to insert account")
	}
	return nil
}
