package horizon

import (
	"gitlab.com/swarmfund/horizon/db2/core"
	"gitlab.com/swarmfund/horizon/db2/history"
	"gitlab.com/swarmfund/horizon/render/hal"
	"gitlab.com/swarmfund/horizon/render/problem"
	"gitlab.com/swarmfund/horizon/resource"
)

// SaleShowAction renders a sale found by its ID.
type SaleShowAction struct {
	Action
	RequestID uint64
	Record    *history.Sale
	offers    []core.Offer
	balances  []core.Balance
	assetPair *core.AssetPair
	result    resource.Sale
}

// JSON is a method for actions.JSON
func (action *SaleShowAction) JSON() {
	action.Do(
		action.EnsureHistoryFreshness,
		action.loadParams,
		action.loadRecord,
		action.populateResult,
		func() {
			hal.Render(action.W, action.result)
		},
	)
}

func (action *SaleShowAction) loadParams() {
	action.RequestID = action.GetUInt64("id")
}

func (action *SaleShowAction) loadRecord() {
	var err error
	action.Record, err = action.HistoryQ().Sales().ByID(action.RequestID)
	if err != nil {
		action.Log.WithError(err).
			WithField("request_id", action.RequestID).
			Error("failed to load sale")
		action.Err = &problem.ServerError
		return
	}

	if action.Record == nil {
		action.Err = &problem.NotFound
		return
	}

	action.offers = make([]core.Offer, 0)
	err = action.CoreQ().Offers().
		ForOrderBookID(action.Record.ID).Select(&action.offers)
	if err != nil {
		action.Log.WithError(err).
			WithField("sale_id", action.Record.ID).
			Error("failed to load offers for sale")
		action.Err = &problem.ServerError
		return
	}

	action.balances, err = action.CoreQ().Balances().
		ByAsset(action.Record.BaseAsset).Select()
	if err != nil {
		action.Log.WithError(err).
			WithField("sale_id", action.Record.ID).
			Error("failed to load base asset balances for sale")
		action.Err = &problem.ServerError
		return
	}

	action.assetPair, err = action.CoreQ().AssetPairs().
		ByCode(action.Record.BaseAsset, action.Record.QuoteAsset)
	if err != nil {
		action.Log.WithError(err).
			WithField("sale_id", action.Record.ID).
			WithField("base", action.Record.BaseAsset).
			WithField("quote", action.Record.QuoteAsset).
			Error("failed to load asset pair for sale")
		action.Err = &problem.ServerError
		return
	}
}

func (action *SaleShowAction) populateResult() {
	action.result.Populate(action.Record)
	err := action.result.PopulateStat(action.offers, action.balances, action.assetPair)
	if err != nil {
		action.Log.WithError(err).
			WithField("request_id", action.RequestID).
			Error("failed to populate stat for sale")
		action.Err = &problem.ServerError
		return
	}
}
