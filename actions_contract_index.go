package horizon
/*
import (
	"gitlab.com/distributed_lab/logan/v3/errors"
	"gitlab.com/swarmfund/horizon/db2"
	"gitlab.com/swarmfund/horizon/db2/history"
	"gitlab.com/swarmfund/horizon/render/hal"
	"gitlab.com/swarmfund/horizon/render/problem"
	"gitlab.com/swarmfund/horizon/resource"
	"gitlab.com/tokend/regources"
)

// TransactionV2IndexAction: pages of transactions

// TransactionV2IndexAction renders a page of ledger resources, identified by
// a normal page query, entry type and effects
type ContractIndexAction struct {
	Action
	PagingParams     db2.PageQuery
	StartTime        uint64
	EndTime          uint64
	State            uint32
	ContractsRecords []regources.Contract
	Page             hal.Page
}

// JSON is a method for actions.JSON
func (action *ContractIndexAction) JSON() {
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

func (action *ContractIndexAction) checkAllowed() {
	action.IsAllowed("")
}

func (action *ContractIndexAction) loadParams() {
	action.ValidateCursorAsDefault()
	action.StartTime = action.GetUInt64("start_time")
	action.EffectFilter = action.GetIntArray("effect")
	action.PagingParams = action.getTxPageQuery()
}

func (action *ContractIndexAction) getTxPageQuery() db2.PageQuery {
	pagingParams := action.GetPageQuery()
	limit := action.GetUInt64("limit")
	if limit > maxTxPagSize {
		pagingParams.Limit = maxTxPagSize
	}

	return pagingParams
}

// getTransactionRecords - returns slice of transactions fetched for ledger changes,
// true - if page of records was full, error - if something bad happened
func (action *ContractIndexAction) getTransactionRecords() ([]regources.Transaction, bool, error) {
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

	var result []regources.Transaction
	for _, record := range transactions {
		tx := resource.PopulateTransactionV2(record, sortedLedgerChanges[record.ID])
		result = append(result, tx)
	}

	return result, isPageFull, nil
}

func (action *ContractIndexAction) loadRecords() {
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
		noUpdatesUntilLedgerSeq = action.TransactionsV2Records[len(action.TransactionsV2Records)-1].Ledger
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
func (action *ContractIndexAction) getLedgerChanges() (map[int64][]history.LedgerChanges, bool, error) {
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

func (action *ContractIndexAction) loadPage() {
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
*/