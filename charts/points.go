package charts

import (
	"time"
)

type Points []Point

// NewPoints returns populated Points object with `count` buckets each of `bucket` duration,
// treating `base` timestamp as current head.
func NewPoints(count uint, bucket time.Duration, base time.Time) Points {
	if count < 2 {
		panic("at least two buckets required for operation")
	}
	if bucket <= 0 {
		panic("bucket duration should be positive")
	}

	points := make([]Point, count)
	for i := uint(0); i < count; i++ {
		// adding 2 makes last bucket shifted to the "future"
		// to accommodate values inserted between shifts
		offset := time.Duration(i-count+2) * bucket
		points[i].Timestamp = base.Add(offset)
	}

	return points
}

func (p Points) BucketDuration() time.Duration {
	return p[1].Timestamp.Sub(p[0].Timestamp)
}

func (p Points) Last() Point {
	return p[len(p)-1]
}

func (p Points) Insert(idx int, value int64) {
	p[idx].Value = &value
}
