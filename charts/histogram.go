package charts

import (
	"time"
)

type Histogram struct {
	Duration time.Duration
	Count    int64

	points   Points
	preceded Point
}

func NewHistogram(duration time.Duration, count int64) *Histogram {
	h := Histogram{
		Duration: duration,
		Count:    count,
	}

	now := time.Now().UTC()
	h.points = make([]Point, h.Count)
	for i := int64(1); i <= count; i++ {
		h.points[count-i].Timestamp = now.Add(-1 * time.Duration(i) * h.bucketLength())
	}

	go h.Ticker()

	return &h
}

func (h *Histogram) bucketLength() time.Duration {
	return h.Duration / time.Duration(h.Count)
}

func (h *Histogram) Ticker() {
	ticker := time.NewTicker(h.bucketLength())
	for ; ; <-ticker.C {
		h.points.Shift()
	}
}

func (h *Histogram) Run(value int64, ts time.Time) {
	idx := h.Count - int64(h.points.Last().Timestamp.Sub(ts)/(h.Duration/time.Duration(h.Count)))
	if idx >= 0 && idx < h.Count {
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
	points := make([]Point, 0, h.Count)
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

	return points
}
