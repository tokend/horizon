package horizon

import (
	"time"

	"fmt"

	"gitlab.com/swarmfund/horizon/db2"
	"gitlab.com/swarmfund/horizon/db2/history"
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
)

// SaleIndexAction renders slice of reviewable requests
type SaleIndexAction struct {
	Action
	Owner          string
	BaseAsset      string
	OpenOnly       bool
	Upcoming       bool
	ReachedSoftCap bool
	NearlyFunded   *int64
	GoalValue      *int64
	SortType       *int64
	Name           string
	Records        []history.Sale
	PagingParams   db2.PageQuery
	Page           hal.Page
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
	action.OpenOnly = action.GetBool("open_only")
	action.Name = action.GetString("name")

	action.NearlyFunded = action.GetOptionalInt64("nearly_funded")
	action.Upcoming = action.GetBool("upcoming")
	action.ReachedSoftCap = action.GetBool("reached_soft_cap")
	action.GoalValue = action.GetOptionalAmount("goal_value")
	action.SortType = action.GetOptionalInt64("sort_by")

	action.Page.Filters = map[string]string{
		"owner":            action.Owner,
		"base_asset":       action.BaseAsset,
		"name":             action.Name,
		"open_only":        action.GetString("open_only"),
		"nearly_funded":    action.GetString("nearly_funded"),
		"upcoming":         action.GetString("upcoming"),
		"reached_soft_cap": action.GetString("reached_soft_cap"),
		"goal_value":       action.GetString("goal_value"),
	}
	fmt.Println("Collect queries")
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

	if action.NearlyFunded != nil && !action.ReachedSoftCap {
		q = q.NearlyFunded(*action.NearlyFunded)
	}

	if action.ReachedSoftCap && action.NearlyFunded == nil {
		q = q.ReachedSoftCap()
	}

	if action.GoalValue != nil {
		q = q.GoalValue(*action.GoalValue)
	}

	sortBy := SortTypeDefaultPage
	if action.SortType != nil {
		sortBy = Sort(*action.SortType)
	}

	switch sortBy {
	case SortTypeDefaultPage:
		q = q.Page(action.PagingParams)
	case SortTypeMostFounded:
		q = q.OrderByCurrentCap(true)
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

	var err error
	action.Records, err = q.Select()
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
