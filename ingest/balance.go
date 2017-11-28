package ingest

import (
	"time"

	"gitlab.com/swarmfund/go/xdr"
	"gitlab.com/swarmfund/horizon/db2/core"
	"github.com/lann/squirrel"
	"github.com/pkg/errors"
)

func balanceUpdated(is *Session, ledgerEntry *xdr.LedgerEntry) error {
	balance := ledgerEntry.Data.Balance
	if balance == nil {
		return errors.New("expected non nil balance")
	}

	if is.Paranoid {
		// seems like we have partial history, ensuring balance exists
		var b core.Balance
		err := is.Ingestion.CoreQ.BalanceByID(&b, balance.BalanceId.AsString())
		if err != nil {
			return errors.Wrap(err, "failed to get balance")
		}
		_, err = is.Ingestion.tryIngestBalance(b.BalanceID, b.Asset, b.AccountID)
		if err != nil {
			return errors.Wrap(err, "failed to ingest balance")
		}
	}

	amount := balance.Amount + balance.Locked

	return is.Ingestion.tryIngestBalanceUpdate(
		balance.BalanceId.AsString(), int64(amount), is.Cursor.Ledger().CloseTime,
	)
}

func balanceCreated(is *Session, ledgerEntry *xdr.LedgerEntry) error {
	balance := ledgerEntry.Data.Balance
	if balance == nil {
		return errors.New("Expected balance not to be nil")
	}

	_, err := is.Ingestion.tryIngestBalance(balance.BalanceId.AsString(),
		string(balance.Asset),
		balance.AccountId.Address())
	if err != nil {
		return errors.Wrap(err, "failed to ingest balance")
	}

	if err := balanceUpdated(is, ledgerEntry); err != nil {
		return errors.Wrap(err, "failed to updated balance")
	}
	return nil
}

func (ingest *Ingestion) tryIngestBalance(
	balanceID, asset, accountID string) (bool, error) {
	result, err := ingest.DB.ExecRaw(`
		insert into history_balances (balance_id, asset, account_id)
		values ($1, $2, $3) on conflict do nothing`,
		balanceID, asset, accountID)
	if err != nil {
		return false, err
	}
	rows, err := result.RowsAffected()
	if err != nil {
		return false, errors.Wrap(err, "failed to get rows affected")
	}
	return rows > 0, nil
}

func (ingest *Ingestion) tryIngestBalanceUpdate(
	balanceID string, amount, closeTime int64) error {
	_, err := ingest.DB.Exec(squirrel.
		Insert("history_balance_updates").
		SetMap(map[string]interface{}{
			"balance_id": balanceID,
			"amount":     amount,
			"updated_at": time.Unix(closeTime, 0).UTC(),
		}))
	return err
}
