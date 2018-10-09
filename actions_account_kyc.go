package horizon

import (
	"gitlab.com/tokend/horizon/render/hal"
	"gitlab.com/tokend/horizon/render/problem"
	"gitlab.com/tokend/horizon/resource"
)

type AccountKYCAction struct {
	Action

	AccountID  string
	AccountKYC resource.AccountKYC
}

func (action *AccountKYCAction) JSON() {
	action.Do(
		action.loadParams,
		action.loadData,
		func() {
			hal.Render(action.W, action.AccountKYC)
		},
	)
}

func (action *AccountKYCAction) loadParams() {
	action.AccountID = action.GetNonEmptyString("id")
}

func (action *AccountKYCAction) loadData() {
	accountKYC, err := action.CoreQ().AccountKYC().ByAddress(action.AccountID)

	if err != nil {
		action.Err = &problem.ServerError
		action.Log.WithStack(err).WithError(err).Error("Failed to load account_kyc by account_id")
		return
	}

	if accountKYC == nil {
		action.Err = &problem.NotFound
		action.Log.WithStack(err).WithError(err).Error("account_kyc not found ")
		return
	}

	action.AccountKYC.Populate(*accountKYC)
}
