package horizon

import (
	"gitlab.com/swarmfund/go/amount"
	"gitlab.com/swarmfund/horizon/render/hal"
	"gitlab.com/swarmfund/horizon/resource"
)

type ChartsAction struct {
	Action

	Code     string
	Resource resource.Charts
}

func (action *ChartsAction) JSON() {
	action.Do(
		action.loadParams,
	)
}

func (action *ChartsAction) loadParams() {
	action.Code = action.GetNonEmptyString("code")
	histograms := action.App.charts.Get(action.Code)
	action.Resource = make(resource.Charts)
	for i, histogram := range histograms {
		points := histogram.Render()
		for _, point := range points {
			action.Resource[i] = append(action.Resource[i], resource.Point{
				Timestamp: point.Timestamp,
				Value:     amount.String(point.Value),
			})
		}
	}
	hal.Render(action.W, action.Resource)
}
