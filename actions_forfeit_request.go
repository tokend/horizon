package horizon

import (
	"bullioncoin.githost.io/development/go/amount"
	"bullioncoin.githost.io/development/go/xdr"
	"gitlab.com/distributed_lab/tokend/horizon/db2/core"
	"gitlab.com/distributed_lab/tokend/horizon/render/hal"
	"gitlab.com/distributed_lab/tokend/horizon/render/problem"
	"gitlab.com/distributed_lab/tokend/horizon/resource/base"
	"github.com/go-errors/errors"
)

type ForfeitRequestAction struct {
	Action
	Account *core.Account
	Amount  int64
	Asset   string

	Response base.ForfeitRequest
}

func (action *ForfeitRequestAction) JSON() {
	action.Do(
		action.loadParams,
		action.calculateAssetForms,
		func() {
			hal.Render(action.W, action.Response)
		},
	)
}

func (action *ForfeitRequestAction) loadParams() {
	action.Account = action.GetCoreAccount("account_id", action.CoreQ())
	action.Amount = action.GetAmount("amount")
	action.Asset = action.GetString("asset")
	if action.Amount == 0 {
		action.SetInvalidField("amount", errors.New("Must not be 0"))
		return
	}
}

func (action *ForfeitRequestAction) calculateAssetForms() {
	asset, err := action.CoreQ().AssetByCode(action.Asset)
	if err != nil {
		action.Log.WithError(err).Error("Failed to get asset by code")
		action.Err = &problem.ServerError
		return
	}

	if asset == nil {
		action.SetInvalidField("asset", &problem.NotFound)
		return
	}

	var totalPercent, totalFixed int64
	for _, assetForm := range asset.AssetForms {
		items := action.Amount / int64(assetForm.Unit)
		if items == 0 {
			continue
		}

		action.Amount = action.Amount % int64(assetForm.Unit)

		percent, fixed, err := action.calculateFeeForAssetForm(items, assetForm)
		if err != nil {
			action.Log.WithError(err).Error("Failed to calculate fee for asset form")
			action.Err = &problem.ServerError
			return
		}

		var resultAssetForm base.ForfeitItem
		resultAssetForm.Populate(string(assetForm.Name), int64(assetForm.Unit), items)
		action.Response.Forms = append(action.Response.Forms, resultAssetForm)

		totalPercent += percent
		totalFixed += fixed
	}

	action.Response.TotalFixedFee = amount.String(totalFixed)
	action.Response.TotalPercentFee = amount.String(totalPercent)
}

func (action *ForfeitRequestAction) calculateFeeForAssetForm(items int64, assetForm xdr.AssetForm) (int64, int64, error) {
	unit := int64(assetForm.Unit)
	formAmount := items * unit
	feeEntry, err := action.CoreQ().FeeByTypeAssetAccount(int(xdr.FeeTypeForfeitFee), action.Asset, unit, action.Account, formAmount)
	if err != nil {
		return 0, 0, err
	}

	if feeEntry == nil {
		return 0, 0, nil
	}

	percentFee, isOverflow := action.CalculatePercentFee(feeEntry.Percent, formAmount)
	if isOverflow {
		return 0, 0, errors.New("Overflow")
	}

	return percentFee, feeEntry.Fixed, nil
}
