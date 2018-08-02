package utils

import (
	"testing"

	"github.com/magiconair/properties/assert"
	"gitlab.com/swarmfund/horizon/db2/core"
	"gitlab.com/tokend/go/xdr"
)

func TestSmartFeeTable_Populate(t *testing.T) {
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
			int(xdr.FeeTypePaymentFee): {
				int64(xdr.PaymentFeeTypeOutgoing): {
					"USD": []core.FeeEntry{
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
		int(xdr.FeeTypePaymentFee): {
			int64(xdr.PaymentFeeTypeOutgoing): {
				"USD": []core.FeeEntry{
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
			},
		},
	}

	sft := NewSmartFeeTable(fees)
	sft.Update(secondaryFees)
	assert.Equal(t, sft, expectedFeeTable)
}
