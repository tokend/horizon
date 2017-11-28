package horizon

import (
	"database/sql"

	"gitlab.com/swarmfund/horizon/db2/core"
	"gitlab.com/swarmfund/horizon/render/hal"
	"gitlab.com/swarmfund/horizon/render/problem"
	"gitlab.com/swarmfund/horizon/resource"
)

type BalanceAssetAction struct {
	Action
	BalanceID string
	Asset     resource.Asset
}

func (action *BalanceAssetAction) JSON() {
	action.Do(
		action.loadParams,
		action.loadData,
		func() {
			hal.Render(action.W, action.Asset)
		},
	)
}

func (action *BalanceAssetAction) loadParams() {
	action.BalanceID = action.GetBalanceIDAsString("balance_id")
}

func (action *BalanceAssetAction) loadData() {
	var result core.Balance
	err := action.CoreQ().BalanceByID(&result, action.BalanceID)
	if err != nil {
		if err == sql.ErrNoRows {
			action.Err = &problem.NotFound
			return
		}

		action.Err = &problem.ServerError
		action.Log.WithError(err).Error("Failed to get balance by id")
		return
	}

	action.Asset.Code = result.Asset
}
