package charts

import (
	"time"
)

type Points []Point

// NewPoints returns populated Points object with `count` buckets each of `bucket` duration,
// treating `base` timestamp as current head.
// Last point will be in "future" to accommodate values inserted between shifts
func NewPoints(count uint, bucket time.Duration, base time.Time) Points {
	if count < 2 {
		panic("at least two buckets required for operation")
	}
	if bucket <= 0 {
		panic("bucket duration should be positive")
	}
	points := make([]Point, count)
	for i := uint(0); i < count-1; i++ {
		points[count-i].Timestamp = base.Add(-1 * time.Duration(i) * bucket)
	}
	points[count-1].Timestamp = base.Add(bucket)
	return points
}

func (p *Points) Shift() {
	// calculate bucket duration
	duration := (*p)[1].Timestamp.Sub((*p)[0].Timestamp)
	*p = append((*p)[1:], Point{
		Timestamp: (*p)[len(*p)-1].Timestamp.Add(duration),
	})
}

func (p Points) Last() Point {
	return p[len(p)-1]
}

func (p Points) Insert(idx int64, value int64) {
	if p[idx].Value != 0 {
		value = (p[idx].Value + value) / 2
	}
	p[idx].Value = value
}
