package charts

import (
	"time"
)

type Histogram struct {
	duration time.Duration
	points   Points
	preceded Point
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
	for ; ; <-ticker.C {
		h.points.Shift()
	}
}

func (h *Histogram) Run(value int64, ts time.Time) {
	offset := h.points.Last().Timestamp.Sub(ts)
	idx := len(h.points) - int(offset/h.points.BucketDuration())
	if idx >= 0 && idx < len(h.points) {
		// point fits into interval
		h.points.Insert(idx, value)
	}
	if idx < 0 {
		// storing latest value before first interval value
		if ts.After(h.preceded.Timestamp) {
			h.preceded = Point{ts, value}
		}
	}
}

func (h *Histogram) Render() []Point {
	points := make([]Point, 0, len(h.points))
	for idx, point := range h.points {
		value := point.Value
		if value == 0 {
			if idx == 0 {
				value = h.preceded.Value
			} else {
				value = points[idx-1].Value
			}
		}
		points = append(points, Point{
			Timestamp: point.Timestamp,
			Value:     value,
		})
	}
	points[len(points)-1].Timestamp = time.Now().UTC()

	return points
}
