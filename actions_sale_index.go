package horizon

import (
	"fmt"
	"time"

	"gitlab.com/swarmfund/horizon/db2"
	"gitlab.com/swarmfund/horizon/db2/history"
	"gitlab.com/swarmfund/horizon/exchange"
	"gitlab.com/swarmfund/horizon/render/hal"
	"gitlab.com/swarmfund/horizon/render/problem"
	"gitlab.com/swarmfund/horizon/resource"
)

type Sort int64

const (
	SortTypeDefaultPage Sort = iota
	SortTypeMostFounded
	SortTypeByEndTime
	SortTypeByPopularity
	SortTypeStartTime
)

// SaleIndexAction renders slice of reviewable requests
type SaleIndexAction struct {
	Action
	Owner        string
	BaseAsset    string
	OpenOnly     bool
	Upcoming     bool
	SortType     *int64
	Name         string
	Records      []history.Sale
	PagingParams db2.PageQuery
	Page         hal.Page
}

// JSON is a method for actions.JSON
func (action *SaleIndexAction) JSON() {
	action.Do(
		action.EnsureHistoryFreshness,
		action.loadParams,
		action.loadRecord,
		action.loadPage,
		func() {
			hal.Render(action.W, action.Page)
		},
	)
}

func (action *SaleIndexAction) loadParams() {
	action.PagingParams = action.GetPageQuery()
	action.Owner = action.GetString("owner")
	action.BaseAsset = action.GetString("base_asset")
	action.Name = action.GetString("name")

	action.OpenOnly = action.GetBool("open_only")
	action.Upcoming = action.GetBool("upcoming")

	action.SortType = action.GetOptionalInt64("sort_by")
	action.Page.Filters = map[string]string{
		"owner":      action.Owner,
		"base_asset": action.BaseAsset,
		"name":       action.Name,
		"open_only":  action.GetString("open_only"),
		"upcoming":   action.GetString("upcoming"),
	}
}

func (action *SaleIndexAction) loadRecord() {
	q := action.HistoryQ().Sales()

	if action.Owner != "" {
		q = q.ForOwner(action.Owner)
	}

	if action.BaseAsset != "" {
		q = q.ForBaseAsset(action.BaseAsset)
	}

	if action.Name != "" {
		q = q.ForName(action.Name)
	}

	if action.OpenOnly {
		q = q.Open(time.Now().UTC())
	}

	if action.Upcoming {
		q = q.Upcoming(time.Now().UTC())
	}

	sortBy := SortTypeDefaultPage
	if action.SortType != nil {
		sortBy = Sort(*action.SortType)
	}

	switch sortBy {
	case SortTypeDefaultPage:
		q = q.Page(action.PagingParams)
	case SortTypeStartTime:
		q = q.OrderByStartTime()
	case SortTypeByEndTime:
		q = q.OrderByEndTime()
	case SortTypeByPopularity:
		values, err := action.CoreQ().OrderBook().InvestorsCount()
		if err != nil {
			action.Log.WithError(err).Error("Unable to load investors count")
			action.Err = &problem.ServerError
			return
		}
		q = q.OrderByPopularity(values)
	default:
		action.SetInvalidField("sort_by", fmt.Errorf("invalid value %d", sortBy))
		return
	}

	converter, err := exchange.NewConverter(action.CoreQ())
	if err != nil {
		action.Log.WithError(err).Error("Failed to init converter")
		action.Err = &problem.ServerError
		return
	}

	action.Records, err = selectSalesWithCurrentCap(q, converter)
	if err != nil {
		action.Log.WithError(err).Error("failed to load sales")
		action.Err = &problem.ServerError
		return
	}
}

func (action *SaleIndexAction) loadPage() {
	for i := range action.Records {
		var res resource.Sale
		res.Populate(&action.Records[i])
		err := populateSaleWithStats(action.Records[i].ID, &res, action.CoreQ())
		if err != nil {
			action.Log.WithError(err).
				WithField("sale_id", action.Records[i].ID).
				Error("failed to populate stat for sale")
			action.Err = &problem.ServerError
			return
		}
		action.Page.Add(&res)
	}

	// with custom sorting type
	// pagination will not work
	if action.SortType != nil {
		// init set empty slice if no records
		action.Page.Init()
		return
	}

	action.Page.BaseURL = action.BaseURL()
	action.Page.BasePath = action.Path()
	action.Page.Limit = action.PagingParams.Limit
	action.Page.Cursor = action.PagingParams.Cursor
	action.Page.Order = action.PagingParams.Order
	action.Page.PopulateLinks()
}
