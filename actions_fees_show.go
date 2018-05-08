package horizon

import (
	"gitlab.com/tokend/go/amount"
	"gitlab.com/swarmfund/horizon/db2/core"
	"gitlab.com/swarmfund/horizon/render/hal"
	"gitlab.com/swarmfund/horizon/render/problem"
	"gitlab.com/swarmfund/horizon/resource"
	"github.com/go-errors/errors"
)

// This file contains the actions:
//
// FeesShowAction: renders fees for operationType
type FeesShowAction struct {
	Action
	FeeType   int
	Asset     string
	Subtype   int64
	AccountID string
	Account   *core.Account

	Amount string

	Fee resource.FeeEntry
}

// JSON is a method for actions.JSON
func (action *FeesShowAction) JSON() {
	action.Do(
		action.loadParams,
		action.loadData,
		func() {
			hal.Render(action.W, action.Fee)
		},
	)
}
func (action *FeesShowAction) loadParams() {
	action.FeeType = int(action.GetInt32("fee_type"))
	action.Asset = action.GetNonEmptyString("asset")
	action.Subtype = action.GetInt64("subtype")
	action.AccountID = action.GetString("account")
	action.Amount = action.GetString("amount")
}

func (action *FeesShowAction) loadData() {
	var err error
	if action.AccountID != "" {
		action.Account, err = action.CoreQ().Accounts().ByAddress(action.AccountID)
	}

	if err != nil {
		action.Log.WithError(err).Error("Failed to load account to get fee")
		action.Err = &problem.ServerError
		return
	}

	am := int64(0)
	if action.Amount != "" {
		amXdr, _ := amount.Parse(action.Amount)
		am = int64(amXdr)
	}
	result, err := action.CoreQ().FeeByTypeAssetAccount(action.FeeType, action.Asset, action.Subtype, action.Account, am)
	if err != nil {
		action.Log.WithError(err).Error("Failed to load fee by asset and type")
		action.Err = &problem.ServerError
		return
	}

	if result == nil {
		result = new(core.FeeEntry)
		result.Asset = action.Asset
		result.FeeType = action.FeeType
	}

	percentFee, isOverflow := action.GetPercentFee(result.Percent, action.Amount)
	if isOverflow {
		action.SetInvalidField("amount", errors.New("is too big - overflow"))
		return
	}

	result.Percent = percentFee

	action.Fee.Populate(*result)
}

func (action *FeesShowAction) GetPercentFee(percentFee int64, rawAmount string) (int64, bool) {
	// request does not require to calculate
	if rawAmount == "" {
		return percentFee, false
	}

	am, err := amount.Parse(rawAmount)
	if err != nil {
		action.SetInvalidField("amount", err)
		return 0, false
	}

	return action.CalculatePercentFee(percentFee, int64(am))
}
