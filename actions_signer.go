package horizon

import (
	"gitlab.com/tokend/horizon/db2/core"
	"gitlab.com/tokend/horizon/render/hal"
	"gitlab.com/tokend/horizon/render/problem"
	"gitlab.com/tokend/horizon/resource"
)

type SignerShowAction struct {
	Action
	Address       string
	SignerAddress string
	CoreSigners   []core.Signer
	Resource      resource.Signer
}

// JSON is a method for actions.JSON
func (action *SignerShowAction) JSON() {
	action.Do(
		action.loadParams,
		action.loadRecord,
		func() {
			hal.Render(action.W, action.Resource)
		},
	)
}

func (action *SignerShowAction) loadParams() {
	action.Address = action.GetString("account_id")
	action.SignerAddress = action.GetString("id")
}

func (action *SignerShowAction) loadRecord() {
	account, err := action.CoreQ().Accounts().ByAddress(action.Address)
	if err != nil {
		action.Log.WithError(err).Error("Failed to load account by address")
		action.Err = &problem.ServerError
		return
	}

	if account == nil {
		action.Err = &problem.NotFound
		return
	}

	action.CoreSigners, action.Err = action.GetSigners(account)
	if action.Err != nil {
		return
	}

	var signer *core.Signer
	for i := range action.CoreSigners {
		if action.CoreSigners[i].Publickey == action.SignerAddress {
			signer = &action.CoreSigners[i]
			break
		}
	}

	if signer == nil {
		action.Err = &problem.NotFound
		return
	}

	action.Resource.Populate(*signer)
}
