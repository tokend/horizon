package horizon

import (
	"gitlab.com/swarmfund/horizon/render/hal"
	"gitlab.com/swarmfund/horizon/render/problem"
	"gitlab.com/tokend/regources"
)

type StatisticsAction struct {
	Action

	Resource regources.SystemStatistics
}

func (action *StatisticsAction) JSON() {
	action.Do(
		// just to make sure core is alive
		action.App.UpdateStellarCoreInfo,
		action.loadExternalSystemPoolCount,
		func() {
			hal.Render(action.W, action.Resource)
		},
	)
}

func (action *StatisticsAction) loadExternalSystemPoolCount() {
	counts, err := action.CoreQ().ExternalSystemAccountIDPool().EntitiesCount()
	if err != nil {
		action.Err = &problem.ServerError
		action.Log.WithError(err).Error("failed to get external system pool counts")
		return
	}

	action.Resource.ExternalSystemPoolEntriesCount = counts
}
