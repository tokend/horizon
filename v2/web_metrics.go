package v2

import (
	"github.com/rcrowley/go-metrics"
	"time"
)

type WebMetrics struct {
	requestCounter metrics.Counter
	failureCounter metrics.Counter

	requestTimer metrics.Timer

	failureMeter metrics.Meter
	successMeter metrics.Meter
}

func NewWebMetrics() *WebMetrics {
	return &WebMetrics{
		requestCounter: metrics.NewCounter(),
		failureCounter: metrics.NewCounter(),

		requestTimer: metrics.NewTimer(),
		failureMeter: metrics.NewMeter(),
		successMeter: metrics.NewMeter(),
	}
}

func (m *WebMetrics) Update(requestDuration time.Duration, responseStatus int) {
	m.requestTimer.Update(requestDuration)

	if responseStatus >= 200 && responseStatus < 400 {
		m.successMeter.Mark(1)
	} else if responseStatus >= 400 && responseStatus < 600 {
		m.failureCounter.Inc(1)
		m.failureMeter.Mark(1)
	}
}

func (m *WebMetrics) ResetTimer() {
	m.requestTimer = metrics.NewTimer()
}
