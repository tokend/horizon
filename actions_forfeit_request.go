package horizon

import (
	"gitlab.com/tokend/horizon/db2/core"
	"gitlab.com/tokend/horizon/render/hal"
	"gitlab.com/tokend/horizon/resource/base"
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