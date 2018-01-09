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
	ticker := time.NewTicker(h.points.BucketDuration())
	for {
		<-ticker.C
		cut := h.points.Shift()

		if h.preceded.Timestamp.After(cut.Timestamp) {
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

// Render fills missing buckets with previously known values
// Guaranteed to return non-nil values
func (h *Histogram) Render() []Point {
	var zero int64 = 0
	points := make([]Point, 0, len(h.points))
	for idx, point := range h.points {
		value := point.Value
		if value == nil {
			if idx == 0 {
				if h.preceded != nil {
					value = h.preceded.Value
				} else {
					value = &zero
				}
			} else {
				value = points[idx-1].Value
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
