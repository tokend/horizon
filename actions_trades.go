package horizon

import (
	"gitlab.com/swarmfund/horizon/db2"
	"gitlab.com/swarmfund/horizon/db2/history"
	"gitlab.com/swarmfund/horizon/render/hal"
	"gitlab.com/swarmfund/horizon/render/problem"
	"gitlab.com/swarmfund/horizon/resource"
	"strconv"
)

type TradesAction struct {
	Action
	BaseAsset    string
	QuoteAsset   string
	PagingParams db2.PageQuery
	OrderBookID uint64

	Trades []history.Trades
	Page   hal.Page
}

// JSON is a method for actions.JSON
func (action *TradesAction) JSON() {
	action.Do(
		action.loadParams,
		action.loadRecords,
		func() {
			hal.Render(action.W, action.Page)
		},
	)
}

func (action *TradesAction) loadParams() {
	action.BaseAsset = action.GetNonEmptyString("base_asset")
	action.QuoteAsset = action.GetNonEmptyString("quote_asset")
	action.OrderBookID = action.GetUInt64("order_book_id")
	action.PagingParams = action.GetPageQuery()
	action.Page.Filters = map[string]string{
		"base_asset":  action.BaseAsset,
		"quote_asset": action.QuoteAsset,
		"order_book_id": strconv.FormatUint(action.OrderBookID, 10),
	}
}

func (action *TradesAction) loadRecords() {
	err := action.HistoryQ().Trades().ForOrderBook(action.OrderBookID).ForPair(action.BaseAsset, action.QuoteAsset).Page(action.PagingParams).Select(&action.Trades)
	if err != nil {
		action.Log.WithError(err).Error("Failed to get trades")
		action.Err = &problem.ServerError
		return
	}

	for i := range action.Trades {
		var result resource.Trades
		result.Populate(&action.Trades[i])
		action.Page.Add(&result)
	}

	action.Page.BaseURL = action.BaseURL()
	action.Page.BasePath = action.Path()
	action.Page.Limit = action.PagingParams.Limit
	action.Page.Cursor = action.PagingParams.Cursor
	action.Page.Order = action.PagingParams.Order
	action.Page.PopulateLinks()

}
