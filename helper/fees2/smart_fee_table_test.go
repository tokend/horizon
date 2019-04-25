package fees2

import (
	"testing"

	"gitlab.com/tokend/horizon/helper"

	"math"

	"github.com/stretchr/testify/assert"
	"gitlab.com/tokend/go/amount"
	"gitlab.com/tokend/go/xdr"
	core "gitlab.com/tokend/horizon/db2/core2"
)

func TestNewSmartFeeTable(t *testing.T) {
	t.Run("valid", func(t *testing.T) {
		fees := []core.Fee{
			{
				FeeType:    int32(xdr.FeeTypePaymentFee),
				Subtype:    int64(xdr.PaymentFeeTypeOutgoing),
				Asset:      "USD",
				LowerBound: 0,
				UpperBound: 5,
			},
			{
				FeeType:    int32(xdr.FeeTypePaymentFee),
				Subtype:    int64(xdr.PaymentFeeTypeOutgoing),
				Asset:      "USD",
				LowerBound: 10,
				UpperBound: 15,
			},
		}
		expected := helper.SmartFeeTable{
			helper.FeeGroup{
				AssetCode: "USD",
				FeeType:   int32(xdr.FeeTypePaymentFee),
				Subtype:   int64(xdr.PaymentFeeTypeOutgoing),
			}: []helper.FeeWrapper{
				{
					Fee: core.Fee{
						FeeType:    int32(xdr.FeeTypePaymentFee),
						Subtype:    int64(xdr.PaymentFeeTypeOutgoing),
						Asset:      "USD",
						LowerBound: 0,
						UpperBound: 5,
					},
				},
				{
					Fee: core.Fee{
						FeeType:    int32(xdr.FeeTypePaymentFee),
						Subtype:    int64(xdr.PaymentFeeTypeOutgoing),
						Asset:      "USD",
						LowerBound: 10,
						UpperBound: 15,
					},
				},
			},
		}
		sft := helper.NewSmartFeeTable(fees)
		assert.Equal(t, sft, expected)
	})
}

func TestSmartFeeTable_GetValuesByAsset(t *testing.T) {
	fees := []core.Fee{
		{
			FeeType:    int32(xdr.FeeTypePaymentFee),
			Subtype:    int64(xdr.PaymentFeeTypeOutgoing),
			Asset:      "USD",
			LowerBound: 0,
			UpperBound: 5,
		},
		{
			FeeType:    int32(xdr.FeeTypePaymentFee),
			Subtype:    int64(xdr.PaymentFeeTypeOutgoing),
			Asset:      "USD",
			LowerBound: 10,
			UpperBound: 15,
		},
	}
	expected := map[string][]helper.FeeWrapper{
		"USD": {
			{
				Fee: core.Fee{
					FeeType:    int32(xdr.FeeTypePaymentFee),
					Subtype:    int64(xdr.PaymentFeeTypeOutgoing),
					Asset:      "USD",
					LowerBound: 0,
					UpperBound: 5,
				},
			},
			{
				Fee: core.Fee{
					FeeType:    int32(xdr.FeeTypePaymentFee),
					Subtype:    int64(xdr.PaymentFeeTypeOutgoing),
					Asset:      "USD",
					LowerBound: 10,
					UpperBound: 15,
				},
			},
		},
	}
	sft := helper.NewSmartFeeTable(fees)
	got := sft.GetValuesByAsset()
	assert.Equal(t, got, expected)

}

