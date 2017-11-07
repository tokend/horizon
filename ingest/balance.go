package ingest

import (
	"time"

	"bullioncoin.githost.io/development/go/xdr"
	"bullioncoin.githost.io/development/horizon/db2/core"
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
		_, err = is.Ingestion.tryIngestBalance(
			b.BalanceID, b.Asset, b.AccountID, b.ExchangeID, b.ExchangeName)
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

	exchangeName := is.CoreInfo.MasterExchangeName
	if !balance.Exchange.Equals(is.CoreInfo.MasterAccountIDXDR) {
		exchangeNameP, err := is.Ingestion.CoreQ.ExchangeName(balance.Exchange.Address())
		if err != nil {
			is.log.WithError(err).Error("Failed to get exchange name")
			return err
		}

		if exchangeNameP == nil {
			err = errors.New("Expected exchange name not to be nil")
			is.log.WithError(err).Error("Failed to get exchange name")
			return err
		}

		exchangeName = *exchangeNameP
	}

	_, err := is.Ingestion.tryIngestBalance(
		balance.BalanceId.AsString(), string(balance.Asset),
		balance.AccountId.Address(), balance.Exchange.Address(), exchangeName)
	if err != nil {
		return errors.Wrap(err, "failed to ingest balance")
	}

	if err := balanceUpdated(is, ledgerEntry); err != nil {
		return errors.Wrap(err, "failed to updated balance")
	}
	return nil
}

func (ingest *Ingestion) tryIngestBalance(
	balanceID, asset, accountID, exchangeID, exchangeName string) (bool, error) {
	result, err := ingest.DB.ExecRaw(`
		insert into history_balances (
			balance_id, asset, account_id, exchange_id, exchange_name)
		values ($1, $2, $3, $4, $5) on conflict do nothing`,
		balanceID, asset, accountID, exchangeID, exchangeName)
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
