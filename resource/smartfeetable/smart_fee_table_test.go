package smartfeetable

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"gitlab.com/swarmfund/horizon/db2/core"
	"gitlab.com/tokend/go/xdr"
)

func TestNewSmartFeeTable(t *testing.T) {
	t.Run("valid", func(t *testing.T) {
		fees := []core.FeeEntry{
			{
				FeeType:    int(xdr.FeeTypePaymentFee),
				Subtype:    int64(xdr.PaymentFeeTypeOutgoing),
				Asset:      "USD",
				LowerBound: 0,
				UpperBound: 5,
			},
			{
				FeeType:    int(xdr.FeeTypePaymentFee),
				Subtype:    int64(xdr.PaymentFeeTypeOutgoing),
				Asset:      "USD",
				LowerBound: 10,
				UpperBound: 15,
			},
		}
		expected := SmartFeeTable{
			FeeGroup{
				AssetCode: "USD",
				FeeType:   int(xdr.FeeTypePaymentFee),
				Subtype:   int64(xdr.PaymentFeeTypeOutgoing),
			}: []core.FeeEntry{
				{
					FeeType:    int(xdr.FeeTypePaymentFee),
					Subtype:    int64(xdr.PaymentFeeTypeOutgoing),
					Asset:      "USD",
					LowerBound: 0,
					UpperBound: 5,
				},
				{
					FeeType:    int(xdr.FeeTypePaymentFee),
					Subtype:    int64(xdr.PaymentFeeTypeOutgoing),
					Asset:      "USD",
					LowerBound: 10,
					UpperBound: 15,
				},
			},
		}
		sft := NewSmartFeeTable(fees)
		assert.Equal(t, sft, expected)
	})
}

func TestSmartFeeTable_GetValuesByAsset(t *testing.T) {
	fees := []core.FeeEntry{
		{
			FeeType:    int(xdr.FeeTypePaymentFee),
			Subtype:    int64(xdr.PaymentFeeTypeOutgoing),
			Asset:      "USD",
			LowerBound: 0,
			UpperBound: 5,
		},
		{
			FeeType:    int(xdr.FeeTypePaymentFee),
			Subtype:    int64(xdr.PaymentFeeTypeOutgoing),
			Asset:      "USD",
			LowerBound: 10,
			UpperBound: 15,
		},
	}
	expected := map[string][]core.FeeEntry{
		"USD": {
			{
				FeeType:    int(xdr.FeeTypePaymentFee),
				Subtype:    int64(xdr.PaymentFeeTypeOutgoing),
				Asset:      "USD",
				LowerBound: 0,
				UpperBound: 5,
			},
			{
				FeeType:    int(xdr.FeeTypePaymentFee),
				Subtype:    int64(xdr.PaymentFeeTypeOutgoing),
				Asset:      "USD",
				LowerBound: 10,
				UpperBound: 15,
			},
		},
	}
	sft := NewSmartFeeTable(fees)
	got := sft.GetValuesByAsset()
	assert.Equal(t, got, expected)

}

