package horizon

import (
	"gitlab.com/swarmfund/horizon/db2"
	"gitlab.com/swarmfund/horizon/db2/history"
	"gitlab.com/swarmfund/horizon/render/hal"
	"gitlab.com/swarmfund/horizon/render/problem"
	"gitlab.com/swarmfund/horizon/resource"
	"gitlab.com/distributed_lab/logan/v3/errors"
)

// TransactionV2IndexAction: pages of transactions

// TransactionV2IndexAction renders a page of ledger resources, identified by
// a normal page query, entry type and effects
type TransactionV2IndexAction struct {
	Action
	EntryTypeFilter       []int
	EffectFilter          []int
	PagingParams          db2.PageQuery
	TransactionsV2Records []resource.TransactionV2
	MetaLedger            history.Ledger
	Page                  hal.Page

}

// JSON is a method for actions.JSON
func (action *TransactionV2IndexAction) JSON() {
	action.Do(
		action.EnsureHistoryFreshness,
		action.checkAllowed,
		action.loadParams,
		action.ValidateCursorWithinHistory,
		action.loadRecords,
		action.loadPage,
		func() {
			hal.Render(action.W, action.Page)
		},
	)
}

func (action *TransactionV2IndexAction) checkAllowed() {
	action.IsAllowed("")
}

func (action *TransactionV2IndexAction) loadParams() {
	action.ValidateCursorAsDefault()
	action.EntryTypeFilter = action.GetIntArray("entry_type")
	action.EffectFilter = action.GetIntArray("effect")
	action.PagingParams = action.getTxPageQuery()
}

func (action *TransactionV2IndexAction) getTxPageQuery() db2.PageQuery {
	pagingParams := action.GetPageQuery()
	limit := action.GetUInt64("limit")
	if limit > maxTxPagSize {
		pagingParams.Limit = maxTxPagSize
	}

	return pagingParams
}

func (action *TransactionV2IndexAction) loadRecords() {
	// memorize ledger sequence before select to prevent data race
	latestLedger := int32(action.App.historyLatestLedgerGauge.Value())

	sortedLedgerChanges, err := action.getLedgerChanges()
	if err != nil {
		action.Log.WithError(err).Error("failed to get ledger changes")
		action.Err = &problem.ServerError
		return
	}

	var transactionsIDs []int64
	for txID := range sortedLedgerChanges {
		transactionsIDs = append(transactionsIDs, txID)
	}

	var transactions []history.Transaction
	err = action.HistoryQ().Transactions().ByTxIDs(transactionsIDs).Select(&transactions)
	if err != nil {
		action.Log.WithError(err).Error("failed to get transactions")
		action.Err = &problem.ServerError
		return
	}

	for _, tx := range transactions {
		transactionV2 := resource.TransactionV2{}
		transactionV2.Populate(tx, sortedLedgerChanges[tx.ID])
		action.TransactionsV2Records = append(action.TransactionsV2Records, transactionV2)
	}

	if uint64(len(transactions)) == action.PagingParams.Limit {
		// we fetched full page, probably there is something ahead
		latestLedger = transactions[len(transactions)-1].LedgerSequence
	}

	// load ledger close time
	if err := action.HistoryQ().LedgerBySequence(&action.MetaLedger, latestLedger); err != nil {
		action.Log.WithError(err).Error("failed to get ledger")
		action.Err = &problem.ServerError
		return
	}
}

func (action *TransactionV2IndexAction) getLedgerChanges() (map[int64][]history.LedgerChanges, error) {
	var ledgerChanges []history.LedgerChanges

	err := action.HistoryQ().LedgerChanges().ByEntryType(action.EntryTypeFilter).
	ByEffects(action.EffectFilter).
	ByTransactionIDs(action.PagingParams, action.EntryTypeFilter, action.EffectFilter).
	Select(&ledgerChanges)

	if err != nil {
		return nil, errors.Wrap(err, "failed to select ledger changes")
	}

	sortedLedgerChanges := map[int64][]history.LedgerChanges{}
	for _, change := range ledgerChanges {
		sortedLedgerChanges[change.TransactionID] = append(sortedLedgerChanges[change.TransactionID], change)
	}

	return sortedLedgerChanges, nil
}

func (action *TransactionV2IndexAction) loadPage() {
	for _, record := range action.TransactionsV2Records {
		action.Page.Add(record)
	}

	action.Page.Embedded.Meta = &hal.PageMeta{
		LatestLedger: &hal.LatestLedgerMeta{
			Sequence: action.MetaLedger.Sequence,
			ClosedAt: action.MetaLedger.ClosedAt,
		},
	}

	action.Page.BaseURL = action.BaseURL()
	action.Page.BasePath = action.Path()
	action.Page.Limit = action.PagingParams.Limit
	action.Page.Cursor = action.PagingParams.Cursor
	action.Page.Order = action.PagingParams.Order
	action.Page.PopulateLinks()
}