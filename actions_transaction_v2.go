package horizon

import (
	"gitlab.com/swarmfund/horizon/db2"
	"gitlab.com/swarmfund/horizon/db2/history"
	"gitlab.com/swarmfund/horizon/render/hal"
	"gitlab.com/swarmfund/horizon/render/problem"
	"gitlab.com/swarmfund/horizon/resource"
	"gitlab.com/distributed_lab/logan/v3/errors"
	"gitlab.com/tokend/regources"
)

// TransactionV2IndexAction: pages of transactions

// TransactionV2IndexAction renders a page of ledger resources, identified by
// a normal page query, entry type and effects
type TransactionV2IndexAction struct {
	Action
	EntryTypeFilter       []int
	EffectFilter          []int
	PagingParams          db2.PageQuery
	TransactionsV2Records []regources.TransactionV2
	// It's guarantied that there is no additional changes
	// which satisfy restriction change_time < NoUpdatesUntilLedger.ClosedAt
	NoUpdatesUntilLedger history.Ledger
	Page                 hal.Page

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

// getTransactionRecords - returns slice of transactions fetched for ledger changes,
// true - if page of records was full, error - if something bad happened
func (action *TransactionV2IndexAction) getTransactionRecords() ([]regources.TransactionV2, bool, error) {
	sortedLedgerChanges, isPageFull, err := action.getLedgerChanges()
	if err != nil {
		return nil, false, errors.Wrap(err, "failed to get ledger changes")
	}

	var transactionsIDs []int64
	for txID := range sortedLedgerChanges {
		transactionsIDs = append(transactionsIDs, txID)
	}

	var transactions []history.Transaction
	err = action.HistoryQ().Transactions().ByTxIDs(transactionsIDs).
		Page(action.PagingParams). // page here is needed to make sure that order of transaction is correct (it will not add or remove any values)
		Select(&transactions)
	if err != nil {
		return nil, false, errors.Wrap(err, "failed to get transactions for ledger changes")
	}

	var result []regources.TransactionV2
	for _, tx := range transactions {
		txV2 := resource.PopulateTransactionV2(tx, sortedLedgerChanges[tx.ID])
		result = append(result, txV2)
	}

	return result, isPageFull, nil
}

func (action *TransactionV2IndexAction) loadRecords() {
	// memorize ledger sequence before select to prevent data race
	noUpdatesUntilLedgerSeq := int32(action.App.historyLatestLedgerGauge.Value())

	var err error
	var isPageFull bool
	action.TransactionsV2Records, isPageFull, err = action.getTransactionRecords()
	if err != nil {
		action.Log.WithError(err).Error("Failed to get transactions v2 records")
		action.Err = &problem.ServerError
		return
	}

	if isPageFull {
		// we fetched full page, probably there is something ahead
		noUpdatesUntilLedgerSeq = action.TransactionsV2Records[len(action.TransactionsV2Records)-1].LedgerSequence
	}

	// load ledger close time
	if err := action.HistoryQ().LedgerBySequence(&action.NoUpdatesUntilLedger, noUpdatesUntilLedgerSeq); err != nil {
		action.Log.WithError(err).Error("failed to get NoUpdatesUntilLedger")
		action.Err = &problem.ServerError
		return
	}
}

// getLedgerChanges - returns map of ledger changes grouped by transaction_id,
// true if page is full, error if something went wrong
func (action *TransactionV2IndexAction) getLedgerChanges() (map[int64][]history.LedgerChanges, bool, error) {
	var ledgerChanges []history.LedgerChanges

	err := action.HistoryQ().LedgerChanges().ByEntryType(action.EntryTypeFilter).
	ByEffects(action.EffectFilter).
	ByTransactionIDs(action.PagingParams, action.EntryTypeFilter, action.EffectFilter).
	Select(&ledgerChanges)

	if err != nil {
		return nil, false, errors.Wrap(err, "failed to select ledger changes")
	}

	sortedLedgerChanges := map[int64][]history.LedgerChanges{}
	for _, change := range ledgerChanges {
		sortedLedgerChanges[change.TransactionID] = append(sortedLedgerChanges[change.TransactionID], change)
	}

	isPageFull := uint64(len(ledgerChanges)) >= action.PagingParams.Limit
	return sortedLedgerChanges, isPageFull, nil
}

func (action *TransactionV2IndexAction) loadPage() {
	for _, record := range action.TransactionsV2Records {
		action.Page.Add(record)
	}

	action.Page.Embedded.Meta = &hal.PageMeta{
		LatestLedger: &hal.LatestLedgerMeta{
			Sequence: action.NoUpdatesUntilLedger.Sequence,
			ClosedAt: action.NoUpdatesUntilLedger.ClosedAt,
		},
	}

	action.Page.BaseURL = action.BaseURL()
	action.Page.BasePath = action.Path()
	action.Page.Limit = action.PagingParams.Limit
	action.Page.Cursor = action.PagingParams.Cursor
	action.Page.Order = action.PagingParams.Order
	action.Page.PopulateLinks()
}