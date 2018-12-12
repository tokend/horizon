package changes

import (
	"gitlab.com/distributed_lab/logan/v3"
	"gitlab.com/distributed_lab/logan/v3/errors"
	"gitlab.com/tokend/go/xdr"
	"gitlab.com/tokend/horizon/db2/history2"
)

type accountProvider interface {
	GetAccount(address xdr.AccountId) (history2.Account, error)
}

type balanceStorage interface {
	InsertBalance(balance history2.Balance) error
}

type balanceChanges struct {
	storage         balanceStorage
	accountProvider accountProvider
	balances        map[xdr.BalanceId]history2.Balance

	balanceSeq          int32
	processingLedgerSeq int32
}

func (c *balanceChanges) Created(lc LedgerChange) error {
	balance := lc.LedgerChange.MustCreated().Data.MustBalance()
	account, err := c.accountProvider.GetAccount(balance.AccountId)
	if err != nil {
		return errors.Wrap(err, "failed to get account for address", logan.F{
			"address": balance.AccountId.Address(),
		})
	}

	balanceSeq := c.nextBalanceSeq(lc.LedgerSeq)
	newBalance := history2.NewBalance(lc.LedgerSeq, balanceSeq, account.ID, balance)
	c.balances[balance.BalanceId] = newBalance
	err = c.storage.InsertBalance(newBalance)
	if err != nil {
		return errors.Wrap(err, "failed to insert balance", logan.F{
			"balance_address": newBalance.BalanceAddress,
		})
	}
	return nil
}

func (c *balanceChanges) nextBalanceSeq(ledgerSeq int32) int32 {
	if c.processingLedgerSeq != ledgerSeq {
		c.processingLedgerSeq = ledgerSeq
		c.balanceSeq = 0
	}

	c.balanceSeq++

	return c.balanceSeq
}
