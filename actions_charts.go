package horizon

import (
	"strings"

	"gitlab.com/tokend/go/amount"
	"gitlab.com/tokend/horizon/charts"
	"gitlab.com/tokend/horizon/render/hal"
	"gitlab.com/tokend/horizon/render/problem"
	"gitlab.com/tokend/horizon/resource"
)

type ChartsAction struct {
	Action

	Code string

	Record   map[string]*charts.Histogram
	Resource resource.Charts
}

func (action *ChartsAction) JSON() {
	action.Do(
		action.loadParams,
		action.loadChart,
		action.renderResource,
		func() {
			hal.Render(action.W, action.Resource)
		},
	)
}

func (action *ChartsAction) loadParams() {
	action.Code = action.GetNonEmptyString("code")
}

func (action *ChartsAction) loadChart() {
	action.Record = action.App.charts.Get(action.Code)
	if action.Record == nil {
		action.Err = &problem.NotFound
		return
	}
}

func (action *ChartsAction) renderResource() {
	action.Resource = make(resource.Charts)
	var points []charts.Point
	for key, histogram := range action.Record {
		if strings.HasSuffix(action.Code, "volume") {
			points = histogram.Render(false)
		} else {
			points = histogram.Render(true)
		}
		for _, point := range points {
			action.Resource[key] = append(action.Resource[key], resource.Point{
				Timestamp: point.Timestamp,
				Value:     amount.String(*point.Value),
			})
		}
	}
}