func TestSmartFeeTable_Update(t *testing.T) {
	t.Run("valid", func(t *testing.T) {
		fees := []core.Fee{
			{
				FeeType:    int32(xdr.FeeTypePaymentFee),
				Subtype:    int64(xdr.PaymentFeeTypeOutgoing),
				Asset:      "USD",
				LowerBound: 0,
				UpperBound: 5,
			},
			{
				FeeType:    int32(xdr.FeeTypePaymentFee),
				Subtype:    int64(xdr.PaymentFeeTypeOutgoing),
				Asset:      "USD",
				LowerBound: 10,
				UpperBound: 15,
			},
		}
		secondaryFees := []core.Fee{
			{
				FeeType:    int32(xdr.FeeTypePaymentFee),
				Subtype:    int64(xdr.PaymentFeeTypeOutgoing),
				Asset:      "USD",
				LowerBound: 2,
				UpperBound: 20,
				Percent:    2,
			},
			{
				FeeType:    int32(xdr.FeeTypePaymentFee),
				Subtype:    int64(xdr.PaymentFeeTypeOutgoing),
				Asset:      "USD",
				LowerBound: 21,
				UpperBound: 25,
				Percent:    3,
			},
		}
		expectedFeeTable := helper.SmartFeeTable{
			helper.FeeGroup{
				AssetCode: "USD",
				FeeType:   int32(xdr.FeeTypePaymentFee),
				Subtype:   int64(xdr.PaymentFeeTypeOutgoing),
			}: []helper.FeeWrapper{
				{
					Fee: core.Fee{
						FeeType:    int32(xdr.FeeTypePaymentFee),
						Subtype:    int64(xdr.PaymentFeeTypeOutgoing),
						Asset:      "USD",
						LowerBound: 0,
						UpperBound: 5,
					},
				},
				{
					Fee: core.Fee{
						FeeType:    int32(xdr.FeeTypePaymentFee),
						Subtype:    int64(xdr.PaymentFeeTypeOutgoing),
						Asset:      "USD",
						LowerBound: 6,
						UpperBound: 9,
						Percent:    2,
					},
				},
				{
					Fee: core.Fee{
						FeeType:    int32(xdr.FeeTypePaymentFee),
						Subtype:    int64(xdr.PaymentFeeTypeOutgoing),
						Asset:      "USD",
						LowerBound: 10,
						UpperBound: 15,
					},
				},
				{
					Fee: core.Fee{
						FeeType:    int32(xdr.FeeTypePaymentFee),
						Subtype:    int64(xdr.PaymentFeeTypeOutgoing),
						Asset:      "USD",
						LowerBound: 16,
						UpperBound: 20,
						Percent:    2,
					},
				},
				{
					Fee: core.Fee{
						FeeType:    int32(xdr.FeeTypePaymentFee),
						Subtype:    int64(xdr.PaymentFeeTypeOutgoing),
						Asset:      "USD",
						LowerBound: 21,
						UpperBound: 25,
						Percent:    3,
					},
				},
			},
		}

		sft := helper.NewSmartFeeTable(fees)
		sft.Update(secondaryFees)
		assert.Equal(t, sft, expectedFeeTable)
	})
	t.Run("no primary fees", func(t *testing.T) {
		var primaryFees []core.Fee
		secondaryFees := []core.Fee{
			{
				FeeType:    int32(xdr.FeeTypePaymentFee),
				Subtype:    int64(xdr.PaymentFeeTypeOutgoing),
				Asset:      "USD",
				LowerBound: 2,
				UpperBound: 20,
				Percent:    2,
			},
			{
				FeeType:    int32(xdr.FeeTypePaymentFee),
				Subtype:    int64(xdr.PaymentFeeTypeOutgoing),
				Asset:      "USD",
				LowerBound: 21,
				UpperBound: 25,
				Percent:    3,
			},
		}
		expectedFeeTable := helper.SmartFeeTable{
			helper.FeeGroup{
				AssetCode: "USD",
				FeeType:   int32(xdr.FeeTypePaymentFee),
				Subtype:   int64(xdr.PaymentFeeTypeOutgoing),
			}: []helper.FeeWrapper{
				{
					Fee: core.Fee{
						FeeType:    int32(xdr.FeeTypePaymentFee),
						Subtype:    int64(xdr.PaymentFeeTypeOutgoing),
						Asset:      "USD",
						LowerBound: 2,
						UpperBound: 20,
						Percent:    2,
					},
				},
				{
					Fee: core.Fee{
						FeeType:    int32(xdr.FeeTypePaymentFee),
						Subtype:    int64(xdr.PaymentFeeTypeOutgoing),
						Asset:      "USD",
						LowerBound: 21,
						UpperBound: 25,
						Percent:    3,
					},
				},
			},
		}

		sft := helper.NewSmartFeeTable(primaryFees)
		sft.Update(secondaryFees)
		assert.Equal(t, sft, expectedFeeTable)
	})
	t.Run("different attributes of secondary fees", func(t *testing.T) {
		fees := []core.Fee{
			{
				FeeType:    int32(xdr.FeeTypePaymentFee),
				Subtype:    int64(xdr.PaymentFeeTypeOutgoing),
				Asset:      "USD",
				LowerBound: 0,
				UpperBound: 5,
			},
			{
				FeeType:    int32(xdr.FeeTypePaymentFee),
				Subtype:    int64(xdr.PaymentFeeTypeOutgoing),
				Asset:      "USD",
				LowerBound: 10,
				UpperBound: 15,
			},
		}
		secondaryFees := []core.Fee{
			{
				FeeType:    int32(xdr.FeeTypeOfferFee),
				Subtype:    int64(0),
				Asset:      "USD",
				LowerBound: 2,
				UpperBound: 20,
				Percent:    2,
			},
			{
				FeeType:    int32(xdr.FeeTypeOfferFee),
				Subtype:    int64(0),
				Asset:      "USD",
				LowerBound: 21,
				UpperBound: 25,
				Percent:    3,
			},
		}
		expectedFeeTable := helper.SmartFeeTable{
			helper.FeeGroup{
				AssetCode: "USD",
				FeeType:   int32(xdr.FeeTypePaymentFee),
				Subtype:   int64(xdr.PaymentFeeTypeOutgoing),
			}: []helper.FeeWrapper{
				{
					Fee: core.Fee{
						FeeType:    int32(xdr.FeeTypePaymentFee),
						Subtype:    int64(xdr.PaymentFeeTypeOutgoing),
						Asset:      "USD",
						LowerBound: 0,
						UpperBound: 5,
					},
				},
				{
					Fee: core.Fee{
						FeeType:    int32(xdr.FeeTypePaymentFee),
						Subtype:    int64(xdr.PaymentFeeTypeOutgoing),
						Asset:      "USD",
						LowerBound: 10,
						UpperBound: 15,
					},
				},
			},

			helper.FeeGroup{
				AssetCode: "USD",
				FeeType:   int32(xdr.FeeTypeOfferFee),
				Subtype:   int64(0),
			}: []helper.FeeWrapper{
				{
					Fee: core.Fee{
						FeeType:    int32(xdr.FeeTypeOfferFee),
						Subtype:    int64(0),
						Asset:      "USD",
						LowerBound: 2,
						UpperBound: 20,
						Percent:    2,
					},
				},
				{
					Fee: core.Fee{
						FeeType:    int32(xdr.FeeTypeOfferFee),
						Subtype:    int64(0),
						Asset:      "USD",
						LowerBound: 21,
						UpperBound: 25,
						Percent:    3,
					},
				},
			},
		}

		sft := helper.NewSmartFeeTable(fees)
		sft.Update(secondaryFees)
		assert.Equal(t, sft, expectedFeeTable)
	})
	t.Run("both nil", func(t *testing.T) {
		var primaryFees []core.Fee
		var secondaryFees []core.Fee
		sft := helper.NewSmartFeeTable(primaryFees)
		sft.Update(secondaryFees)
		assert.Empty(t, sft)
	})

	t.Run("valid", func(t *testing.T) {
		fees := []core.Fee{
			{
				FeeType:    int32(xdr.FeeTypePaymentFee),
				Subtype:    int64(xdr.PaymentFeeTypeOutgoing),
				Asset:      "USD",
				LowerBound: 0,
				UpperBound: 815 * amount.One,
			},
			{
				FeeType:    int32(xdr.FeeTypePaymentFee),
				Subtype:    int64(xdr.PaymentFeeTypeOutgoing),
				Asset:      "USD",
				LowerBound: 816 * amount.One,
				UpperBound: 1000 * amount.One,
			},
		}
		secondaryFees := []core.Fee{}
		expectedFeeGroup := helper.FeeGroup{
			AssetCode: "USD",
			FeeType:   int32(xdr.FeeTypePaymentFee),
			Subtype:   int64(xdr.PaymentFeeTypeOutgoing),
		}
		expectedFeeTable := helper.SmartFeeTable{
			helper.FeeGroup{
				AssetCode: "USD",
				FeeType:   int32(xdr.FeeTypePaymentFee),
				Subtype:   int64(xdr.PaymentFeeTypeOutgoing),
			}: []helper.FeeWrapper{
				{
					Fee: core.Fee{
						FeeType:    int32(xdr.FeeTypePaymentFee),
						Subtype:    int64(xdr.PaymentFeeTypeOutgoing),
						Asset:      "USD",
						LowerBound: 0,
						UpperBound: 815 * amount.One,
					},
				},
				{
					Fee: core.Fee{
						FeeType:    int32(xdr.FeeTypePaymentFee),
						Subtype:    int64(xdr.PaymentFeeTypeOutgoing),
						Asset:      "USD",
						LowerBound: 815*amount.One + 1,
						UpperBound: 816*amount.One - 1,
					},
					NotExists: true,
				},
				{
					Fee: core.Fee{
						FeeType:    int32(xdr.FeeTypePaymentFee),
						Subtype:    int64(xdr.PaymentFeeTypeOutgoing),
						Asset:      "USD",
						LowerBound: 816 * amount.One,
						UpperBound: 1000 * amount.One,
					},
				},
				{
					Fee: core.Fee{
						FeeType:    int32(xdr.FeeTypePaymentFee),
						Subtype:    int64(xdr.PaymentFeeTypeOutgoing),
						Asset:      "USD",
						LowerBound: 1000*amount.One + 1,
						UpperBound: math.MaxInt64,
					},
					NotExists: true,
				},
			},
		}

		sft := helper.NewSmartFeeTable(fees)
		sft.Update(secondaryFees)
		sft.AddZeroFees([]string{"USD"})

		assert.Equal(t, sft[expectedFeeGroup], expectedFeeTable[expectedFeeGroup])
	})
}
