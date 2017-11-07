package ingest

import (
	"encoding/json"
	"fmt"
	"time"

	"bullioncoin.githost.io/development/go/xdr"
	"bullioncoin.githost.io/development/horizon/db2/core"
	"bullioncoin.githost.io/development/horizon/db2/history"
	"bullioncoin.githost.io/development/horizon/db2/sqx"
	"bullioncoin.githost.io/development/horizon/ingest/participants"
	"github.com/guregu/null"
	sq "github.com/lann/squirrel"
	"github.com/pkg/errors"
)

// Clear removes data from the ledger
func (ingest *Ingestion) Clear(start int64, end int64) error {
	clear := ingest.DB.DeleteRange
	err := clear(start, end, "history_operation_participants", "history_operation_id")
	if err != nil {
		return err
	}
	err = clear(start, end, "history_operations", "id")
	if err != nil {
		return err
	}
	err = clear(start, end, "history_transaction_participants", "history_transaction_id")
	if err != nil {
		return err
	}
	err = clear(start, end, "history_transactions", "id")
	if err != nil {
		return err
	}
	err = clear(start, end, "history_ledgers", "id")
	if err != nil {
		return err
	}
	err = clear(start, end, "history_emission_requests", "id")
	if err != nil {
		return err
	}
	err = clear(start, end, "history_payment_requests", "id")
	if err != nil {
		return err
	}

	return nil
}

// Close finishes the current transaction and finishes this ingestion.
func (ingest *Ingestion) Close() error {
	return ingest.commit()
}

// Flush writes the currently buffered rows to the db, and if successful
// starts a new transaction.
func (ingest *Ingestion) Flush() error {
	err := ingest.commit()
	if err != nil {
		return err
	}

	return ingest.Start()
}

// Ledger adds a ledger to the current ingestion
func (ingest *Ingestion) Ledger(
	id int64,
	header *core.LedgerHeader,
	txs int,
	ops int,
) error {

	sql := ingest.ledgers.Values(
		CurrentVersion,
		id,
		header.Sequence,
		header.LedgerHash,
		null.NewString(header.PrevHash, header.Sequence > 1),
		0, // TODO remove
		0, // TODO remove
		header.Data.BaseFee,
		header.Data.BaseReserve,
		header.Data.MaxTxSetSize,
		time.Unix(header.CloseTime, 0).UTC(),
		txs,
		ops,
	)

	_, err := ingest.DB.Exec(sql)
	if err != nil {
		return err
	}

	return nil
}

// Operation ingests the provided operation data into a new row in the
// `history_operations` table
func (ingest *Ingestion) Operation(
	id int64,
	txid int64,
	order int32,
	source xdr.AccountId,
	typ xdr.OperationType,
	details map[string]interface{},
	ledgerCloseTime int64,
	identifier uint64,
	state int,
) error {
	djson, err := json.Marshal(details)
	if err != nil {
		return err
	}

	sql := ingest.operations.Values(id,
		txid, order, source.Address(),
		typ, djson, time.Unix(ledgerCloseTime, 0).UTC(),
		identifier, state)
	_, err = ingest.DB.Exec(sql)
	if err != nil {
		return err
	}

	return nil
}

// OperationParticipants ingests the provided accounts `aids` as participants of
// operation with id `op`, creating a new row in the
// `history_operation_participants` table.
func (ingest *Ingestion) OperationParticipants(op int64, opParticipants []participants.Participant) error {
	sql := ingest.operation_participants
	for _, opParticipant := range opParticipants {
		var djson *[]byte
		var err error
		if opParticipant.Details != nil {
			marshalledDetails, err := json.Marshal(opParticipant.Details)
			djson = &marshalledDetails
			if err != nil {
				return err
			}
		}
		haid, err := ingest.getParticipantID(opParticipant.AccountID)
		if err != nil {
			return err
		}
		sql = sql.Values(op, haid, opParticipant.BalanceID.AsString(), djson)
	}

	_, err := ingest.DB.Exec(sql)
	if err != nil {
		return err
	}

	return nil
}

// Rollback aborts this ingestions transaction
func (ingest *Ingestion) Rollback() (err error) {
	err = ingest.DB.Rollback()
	return
}

// Start makes the ingestion reeady, initializing the insert builders and tx
func (ingest *Ingestion) Start() (err error) {
	err = ingest.DB.Begin()
	if err != nil {
		return
	}

	ingest.createInsertBuilders()

	return
}

