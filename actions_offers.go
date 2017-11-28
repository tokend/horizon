package horizon

import (
	"strconv"

	"gitlab.com/swarmfund/horizon/db2"
	"gitlab.com/swarmfund/horizon/db2/core"
	"gitlab.com/swarmfund/horizon/render/hal"
	"gitlab.com/swarmfund/horizon/render/problem"
	"gitlab.com/swarmfund/horizon/resource"
	"github.com/go-errors/errors"
)

type OffersAction struct {
	Action
	AccountID    string
	BaseAsset    string
	QuoteAsset   string
	IsBuy        *bool
	PagingParams db2.PageQuery
	OfferID      string

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
	if (action.BaseAsset == "") != (action.QuoteAsset == "") {
		action.SetInvalidField("base_asset", errors.New("base and quote assets must be both set or both not set"))
		return
	}
	action.PagingParams = action.GetPageQuery()
	action.Page.Filters = map[string]string{
		"offer_id":    action.OfferID,
		"base_asset":  action.BaseAsset,
		"quote_asset": action.QuoteAsset,
	}

	if action.IsBuy != nil {
		action.Page.Filters["is_buy"] = strconv.FormatBool(*action.IsBuy)
	}
}

func (action *OffersAction) checkAllowed() {
	action.IsAllowed(action.AccountID)
}

func (action *OffersAction) loadRecords() {
	q := action.CoreQ().Offers().ForAccount(action.AccountID)
	if action.BaseAsset != "" {
		q.ForAssets(action.BaseAsset, action.QuoteAsset)
	}

	if action.IsBuy != nil {
		q.IsBuy(*action.IsBuy)
	}

	if action.OfferID != "" {
		q.ForOfferID(action.OfferID)
	}

	q = q.Page(action.PagingParams)

	err := q.Select(&action.CoreRecords)
	if err != nil {
		action.Log.WithError(err).Error("Failed to get offers from core DB")
		action.Err = &problem.ServerError
		return
	}

	for i := range action.CoreRecords {
		var result resource.Offer
		result.Populate(&action.CoreRecords[i])
		action.Page.Add(&result)
	}

	action.Page.BaseURL = action.BaseURL()
	action.Page.BasePath = action.Path()
	action.Page.Limit = action.PagingParams.Limit
	action.Page.Cursor = action.PagingParams.Cursor
	action.Page.Order = action.PagingParams.Order
	action.Page.PopulateLinks()

}
