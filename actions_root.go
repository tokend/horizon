package horizon

import (
	"gitlab.com/tokend/go/xdr"
	"gitlab.com/tokend/horizon/ledger"
	"gitlab.com/tokend/horizon/render/hal"
	"gitlab.com/tokend/horizon/render/problem"
	"gitlab.com/tokend/horizon/resource"
)

// RootAction provides a summary of the horizon instance and links to various
// useful endpoints
type RootAction struct {
	Action
}

// JSON renders the json response for RootAction
func (action *RootAction) JSON() {
	action.App.UpdateCoreInfo()

	if action.App.CoreInfo == nil {
		action.Err = &problem.ServerOverCapacity
		return
	}

	var res resource.Root
	res.PopulateLedgerState(action.Ctx, ledger.CurrentState())

	res.NetworkPassphrase = action.App.CoreInfo.NetworkPassphrase
	res.AdminAccountID = action.App.CoreInfo.AdminAccountID
	res.MasterExchangeName = action.App.CoreInfo.MasterExchangeName
	res.TxExpirationPeriod = action.App.CoreInfo.TxExpirationPeriod
	res.XDRRevision = xdr.Revision

	hal.Render(action.W, res)
}
