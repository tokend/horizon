package charts

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestPoints_Shift(t *testing.T) {
	points := Points{
		{time.Time{}.Add(1), 1},
		{time.Time{}.Add(2), 2},
		{time.Time{}.Add(3), 3},
	}
	l := len(points)
	points.Shift()
	assert.Len(t, points, l)
	assert.EqualValues(t, 0, points[l-1].Value)
	assert.InEpsilon(t, time.Now().UTC().Unix(), points[l-1].Timestamp.Unix(), 1)
}
