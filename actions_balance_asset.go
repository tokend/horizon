package horizon

import (
	"database/sql"

	"bullioncoin.githost.io/development/horizon/db2/core"
	"bullioncoin.githost.io/development/horizon/render/hal"
	"bullioncoin.githost.io/development/horizon/render/problem"
	"bullioncoin.githost.io/development/horizon/resource"
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
