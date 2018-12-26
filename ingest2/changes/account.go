package changes

import (
	"gitlab.com/distributed_lab/logan/v3/errors"
	"gitlab.com/tokend/go/xdr"
	history "gitlab.com/tokend/horizon/db2/history2"
)

type accountStorage interface {
	InsertAccount(rawAccountID xdr.AccountId, account history.Account) error
	MustAccount(address xdr.AccountId) (history.Account)
}

type accountHandler struct {
	storage accountStorage
}

func newAccountHandler(storage accountStorage) *accountHandler {
	return &accountHandler{
		storage: storage,
	}
}

//Created - stores new account to storage
func (p *accountHandler) Created(lc ledgerChange) error {
	account := lc.LedgerChange.MustCreated().Data.MustAccount()
	newAccount := history.NewAccount(uint64(account.AccountSeqId), account.AccountId.Address())
	err := p.storage.InsertAccount(account.AccountId, newAccount)
	if err != nil {
		return errors.Wrap(err, "failed to insert account")
	}
	return nil
}
