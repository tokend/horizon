package horizon

import (
	"gitlab.com/swarmfund/horizon/db2"
	"gitlab.com/swarmfund/horizon/db2/history"
	"gitlab.com/swarmfund/horizon/render/hal"
	"gitlab.com/swarmfund/horizon/render/problem"
	"gitlab.com/swarmfund/horizon/resource"
)

// This file contains the actions:
//
// TransactionIndexAction: pages of transactions
// TransactionShowAction: single transaction by sequence, by hash or id

// TransactionIndexAction renders a page of ledger resources, identified by
// a normal page query.
type TransactionV2IndexAction struct {
	Action
	EntryTypeFilter []int
	EffectFilter    []int
	PagingParams    db2.PageQuery
	Records         []history.Transaction
	MetaLedger      history.Ledger
	Page            hal.Page
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

	action.PagingParams = action.GetPageQuery()

	limit := action.GetUInt64("limit")
	if limit > db2.MaxPageSize {
		action.PagingParams.Limit = limit
		if limit > maxTxPagSize {
			action.PagingParams.Limit = maxTxPagSize
		}
	}
}

func (action *TransactionV2IndexAction) loadRecords() {
	q := action.HistoryQ()
	txs := q.Transactions()

	// memorize ledger sequence before select to prevent data race
	latestLedger := int32(action.App.historyLatestLedgerGauge.Value())

	err := txs.Page(action.PagingParams).Select(&action.Records)
	if err != nil {
		action.Log.WithError(err).Error("failed to get transactions")
		action.Err = &problem.ServerError
		return
	}

	if uint64(len(action.Records)) == action.PagingParams.Limit {
		// we fetched full page, probably there is something ahead
		latestLedger = action.Records[len(action.Records)-1].LedgerSequence
	}

	// load ledger close time
	if err := action.HistoryQ().LedgerBySequence(&action.MetaLedger, latestLedger); err != nil {
		action.Log.WithError(err).Error("failed to get ledger")
		action.Err = &problem.ServerError
		return
	}
}

func (action *TransactionV2IndexAction) loadPage() {
	for _, record := range action.Records {
		var res resource.TransactionV2
		res.Populate(action.Ctx, record)
		action.Page.Add(res)
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