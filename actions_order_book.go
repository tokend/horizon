package horizon

import (
	"bullioncoin.githost.io/development/horizon/db2/core"
	"bullioncoin.githost.io/development/horizon/render/hal"
	"bullioncoin.githost.io/development/horizon/render/problem"
	"bullioncoin.githost.io/development/horizon/resource"
	"strconv"
)

type OrderBookAction struct {
	Action
	BaseAsset  string
	QuoteAsset string
	IsBuy      bool

	CoreRecords []core.OrderBookEntry
	Page        hal.Page
}

// JSON is a method for actions.JSON
func (action *OrderBookAction) JSON() {
	action.Do(
		action.loadParams,
		action.loadRecords,
		func() {
			hal.Render(action.W, action.Page)
		},
	)
}

func (action *OrderBookAction) loadParams() {
	action.BaseAsset = action.GetNonEmptyString("base_asset")
	action.QuoteAsset = action.GetNonEmptyString("quote_asset")
	action.IsBuy = action.GetBool("is_buy")
	action.Page.Filters = map[string]string{
		"base_asset":  action.BaseAsset,
		"quote_asset": action.QuoteAsset,
		"is_buy":      strconv.FormatBool(action.IsBuy),
	}
}

func (action *OrderBookAction) loadRecords() {
	err := action.CoreQ().OrderBook().ForAssets(action.BaseAsset, action.QuoteAsset).Direction(action.IsBuy).Select(&action.CoreRecords)
	if err != nil {
		action.Log.WithError(err).Error("Failed to get offers from core DB")
		action.Err = &problem.ServerError
		return
	}

	for i := range action.CoreRecords {
		var result resource.OrderBookEntry
		result.Populate(&action.CoreRecords[i], action.BaseAsset, action.QuoteAsset, action.IsBuy)
		action.Page.Add(&result)
	}

	action.Page.BaseURL = action.BaseURL()
	action.Page.BasePath = action.Path()
	action.Page.PopulateLinks()

}