func TestSmartFeeTable_Update(t *testing.T) {
	t.Run("valid", func(t *testing.T) {
		fees := []core.FeeEntry{
			{
				FeeType:    int(xdr.FeeTypePaymentFee),
				Subtype:    int64(xdr.PaymentFeeTypeOutgoing),
				Asset:      "USD",
				LowerBound: 0,
				UpperBound: 5,
			},
			{
				FeeType:    int(xdr.FeeTypePaymentFee),
				Subtype:    int64(xdr.PaymentFeeTypeOutgoing),
				Asset:      "USD",
				LowerBound: 10,
				UpperBound: 15,
			},
		}
		secondaryFees := []core.FeeEntry{
			{
				FeeType:    int(xdr.FeeTypePaymentFee),
				Subtype:    int64(xdr.PaymentFeeTypeOutgoing),
				Asset:      "USD",
				LowerBound: 2,
				UpperBound: 20,
				Percent:    2,
			},
			{
				FeeType:    int(xdr.FeeTypePaymentFee),
				Subtype:    int64(xdr.PaymentFeeTypeOutgoing),
				Asset:      "USD",
				LowerBound: 21,
				UpperBound: 25,
				Percent:    3,
			},
		}
		expectedFeeTable := SmartFeeTable{
			FeeGroup{
				AssetCode: "USD",
				FeeType:   int(xdr.FeeTypePaymentFee),
				Subtype:   int64(xdr.PaymentFeeTypeOutgoing),
			}: []core.FeeEntry{
				{
					LowerBound: 0,
					UpperBound: 5,
					FeeType:    int(xdr.FeeTypePaymentFee),
					Subtype:    int64(xdr.PaymentFeeTypeOutgoing),
					Asset:      "USD",
				},
				{
					LowerBound: 6,
					UpperBound: 9,
					Percent:    2,
					FeeType:    int(xdr.FeeTypePaymentFee),
					Subtype:    int64(xdr.PaymentFeeTypeOutgoing),
					Asset:      "USD",
				},
				{
					LowerBound: 10,
					UpperBound: 15,
					FeeType:    int(xdr.FeeTypePaymentFee),
					Subtype:    int64(xdr.PaymentFeeTypeOutgoing),
					Asset:      "USD",
				},
				{
					LowerBound: 16,
					UpperBound: 20,
					Percent:    2,
					FeeType:    int(xdr.FeeTypePaymentFee),
					Subtype:    int64(xdr.PaymentFeeTypeOutgoing),
					Asset:      "USD",
				},
				{
					LowerBound: 21,
					UpperBound: 25,
					Percent:    3,
					FeeType:    int(xdr.FeeTypePaymentFee),
					Subtype:    int64(xdr.PaymentFeeTypeOutgoing),
					Asset:      "USD",
				},
			},
		}

		sft := NewSmartFeeTable(fees)
		sft.Update(secondaryFees)
		assert.Equal(t, sft, expectedFeeTable)
	})
	t.Run("no primary fees", func(t *testing.T) {
		var primaryFees []core.FeeEntry
		secondaryFees := []core.FeeEntry{
			{
				FeeType:    int(xdr.FeeTypePaymentFee),
				Subtype:    int64(xdr.PaymentFeeTypeOutgoing),
				Asset:      "USD",
				LowerBound: 2,
				UpperBound: 20,
				Percent:    2,
			},
			{
				FeeType:    int(xdr.FeeTypePaymentFee),
				Subtype:    int64(xdr.PaymentFeeTypeOutgoing),
				Asset:      "USD",
				LowerBound: 21,
				UpperBound: 25,
				Percent:    3,
			},
		}
		expectedFeeTable := SmartFeeTable{
			FeeGroup{
				AssetCode: "USD",
				FeeType:   int(xdr.FeeTypePaymentFee),
				Subtype:   int64(xdr.PaymentFeeTypeOutgoing),
			}: []core.FeeEntry{
				{
					FeeType:    int(xdr.FeeTypePaymentFee),
					Subtype:    int64(xdr.PaymentFeeTypeOutgoing),
					Asset:      "USD",
					LowerBound: 2,
					UpperBound: 20,
					Percent:    2,
				},
				{
					FeeType:    int(xdr.FeeTypePaymentFee),
					Subtype:    int64(xdr.PaymentFeeTypeOutgoing),
					Asset:      "USD",
					LowerBound: 21,
					UpperBound: 25,
					Percent:    3,
				},
			},
		}

		sft := NewSmartFeeTable(primaryFees)
		sft.Update(secondaryFees)
		assert.Equal(t, sft, expectedFeeTable)
	})
	t.Run("different attributes of secondary fees", func(t *testing.T) {
		fees := []core.FeeEntry{
			{
				FeeType:    int(xdr.FeeTypePaymentFee),
				Subtype:    int64(xdr.PaymentFeeTypeOutgoing),
				Asset:      "USD",
				LowerBound: 0,
				UpperBound: 5,
			},
			{
				FeeType:    int(xdr.FeeTypePaymentFee),
				Subtype:    int64(xdr.PaymentFeeTypeOutgoing),
				Asset:      "USD",
				LowerBound: 10,
				UpperBound: 15,
			},
		}
		secondaryFees := []core.FeeEntry{
			{
				FeeType:    int(xdr.FeeTypeOfferFee),
				Subtype:    int64(0),
				Asset:      "USD",
				LowerBound: 2,
				UpperBound: 20,
				Percent:    2,
			},
			{
				FeeType:    int(xdr.FeeTypeOfferFee),
				Subtype:    int64(0),
				Asset:      "USD",
				LowerBound: 21,
				UpperBound: 25,
				Percent:    3,
			},
		}
		expectedFeeTable := SmartFeeTable{
			FeeGroup{
				AssetCode: "USD",
				FeeType:   int(xdr.FeeTypePaymentFee),
				Subtype:   int64(xdr.PaymentFeeTypeOutgoing),
			}: []core.FeeEntry{
				{
					LowerBound: 0,
					UpperBound: 5,
					FeeType:    int(xdr.FeeTypePaymentFee),
					Subtype:    int64(xdr.PaymentFeeTypeOutgoing),
					Asset:      "USD",
				},
				{
					LowerBound: 10,
					UpperBound: 15,
					FeeType:    int(xdr.FeeTypePaymentFee),
					Subtype:    int64(xdr.PaymentFeeTypeOutgoing),
					Asset:      "USD",
				},
			},
			FeeGroup{
				AssetCode: "USD",
				FeeType:   int(xdr.FeeTypeOfferFee),
				Subtype:   int64(0),
			}: []core.FeeEntry{
				{
					LowerBound: 2,
					UpperBound: 20,
					Percent:    2,
					FeeType:    int(xdr.FeeTypeOfferFee),
					Subtype:    int64(0),
					Asset:      "USD",
				},
				{
					FeeType:    int(xdr.FeeTypeOfferFee),
					Subtype:    int64(0),
					Asset:      "USD",
					LowerBound: 21,
					UpperBound: 25,
					Percent:    3,
				},
			},
		}

		sft := NewSmartFeeTable(fees)
		sft.Update(secondaryFees)
		assert.Equal(t, sft, expectedFeeTable)
	})
	t.Run("both nil", func(t *testing.T) {
		var primaryFees []core.FeeEntry
		var secondaryFees []core.FeeEntry
		sft := NewSmartFeeTable(primaryFees)
		sft.Update(secondaryFees)
		assert.Empty(t, sft)
	})
}
