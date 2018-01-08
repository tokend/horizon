package charts

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestNewPointsPanics(t *testing.T) {
	cases := []struct {
		name   string
		count  uint
		bucket time.Duration
	}{
		{"invalid count", 1, 10},
		{"negative duration", 2, -10},
		{"zero duration", 2, 0},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			for i := uint(0); i < 2; i++ {
				assert.Panics(t, func() {
					NewPoints(tc.count, tc.bucket, time.Now())
				})
			}
		})
	}
}

func TestNewPoints(t *testing.T) {
	base := time.Unix(1337, 0)
	points := NewPoints(3, 1*time.Second, base)
	assert.Len(t, points, 3)
	assert.True(t, points[2].Timestamp.After(base))
	assert.Equal(t, points[2].Timestamp, base.Add(1*time.Second))
	assert.Equal(t, points[1].Timestamp, base)
}

//insert in not nil element
func TestPoints_Insert(t *testing.T) {

	inputVal := []int64{1, 1, 1, 1, 2, 1}

	durations := []time.Duration{
		5 * time.Minute, 15 * time.Minute, 25 * time.Minute,
		35 * time.Minute, 45 * time.Minute, 55 * time.Minute,
	}

	//p = 1, 1, 0, 2, 1
	p := make(Points, 6)
	for i, d := range durations {
		p[i].Timestamp = time.Now().Add(-d)
		p[i].Value = &inputVal[i]
	}

	p.Insert(2, 11)
	assert.EqualValues(t, 11, *p[2].Value)
}
