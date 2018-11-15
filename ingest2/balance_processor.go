package ingest2

import (
	history "gitlab.com/tokend/horizon/db2/history2"
	"gitlab.com/distributed_lab/logan/v3"
	"gitlab.com/tokend/go/xdr"
	"gitlab.com/distributed_lab/logan/v3/errors"
)

type accountProvider interface {
	GetAccount(address xdr.AccountId) (history.Account, error)
}

type balancesStorage interface {
	InsertBalances(balances []history.Balance) error
}

type balancesProcessor struct {
	log *logan.Entry
	balances map[xdr.BalanceId]history.Balance
	newBalances []history.Balance

	accountProvider

	balanceSeq int32
	processingLedgerSeq int32
}

func (p *balancesProcessor) Consume(it ledgerChangesIteration) (error) {
	lc := it.LedgerChange
	if lc.Type != xdr.LedgerEntryChangeTypeCreated {
		return nil
	}

	if lc.Created.Data.Type != xdr.LedgerEntryTypeBalance {
		return nil
	}

	balance := lc.Created.Data.MustBalance()

	account, err := p.accountProvider.GetAccount(balance.AccountId)
	if err != nil {
		return errors.Wrap(err, "failed to get account for address", logan.Field("address", balance.AccountId.Address()))
	}

	balanceSeq := p.nextBalanceSeq(it.LedgerSeq)
	newBalance := history.NewBalance(it.LedgerSeq, balanceSeq, account.ID, balance)
	p.balances[balance.BalanceId] = newBalance
	p.newBalances = append(p.newBalances, newBalance)
	return nil
}

func (p *balancesProcessor) nextBalanceSeq(ledgerSeq int32) int32 {
	if p.processingLedgerSeq != ledgerSeq {
		p.processingLedgerSeq = ledgerSeq
		p.balanceSeq = 0
	}

	p.balanceSeq++

	return p.balanceSeq
}

func (p *balancesProcessor) Store(storage balancesStorage) error {
	err := storage.InsertBalances(p.newBalances)
	if err != nil {
		return errors.Wrap(err, "failed to insert balances")
	}

	p.newBalances = p.newBalances[0:0]
	return nil
}
