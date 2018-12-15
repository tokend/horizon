package changes

import (
	"gitlab.com/distributed_lab/logan/v3"
	"gitlab.com/distributed_lab/logan/v3/errors"
	"gitlab.com/tokend/go/xdr"
	"gitlab.com/tokend/horizon/db2/history2"
)


type balanceStorage interface {
	InsertBalance(rawID xdr.BalanceId, balance history2.Balance) error
}

type balanceHandler struct {
	storage         balanceStorage
	accountStorage accountStorage

	balanceSeq          int32
	processingLedgerSeq int32
}

func newBalanceHandler(accountStorage accountStorage, storage balanceStorage) (*balanceHandler){
	return &balanceHandler{
		storage:storage,
		accountStorage:accountStorage,
	}
}

func (c *balanceHandler) Created(lc ledgerChange) error {
	balance := lc.LedgerChange.MustCreated().Data.MustBalance()
	account, err := c.accountStorage.GetAccount(balance.AccountId)
	if err != nil {
		return errors.Wrap(err, "failed to get account for address", logan.F{
			"address": balance.AccountId.Address(),
		})
	}

	newBalanceID := c.nextBalanceID(lc.LedgerSeq)
	newBalance := history2.NewBalance(newBalanceID, account.ID, balance)
	err = c.storage.InsertBalance(balance.BalanceId, newBalance)
	if err != nil {
		return errors.Wrap(err, "failed to insert balance", logan.F{
			"balance_address": newBalance.BalanceAddress,
		})
	}
	return nil
}

func (c *balanceHandler) nextBalanceID(ledgerSeq int32) int64 {
	if c.processingLedgerSeq != ledgerSeq {
		c.processingLedgerSeq = ledgerSeq
		c.balanceSeq = 0
	}

	c.balanceSeq++

	// we should never end up in situation when there is more than 16777216 new balances in one ledger
	return int64(ledgerSeq)<<24 | int64(c.balanceSeq)
}