// Transaction ingests the provided transaction data into a new row in the
// `history_transactions` table
func (ingest *Ingestion) Transaction(
	ledger *core.LedgerHeader,
	id int64,
	tx *core.Transaction,
	fee *core.TransactionFee,
) error {

	sql := ingest.transactions.Values(
		id,
		tx.TransactionHash,
		tx.LedgerSequence,
		tx.Index,
		tx.SourceAddress(),
		tx.Salt(),
		tx.Fee(),
		len(tx.Envelope.Tx.Operations),
		tx.EnvelopeXDR(),
		tx.ResultXDR(),
		tx.ResultMetaXDR(),
		fee.ChangesXDR(),
		sqx.StringArray(tx.Base64Signatures()),
		ingest.formatTimeBounds(tx.Envelope.Tx.TimeBounds),
		tx.MemoType(),
		tx.Memo(),
		time.Unix(ledger.CloseTime, 0),
	)

	_, err := ingest.DB.Exec(sql)
	if err != nil {
		return err
	}

	return nil
}

// TransactionParticipants ingests the provided account ids as participants of
// transaction with id `tx`, creating a new row in the
// `history_transaction_participants` table.
func (ingest *Ingestion) TransactionParticipants(tx int64, aids []xdr.AccountId) error {
	sql := ingest.transaction_participants

	for _, aid := range aids {
		haid, err := ingest.getParticipantID(aid)
		if err != nil {
			return err
		}
		sql = sql.Values(tx, haid)
	}

	_, err := ingest.DB.Exec(sql)
	if err != nil {
		return err
	}

	return nil
}

func (ingest *Ingestion) createInsertBuilders() {
	ingest.ledgers = sq.Insert("history_ledgers").Columns(
		"importer_version",
		"id",
		"sequence",
		"ledger_hash",
		"previous_ledger_hash",
		"total_coins",
		"fee_pool",
		"base_fee",
		"base_reserve",
		"max_tx_set_size",
		"closed_at",
		"transaction_count",
		"operation_count",
	)

	ingest.transactions = sq.Insert("history_transactions").Columns(
		"id",
		"transaction_hash",
		"ledger_sequence",
		"application_order",
		"account",
		"salt",
		"fee_paid",
		"operation_count",
		"tx_envelope",
		"tx_result",
		"tx_meta",
		"tx_fee_meta",
		"signatures",
		"time_bounds",
		"memo_type",
		"memo",
		"ledger_close_time",
	)

	ingest.transaction_participants = sq.Insert("history_transaction_participants").Columns(
		"history_transaction_id",
		"history_account_id",
	)

	ingest.operations = sq.Insert("history_operations").Columns(
		"id",
		"transaction_id",
		"application_order",
		"source_account",
		"type",
		"details",
		"ledger_close_time",
		"identifier",
		"state",
	)

	ingest.operation_participants = sq.Insert("history_operation_participants").Columns(
		"history_operation_id",
		"history_account_id",
		"balance_id",
		"effects",
	)

	ingest.trades = sq.Insert("history_trades").Columns(
		"base_asset",
		"quote_asset",
		"base_amount",
		"quote_amount",
		"price",
		"created_at",
	)

	ingest.priceHistory = sq.Insert("history_price").Columns("base_asset", "quote_asset", "timestamp", "price")

	ingest.emission_requests = sq.Insert("history_emission_requests").Columns(
		"request_id",
		"reference",
		"issuer",
		"receiver",
		"amount",
		"asset",
		"approved",
		"created_at",
		"updated_at",
	)

	ingest.payment_requests = sq.Insert("history_payment_requests").Columns(
		"payment_id",
		"exchange",
		"accepted",
		"details",
		"created_at",
		"updated_at",
		"request_type",
	)

}

func (ingest *Ingestion) commit() error {
	err := ingest.DB.Commit()
	if err != nil {
		return err
	}

	return nil
}

func (ingest *Ingestion) tryIngestAccount(aid string) (result int64, err error) {
	q := history.Q{Repo: ingest.DB}
	var existing history.Account
	err = q.AccountByAddress(&existing, aid)

	if err != nil && !q.NoRows(err) {
		return
	}

	// already imported, return the found value
	if !q.NoRows(err) {
		result = existing.ID
		return
	}

	coreQ := core.Q{Repo: ingest.CoreDB}
	account, err := coreQ.Accounts().ByAddress(aid)
	if err != nil {
		return 0, errors.Wrap(err, "failed to get core balance")
	}
	err = ingest.DB.GetRaw(
		&result,
		`INSERT INTO history_accounts (address, account_type) VALUES (?, ?) RETURNING id`,
		aid, account.AccountType,
	)
	return
}

func (ingest *Ingestion) getParticipantID(
	aid xdr.AccountId,
) (int64, error) {
	return ingest.tryIngestAccount(aid.Address())
}

func (ingest *Ingestion) formatTimeBounds(bounds xdr.TimeBounds) interface{} {
	return sq.Expr("?::int8range", fmt.Sprintf("[%d,%d]", bounds.MinTime, bounds.MaxTime))
}
