package utils

import (
	"math"
	"testing"

	"github.com/magiconair/properties/assert"
	"github.com/stretchr/testify/require"
	"gitlab.com/swarmfund/horizon/db2/core"
)

var Tests = []struct {
	coreFees []core.FeeEntry
	zeroFee  core.FeeEntry
	expected []core.FeeEntry
}{
	{nil, core.FeeEntry{}, nil},
	// gap at the begging, mid, end
	{[]core.FeeEntry{
		{LowerBound: 402, UpperBound: 500, Percent: 4},
		{LowerBound: 100, UpperBound: 100, Percent: 2},
		{LowerBound: 101, UpperBound: 200, Percent: 3},
		{LowerBound: 300, UpperBound: 400, Percent: 4},
	},
		core.FeeEntry{},
		[]core.FeeEntry{
			{LowerBound: 0, UpperBound: 99},
			{LowerBound: 100, UpperBound: 100, Percent: 2},
			{LowerBound: 101, UpperBound: 200, Percent: 3},
			{LowerBound: 201, UpperBound: 299, Percent: 0},
			{LowerBound: 300, UpperBound: 400, Percent: 4},
			{LowerBound: 401, UpperBound: 401, Percent: 0},
			{LowerBound: 402, UpperBound: 500, Percent: 4},
			{LowerBound: 501, UpperBound: math.MaxInt64, Percent: 0},
		}},
}

func TestParse(t *testing.T) {
	for _, v := range Tests {
		actualResult := FillFeeGaps(v.coreFees, v.zeroFee)
		require.Equal(t, len(v.expected), len(actualResult))
		for i := range v.expected {
			require.Equal(t, v.expected[i], actualResult[i])
		}
	}
}

func TestSmartFillFeeGaps(t *testing.T) {
	t.Run("one gap", func(t *testing.T) {
		primaryFees := []core.FeeEntry{
			{
				LowerBound: 0,
				UpperBound: 5,
			},
			{
				LowerBound: 10,
				UpperBound: 15,
			},
		}
		secondaryFees := []core.FeeEntry{
			{
				LowerBound: 2,
				UpperBound: 20,
				Percent:    2,
			},
		}

		expected := []core.FeeEntry{
			{
				LowerBound: 0,
				UpperBound: 5,
			},
			{
				LowerBound: 6,
				UpperBound: 9,
				Percent:    2,
			},
			{
				LowerBound: 10,
				UpperBound: 15,
			},
			{
				LowerBound: 16,
				UpperBound: 20,
				Percent:    2,
			},
		}

		got := SmartFillFeeGaps(primaryFees, secondaryFees)

		assert.Equal(t, got, expected)
	})
	t.Run("two gaps", func(t *testing.T) {
		primaryFees := []core.FeeEntry{
			{
				LowerBound: 0,
				UpperBound: 5,
			},
			{
				LowerBound: 10,
				UpperBound: 15,
			},
			{
				LowerBound: 20,
				UpperBound: 22,
			},
		}
		secondaryFees := []core.FeeEntry{
			{
				LowerBound: 2,
				UpperBound: 20,
				Percent:    2,
			},
			{
				LowerBound: 21,
				UpperBound: 25,
				Percent:    3,
			},
		}

		expected := []core.FeeEntry{
			{
				LowerBound: 0,
				UpperBound: 5,
			},
			{
				LowerBound: 6,
				UpperBound: 9,
				Percent:    2,
			},
			{
				LowerBound: 10,
				UpperBound: 15,
			},
			{
				LowerBound: 16,
				UpperBound: 19,
				Percent:    2,
			},
			{
				LowerBound: 20,
				UpperBound: 22,
			},
			{
				LowerBound: 23,
				UpperBound: 25,
				Percent:    3,
			},
		}

		got := SmartFillFeeGaps(primaryFees, secondaryFees)

		assert.Equal(t, got, expected)
	})
	t.Run("no primary fees", func(t *testing.T) {
		primaryFees := []core.FeeEntry{}
		secondaryFees := []core.FeeEntry{
			{
				LowerBound: 2,
				UpperBound: 20,
				Percent:    2,
			},
		}

		expected := []core.FeeEntry{
			{
				LowerBound: 2,
				UpperBound: 20,
				Percent:    2,
			},
		}

		got := SmartFillFeeGaps(primaryFees, secondaryFees)

		assert.Equal(t, got, expected)
	})
	t.Run("no secondary fees", func(t *testing.T) {
		primaryFees := []core.FeeEntry{
			{
				LowerBound: 0,
				UpperBound: 5,
			},
			{
				LowerBound: 10,
				UpperBound: 15,
			},
		}
		expected := []core.FeeEntry{
			{
				LowerBound: 0,
				UpperBound: 5,
			},

			{
				LowerBound: 10,
				UpperBound: 15,
			},
		}

		got := SmartFillFeeGaps(primaryFees, nil)

		assert.Equal(t, got, expected)
	})
	t.Run("no overlap", func(t *testing.T) {
		primaryFees := []core.FeeEntry{
			{
				LowerBound: 0,
				UpperBound: 5,
			},
			{
				LowerBound: 10,
				UpperBound: 15,
			},
		}
		secondaryFees := []core.FeeEntry{
			{
				LowerBound: 0,
				UpperBound: 5,
				Percent:    2,
			},
		}

		expected := []core.FeeEntry{
			{
				LowerBound: 0,
				UpperBound: 5,
			},
			{
				LowerBound: 10,
				UpperBound: 15,
			},
		}

		got := SmartFillFeeGaps(primaryFees, secondaryFees)

		assert.Equal(t, got, expected)
	})
}
