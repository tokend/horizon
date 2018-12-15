package changes

import (
	"gitlab.com/distributed_lab/logan/v3/errors"
	"gitlab.com/tokend/go/xdr"
	history "gitlab.com/tokend/horizon/db2/history2"
)

type accountStorage interface {
	InsertAccount(rawAccountID xdr.AccountId, account history.Account) error
	GetAccount(address xdr.AccountId) (history.Account, error)
}

type accountHandler struct {
	storage accountStorage

	accountSeq          int32
	processingLedgerSeq int32
}

func newAccountHandler(storage accountStorage) *accountHandler {
	return &accountHandler{
		storage: storage,
	}
}

func (p *accountHandler) Created(lc ledgerChange) error {
	account := lc.LedgerChange.MustCreated().Data.MustAccount()
	newAccountID := p.nextAccountID(lc.LedgerSeq)
	newAccount := history.NewAccount(newAccountID, account.AccountId.Address(), int32(account.AccountType))
	err := p.storage.InsertAccount(account.AccountId, newAccount)
	if err != nil {
		return errors.Wrap(err, "failed to insert account")
	}
	return nil
}

func (p *accountHandler) nextAccountID(ledgerSeq int32) int64 {
	if p.processingLedgerSeq != ledgerSeq {
		p.processingLedgerSeq = ledgerSeq
		p.accountSeq = 0
	}

	p.accountSeq++
	return int64(ledgerSeq)<<32 | int64(p.accountSeq)
}
