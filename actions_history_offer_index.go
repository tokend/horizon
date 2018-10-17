package horizon

import (
	"fmt"

	"gitlab.com/tokend/horizon/db2"
	"gitlab.com/tokend/horizon/render/hal"
	"gitlab.com/tokend/horizon/render/problem"
	"gitlab.com/tokend/horizon/resource"
)

type OfferState int64

const (
	OfferStateAny OfferState = iota
	OfferStateNoMatches
	OfferStatePartiallyMatched
	OfferStateFullyMatched
	OfferStateCanceled
)

type HistoryOfferIndexAction struct {
	Action
	OwnerID      string
	BaseAsset    string
	QuoteAsset   string
	IsBuy        *bool
	OfferState   OfferState
	PagingParams db2.PageQuery
	Page         hal.Page
}

// JSON is a method for actions.JSON
func (action *HistoryOfferIndexAction) JSON() {
	action.Do(
		action.loadParams,
		action.checkAllowed,
		action.loadRecords,
		action.loadPage,
		func() {
			hal.Render(action.W, action.Page)
		},
	)
}

func (action *HistoryOfferIndexAction) checkAllowed() {
	action.IsAllowed(action.OwnerID)
}

func (action *HistoryOfferIndexAction) loadParams() {
	action.OwnerID = action.GetString("owner_id")
	action.BaseAsset = action.GetString("base_asset")
	action.QuoteAsset = action.GetString("quote_asset")
	action.IsBuy = action.GetOptionalBool("is_buy")

	action.OfferState = OfferState(action.GetInt64("state"))

	action.PagingParams = action.GetPageQuery()
	action.Page.Filters = map[string]string{
		"owner_id":    action.OwnerID,
		"base_asset":  action.BaseAsset,
		"quote_asset": action.QuoteAsset,
		"is_buy":      action.GetString("is_buy"),
	}
}

func (action *HistoryOfferIndexAction) loadRecords() {
	q := action.HistoryQ().Offers()
	if action.OwnerID != "" {
		q = q.ForOwnerID(action.OwnerID)
	}
	if action.BaseAsset != "" {
		q = q.ForBase(action.BaseAsset)
	}
	if action.QuoteAsset != "" {
		q = q.ForQuote(action.QuoteAsset)
	}
	if action.IsBuy != nil {
		q = q.ForIsBuy(*action.IsBuy)
	}

	switch action.OfferState {
	case OfferStateAny:
		// all offers are appropriate
	case OfferStateNoMatches:
		q = q.NoMatches().Canceled(false)
	case OfferStatePartiallyMatched:
		q = q.PartiallyMatched().Canceled(false)
	case OfferStateFullyMatched:
		q = q.FullyMatched()
	case OfferStateCanceled:
		q = q.Canceled(true)
	default:
		action.SetInvalidField("state", fmt.Errorf("invalid value %d", action.OfferState))
		return
	}

	historyOffers, err := q.Page(action.PagingParams).Select()
	if err != nil {
		action.Log.WithError(err).Error("failed to select history offers")
		action.Err = &problem.ServerError
		return
	}

	if len(historyOffers) == 0 {
		action.Err = &problem.NotFound
		return
	}

	for _, offer := range historyOffers {
		action.Page.Add(resource.PopulateHistoryOffer(offer))
	}
}

func (action *HistoryOfferIndexAction) loadPage() {
	action.Page.BaseURL = action.BaseURL()
	action.Page.BasePath = action.Path()
	action.Page.Limit = action.PagingParams.Limit
	action.Page.Cursor = action.PagingParams.Cursor
	action.Page.Order = action.PagingParams.Order
	action.Page.PopulateLinks()
}
