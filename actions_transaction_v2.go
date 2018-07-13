package horizon

import (
	"gitlab.com/swarmfund/horizon/db2"
	"gitlab.com/swarmfund/horizon/db2/history"
	"gitlab.com/swarmfund/horizon/render/hal"
	"gitlab.com/swarmfund/horizon/render/problem"
	"gitlab.com/swarmfund/horizon/resource"
	"strconv"
	"gitlab.com/distributed_lab/logan/v3/errors"
)

// This file contains the actions:
//
// TransactionIndexAction: pages of transactions
// TransactionShowAction: single transaction by sequence, by hash or id

// TransactionIndexAction renders a page of ledger resources, identified by
// a normal page query.
type TransactionV2IndexAction struct {
	Action
	EntryTypeFilter       []int
	EffectFilter          []int
	PagingParams          db2.PageQuery
	TransactionsRecords   []history.Transaction
	LedgerChangesRecords  []history.LedgerChanges
	TransactionsV2Records []resource.TransactionV2
	MetaLedger            history.Ledger
	Page                  hal.Page

}

// JSON is a method for actions.JSON
func (action *TransactionV2IndexAction) JSON() {
	action.Do(
		action.EnsureHistoryFreshness,
		action.loadParams,
		action.ValidateCursorWithinHistory,
		action.loadRecords,
		action.loadPage,
		func() {
			hal.Render(action.W, action.Page)
		},
	)
}

func (action *TransactionV2IndexAction) loadParams() {
	action.ValidateCursorAsDefault()
	var err error
	action.GetAddress()
	entryTypeFilterStr := action.BaseURL().Query()["entry_type"]
	action.Log.Info(entryTypeFilterStr)
	action.EntryTypeFilter, err = getIntArrayFromStringArray(entryTypeFilterStr)
	if err != nil {
		action.Log.WithError(err).Error("failed to get entry type filter")
		action.Err = &problem.BadRequest
		return
	}

	effectFilterStr := action.BaseURL().Query()["effect"]
	action.EffectFilter, err = getIntArrayFromStringArray(effectFilterStr)
	if err != nil {
		action.Log.WithError(err).Error("failed to get effect filter")
		action.Err = &problem.BadRequest
		return
	}

	action.Log.Info(action.EffectFilter)

	action.PagingParams = action.GetPageQuery()

	limit := action.GetUInt64("limit")
	if limit > db2.MaxPageSize {
		action.PagingParams.Limit = limit
		if limit > maxTxPagSize {
			action.PagingParams.Limit = maxTxPagSize
		}
	}
}

func getIntArrayFromStringArray(input []string) (result []int, err error) {
	for _, str := range input {
		value, err := strconv.Atoi(str)
		if err != nil {
			return nil, errors.New("failed to convert entry type")
		}

		result = append(result, value)
	}

	return
}

func (action *TransactionV2IndexAction) loadRecords() {
	q := action.HistoryQ()

	// memorize ledger sequence before select to prevent data race
	latestLedger := int32(action.App.historyLatestLedgerGauge.Value())

	lc := q.LedgerChanges()

	var ledgerChanges []history.LedgerChanges

	err := lc.ByEntryType(action.EntryTypeFilter).
		ByEffects(action.EffectFilter).
		ByTransactionIDs(action.PagingParams, action.EntryTypeFilter, action.EffectFilter).
		Select(&ledgerChanges)

	if err != nil {
		action.Log.WithError(err).Error("failed to get ledger changes")
		action.Err = &problem.ServerError
		return
	}

	 sortedLedgerChanges := map[int64][]history.LedgerChanges{}

	for _, change := range ledgerChanges {
		sortedLedgerChanges[change.TransactionID] = append(sortedLedgerChanges[change.TransactionID], change)
	}

	for txID, changes := range sortedLedgerChanges {
		var tx history.Transaction
		action.Err = q.TransactionByHashOrID(&tx, string(txID))
		transactionV2 := resource.TransactionV2{}
		transactionV2.Populate(tx, changes)
		action.TransactionsV2Records = append(action.TransactionsV2Records, transactionV2)
	}

	// load ledger close time
	if err := action.HistoryQ().LedgerBySequence(&action.MetaLedger, latestLedger); err != nil {
		action.Log.WithError(err).Error("failed to get ledger")
		action.Err = &problem.ServerError
		return
	}
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