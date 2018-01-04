package charts

import (
	"time"
)

type Histogram struct {
	Duration time.Duration
	Count    int64
	points   Points
}

func NewHistogram(duration time.Duration, count int64) *Histogram {
	h := Histogram{
		Duration: duration,
		Count:    count,
	}

	now := time.Now().UTC()
	h.points = make([]Point, h.Count)
	for i := time.Duration(count); i > 0; i-- {
		h.points[i-1].Timestamp = now.Add(-1 * (i - 1) * h.bucketLength())
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
		h.points.Insert(idx, value)
	}
}

func (h *Histogram) Render() []Point {
	points := make([]Point, 0, h.Count-1)
	for i := 1; i < len(h.points); i++ {
		value := int64(h.points[i].Value)
		if h.points[i].Value == 0 {
			value = int64(h.points[i-1].Value)
		}
		points = append(points, Point{
			Timestamp: h.points[i].Timestamp,
			Value:     value,
		})
	}
	return points
}
