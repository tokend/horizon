package horizon

import (
	"time"

	"gitlab.com/swarmfund/horizon/db2"
	"gitlab.com/swarmfund/horizon/db2/history"
	"gitlab.com/swarmfund/horizon/render/hal"
	"gitlab.com/swarmfund/horizon/render/problem"
	"gitlab.com/swarmfund/horizon/resource"
)

// SaleIndexAction renders slice of reviewable requests
type SaleIndexAction struct {
	Action
	Owner        string
	BaseAsset    string
	OpenOnly     bool
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
	action.OpenOnly = action.GetBool("open_only")
	action.Name = action.GetString("name")

	action.Page.Filters = map[string]string{
		"owner":      action.Owner,
		"base_asset": action.BaseAsset,
		"name":       action.Name,
		"open_only":  action.GetString("open_only"),
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

	q = q.Page(action.PagingParams)
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

	action.Page.BaseURL = action.BaseURL()
	action.Page.BasePath = action.Path()
	action.Page.Limit = action.PagingParams.Limit
	action.Page.Cursor = action.PagingParams.Cursor
	action.Page.Order = action.PagingParams.Order
	action.Page.PopulateLinks()
}
