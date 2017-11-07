package horizon

import (
	"bullioncoin.githost.io/development/go/xdr"
	"bullioncoin.githost.io/development/horizon/db2/core"
	"bullioncoin.githost.io/development/horizon/render/hal"
	"bullioncoin.githost.io/development/horizon/render/problem"
	"bullioncoin.githost.io/development/horizon/resource"
	"database/sql"
)

type AccountTypeLimitsAllAction struct {
	Action

	AccountTypeLimits map[int32]*resource.Limits
}

func (action *AccountTypeLimitsAllAction) JSON() {
	action.Do(
		action.prepareList,
		action.loadData,
		func() {
			hal.Render(action.W, action.AccountTypeLimits)
		},
	)
}

func (action *AccountTypeLimitsAllAction) prepareList() {
	action.AccountTypeLimits = make(map[int32]*resource.Limits)
	defaultLimits := core.DefaultLimits()
	for _, accountType := range xdr.AccountTypeAll {
		var limits resource.Limits
		limits.Populate(defaultLimits)
		action.AccountTypeLimits[int32(accountType)] = &limits
	}
}

func (action *AccountTypeLimitsAllAction) loadData() {
	actualAccountTypeLimits := []core.AccountTypeLimits{}
	err := action.CoreQ().AccountTypeLimits().Select(&actualAccountTypeLimits)
	if err != nil {
		if err != sql.ErrNoRows {
			action.Err = &problem.ServerError
			action.Log.WithStack(err).WithError(err).Error("Could not get default limits from the database")
			return
		}

		err = nil
	}

	for _, actualEntry := range actualAccountTypeLimits {
		action.AccountTypeLimits[actualEntry.AccountType].Populate(actualEntry.Limits)
	}
}
