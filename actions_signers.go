package horizon

import (
	"bullioncoin.githost.io/development/horizon/db2/core"
	"bullioncoin.githost.io/development/horizon/render/hal"
	"bullioncoin.githost.io/development/horizon/render/problem"
	"bullioncoin.githost.io/development/horizon/resource"
)

type SignersIndexAction struct {
	Action
	Address     string
	CoreSigners []core.Signer
	Resource    resource.Signers
}

// JSON is a method for actions.JSON
func (action *SignersIndexAction) JSON() {
	action.Do(
		action.loadParams,
		action.loadRecord,
		func() {
			hal.Render(action.W, action.Resource)
		},
	)
}

func (action *SignersIndexAction) loadParams() {
	action.Address = action.GetString("id")
}

func (action *SignersIndexAction) loadRecord() {
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

	action.Resource.Populate(
		action.CoreSigners,
	)
}
