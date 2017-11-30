package ingest

import (
	"github.com/pkg/errors"
	"gitlab.com/swarmfund/go/xdr"
	"gitlab.com/swarmfund/horizon/db2/core"
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
		_, err = is.Ingestion.TryIngestBalance(b.BalanceID, b.Asset, b.AccountID)
		if err != nil {
			return errors.Wrap(err, "failed to ingest balance")
		}
	}

	amount := balance.Amount + balance.Locked

	return is.Ingestion.TryIngestBalanceUpdate(
		balance.BalanceId.AsString(), int64(amount), is.Cursor.Ledger().CloseTime,
	)
}

func balanceCreated(is *Session, ledgerEntry *xdr.LedgerEntry) error {
	balance := ledgerEntry.Data.Balance
	if balance == nil {
		return errors.New("Expected balance not to be nil")
	}

	_, err := is.Ingestion.TryIngestBalance(balance.BalanceId.AsString(),
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
