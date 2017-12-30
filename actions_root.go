package horizon

import (
	"gitlab.com/swarmfund/horizon/ledger"
	"gitlab.com/swarmfund/horizon/render/hal"
	"gitlab.com/swarmfund/horizon/resource"
)

// RootAction provides a summary of the horizon instance and links to various
// useful endpoints
type RootAction struct {
	counter []int
	Action
}

// JSON renders the json response for RootAction
func (action *RootAction) JSON() {
	action.counter = append(action.counter, 2)
	action.App.UpdateStellarCoreInfo()

	var res resource.Root
	res.PopulateLedgerState(action.Ctx, ledger.CurrentState())
	res.NetworkPassphrase = action.App.CoreInfo.NetworkPassphrase
	res.CommissionAccountID = action.App.CoreInfo.CommissionAccountID
	res.MasterAccountID = action.App.CoreInfo.MasterAccountID
	res.OperationalAccountID = action.App.CoreInfo.OperationalAccountID
	res.MasterExchangeName = action.App.CoreInfo.MasterExchangeName
	res.TxExpirationPeriod = action.App.CoreInfo.TxExpirationPeriod

	action.Log.Error(action.counter)
	hal.Render(action.W, res)
}
