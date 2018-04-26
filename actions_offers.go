package horizon

import (
	"github.com/go-errors/errors"
	"gitlab.com/swarmfund/horizon/db2/core"
	"gitlab.com/swarmfund/horizon/render/hal"
	"gitlab.com/swarmfund/horizon/render/problem"
	"gitlab.com/swarmfund/horizon/resource"
)

type OffersAction struct {
	Action
	AccountID         string
	BaseAsset         string
	QuoteAsset        string
	IsBuy             *bool
	OfferID           string
	OrderBookID       *uint64
	OnlyPrimaryMarket bool

	CoreRecords []core.Offer
	Page        hal.Page
}

// JSON is a method for actions.JSON
func (action *OffersAction) JSON() {
	action.Do(
		action.loadParams,
		action.checkAllowed,
		action.loadRecords,
		func() {
			hal.Render(action.W, action.Page)
		},
	)
}

func (action *OffersAction) loadParams() {
	action.AccountID = action.GetNonEmptyString("account_id")
	action.BaseAsset = action.GetString("base_asset")
	action.QuoteAsset = action.GetString("quote_asset")
	action.IsBuy = action.GetOptionalBool("is_buy")
	action.OfferID = action.GetString("offer_id")
	action.OrderBookID = action.GetOptionalUint64("order_book_id")
	if (action.BaseAsset == "") != (action.QuoteAsset == "") {
		action.SetInvalidField("base_asset", errors.New("base and quote assets must be both set or both not set"))
		return
	}

	action.OnlyPrimaryMarket = action.GetBool("only_primary")
}

func (action *OffersAction) checkAllowed() {
	action.IsAllowed(action.AccountID)
}

func (action *OffersAction) loadRecords() {
	q := action.CoreQ().Offers().ForAccount(action.AccountID)
	if action.BaseAsset != "" {
		q = q.ForAssets(action.BaseAsset, action.QuoteAsset)
	}

	if action.IsBuy != nil {
		q = q.IsBuy(*action.IsBuy)
	}

	if action.OfferID != "" {
		q = q.ForOfferID(action.OfferID)
	}

	if action.OrderBookID != nil {
		q = q.ForOrderBookID(*action.OrderBookID)
	}

	if action.OnlyPrimaryMarket {
		q = q.OnlyPrimaryMarket()
	}

	err := q.Select(&action.CoreRecords)
	if err != nil {
		action.Log.WithError(err).Error("Failed to get offers from core DB")
		action.Err = &problem.ServerError
		return
	}

	action.Page.Init()
	for i := range action.CoreRecords {
		var result resource.Offer
		result.Populate(&action.CoreRecords[i])
		action.Page.Add(&result)
	}

}
