package horizon

import (
	"bullioncoin.githost.io/development/horizon/ledger"
	"bullioncoin.githost.io/development/horizon/render/hal"
	"bullioncoin.githost.io/development/horizon/resource"
)

// RootAction provides a summary of the horizon instance and links to various
// useful endpoints
type RootAction struct {
	Action
}

// JSON renders the json response for RootAction
func (action *RootAction) JSON() {
	action.App.UpdateStellarCoreInfo()

	var res resource.Root
	res.PopulateLedgerState(action.Ctx, ledger.CurrentState())
	res.NetworkPassphrase = action.App.CoreInfo.NetworkPassphrase
	res.CommissionAccountID = action.App.CoreInfo.CommissionAccountID
	res.MasterAccountID = action.App.CoreInfo.MasterAccountID
	res.OperationalAccountID = action.App.CoreInfo.OperationalAccountID
	res.StorageFeeAccountID = action.App.CoreInfo.StorageFeeManageAccountID
	res.MasterExchangeName = action.App.CoreInfo.MasterExchangeName
	res.TxExpirationPeriod = action.App.CoreInfo.TxExpirationPeriod

	hal.Render(action.W, res)
}
