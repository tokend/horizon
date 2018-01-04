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
}

// JSON is a method for actions.JSON
func (action *SaleShowAction) JSON() {
	action.Do(
		action.EnsureHistoryFreshness,
		action.loadParams,
		action.loadRecord,
		func() {
			var res resource.Sale
			res.Populate(action.Record)
			res.PopulateStatistic(action.offers)
			hal.Render(action.W, res)
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
		action.Log.WithError(err).WithField("request_id", action.RequestID).Error("failed to load sale")
		action.Err = &problem.ServerError
		return
	}

	if action.Record == nil {
		action.Err = &problem.NotFound
		return
	}

	action.offers = make([]core.Offer, 0)
	err = action.CoreQ().
		Offers().
		ForOrderBookID(action.Record.ID).
		Select(&action.offers)
}
