package charts

import (
	"time"
)

type Histogram struct {
	duration time.Duration
	points   Points
	preceded *Point
}

func NewHistogram(duration time.Duration, count uint) *Histogram {
	h := Histogram{
		duration: duration,
		points:   NewPoints(count, duration/time.Duration(count), time.Now().UTC()),
	}

	go h.Ticker()

	return &h
}

func (h *Histogram) Ticker() {
	// intentionally no defer, no way to recover with proper state
	ticker := time.NewTicker(h.points.BucketDuration())
	for {
		<-ticker.C
		cut := h.points.Shift()

		if cut.Value != nil && (h.preceded == nil || cut.Timestamp.After(h.preceded.Timestamp)) {
			h.preceded = &cut
		}
	}
}

func (h Histogram) getIndex(ts time.Time) int {
	offset := h.points.Last().Timestamp.Sub(ts)
	idx := len(h.points) - int(offset/h.points.BucketDuration()) - 1
	return idx
}

func (h *Histogram) Run(value int64, ts time.Time) {
	idx := h.getIndex(ts)
	if idx >= 0 && idx < len(h.points) {
		// point fits into interval
		h.points.Insert(idx, value)
	}
	if idx < 0 {
		// storing latest value before first interval value
		if h.preceded == nil || ts.After(h.preceded.Timestamp) {
			h.preceded = &Point{ts, &value}
		}
	}
}

func (h *Histogram) Add(value int64, ts time.Time) {
	idx := h.getIndex(ts)
	if idx >= 0 && idx < len(h.points) {
		// point fits into interval
		h.points.Add(idx, value)
	}
	if idx < 0 {
		// storing latest value before first interval value
		if h.preceded == nil || ts.After(h.preceded.Timestamp) {
			h.preceded = &Point{ts, &value}
		}
	}
}

// Render fills missing buckets with previously known values
// Guaranteed to return non-nil values
func (h *Histogram) Render(fill bool) []Point {
	var zero int64 = 0
	points := make([]Point, 0, len(h.points))
	for idx, point := range h.points {
		value := point.Value
		if value == nil {
			if idx == 0 {
				if h.preceded != nil && fill {
					value = h.preceded.Value
				} else {
					value = &zero
				}
			} else {
				if fill {
					value = points[idx-1].Value
				} else {
					value = &zero
				}
			}
		}
		if value == nil {
			// marks issue with data provider layer
			panic("no initial value has been set")
		}
		points = append(points, Point{
			Timestamp: point.Timestamp,
			Value:     value,
		})
	}
	points[len(points)-1].Timestamp = time.Now().UTC()

	return points
}
