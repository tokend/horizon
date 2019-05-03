package handlers

import (
	"github.com/magiconair/properties/assert"
	"gitlab.com/tokend/regources/generated"
	"testing"
)

func TestSortConvertedStates(t *testing.T) {
	initialStates := []regources.ConvertedBalanceState{
		{
			Attributes: regources.ConvertedBalanceStateAttributes{
				ConvertedAmounts: regources.BalanceStateAttributeAmounts{
					Available: regources.Amount(0),
				},
				InitialAmounts: regources.BalanceStateAttributeAmounts{
					Available: regources.Amount(0),
				},
				IsConverted: false,
			},
		},
		{
			Attributes: regources.ConvertedBalanceStateAttributes{
				ConvertedAmounts: regources.BalanceStateAttributeAmounts{
					Available: regources.Amount(0),
				},
				InitialAmounts: regources.BalanceStateAttributeAmounts{
					Available: regources.Amount(200000),
				},
				IsConverted: false,
			},
		},
		{
			Attributes: regources.ConvertedBalanceStateAttributes{
				ConvertedAmounts: regources.BalanceStateAttributeAmounts{
					Available: regources.Amount(0),
				},
				InitialAmounts: regources.BalanceStateAttributeAmounts{
					Available: regources.Amount(3500000),
				},
				IsConverted: false,
			},
		},
		{
			Attributes: regources.ConvertedBalanceStateAttributes{
				ConvertedAmounts: regources.BalanceStateAttributeAmounts{
					Available: regources.Amount(4000000),
				},
				InitialAmounts: regources.BalanceStateAttributeAmounts{
					Available: regources.Amount(5000000),
				},
				IsConverted: true,
			},
		},
		{
			Attributes: regources.ConvertedBalanceStateAttributes{
				ConvertedAmounts: regources.BalanceStateAttributeAmounts{
					Available: regources.Amount(0),
				},
				InitialAmounts: regources.BalanceStateAttributeAmounts{
					Available: regources.Amount(4000000),
				},
				IsConverted: false,
			},
		},
		{
			Attributes: regources.ConvertedBalanceStateAttributes{
				ConvertedAmounts: regources.BalanceStateAttributeAmounts{
					Available: regources.Amount(2000000),
				},
				InitialAmounts: regources.BalanceStateAttributeAmounts{
					Available: regources.Amount(3000000),
				},
				IsConverted: true,
			},
		},
		{
			Attributes: regources.ConvertedBalanceStateAttributes{
				ConvertedAmounts: regources.BalanceStateAttributeAmounts{
					Available: regources.Amount(6000000),
				},
				InitialAmounts: regources.BalanceStateAttributeAmounts{
					Available: regources.Amount(7000000),
				},
				IsConverted: true,
			},
		},
		{
			Attributes: regources.ConvertedBalanceStateAttributes{
				ConvertedAmounts: regources.BalanceStateAttributeAmounts{
					Available: regources.Amount(0),
				},
				InitialAmounts: regources.BalanceStateAttributeAmounts{
					Available: regources.Amount(3200000),
				},
				IsConverted: false,
			},
		},
		{
			Attributes: regources.ConvertedBalanceStateAttributes{
				ConvertedAmounts: regources.BalanceStateAttributeAmounts{
					Available: regources.Amount(0),
				},
				InitialAmounts: regources.BalanceStateAttributeAmounts{
					Available: regources.Amount(5000000),
				},
				IsConverted: false,
			},
		},
		{
			Attributes: regources.ConvertedBalanceStateAttributes{
				ConvertedAmounts: regources.BalanceStateAttributeAmounts{
					Available: regources.Amount(5500000),
				},
				InitialAmounts: regources.BalanceStateAttributeAmounts{
					Available: regources.Amount(72000000),
				},
				IsConverted: true,
			},
		},
		{
			Attributes: regources.ConvertedBalanceStateAttributes{
				ConvertedAmounts: regources.BalanceStateAttributeAmounts{
					Available: regources.Amount(0),
				},
				InitialAmounts: regources.BalanceStateAttributeAmounts{
					Available: regources.Amount(0),
				},
				IsConverted: true,
			},
		},
	}

	expectedStates := []regources.ConvertedBalanceState{
		{
			Attributes: regources.ConvertedBalanceStateAttributes{
				ConvertedAmounts: regources.BalanceStateAttributeAmounts{
					Available: regources.Amount(6000000),
				},
				InitialAmounts: regources.BalanceStateAttributeAmounts{
					Available: regources.Amount(7000000),
				},
				IsConverted: true,
			},
		},
		{
			Attributes: regources.ConvertedBalanceStateAttributes{
				ConvertedAmounts: regources.BalanceStateAttributeAmounts{
					Available: regources.Amount(5500000),
				},
				InitialAmounts: regources.BalanceStateAttributeAmounts{
					Available: regources.Amount(72000000),
				},
				IsConverted: true,
			},
		},
		{
			Attributes: regources.ConvertedBalanceStateAttributes{
				ConvertedAmounts: regources.BalanceStateAttributeAmounts{
					Available: regources.Amount(4000000),
				},
				InitialAmounts: regources.BalanceStateAttributeAmounts{
					Available: regources.Amount(5000000),
				},
				IsConverted: true,
			},
		},
		{
			Attributes: regources.ConvertedBalanceStateAttributes{
				ConvertedAmounts: regources.BalanceStateAttributeAmounts{
					Available: regources.Amount(2000000),
				},
				InitialAmounts: regources.BalanceStateAttributeAmounts{
					Available: regources.Amount(3000000),
				},
				IsConverted: true,
			},
		},
		{
			Attributes: regources.ConvertedBalanceStateAttributes{
				ConvertedAmounts: regources.BalanceStateAttributeAmounts{
					Available: regources.Amount(0),
				},
				InitialAmounts: regources.BalanceStateAttributeAmounts{
					Available: regources.Amount(5000000),
				},
				IsConverted: false,
			},
		},
		{
			Attributes: regources.ConvertedBalanceStateAttributes{
				ConvertedAmounts: regources.BalanceStateAttributeAmounts{
					Available: regources.Amount(0),
				},
				InitialAmounts: regources.BalanceStateAttributeAmounts{
					Available: regources.Amount(4000000),
				},
				IsConverted: false,
			},
		},
		{
			Attributes: regources.ConvertedBalanceStateAttributes{
				ConvertedAmounts: regources.BalanceStateAttributeAmounts{
					Available: regources.Amount(0),
				},
				InitialAmounts: regources.BalanceStateAttributeAmounts{
					Available: regources.Amount(3500000),
				},
				IsConverted: false,
			},
		},
		{
			Attributes: regources.ConvertedBalanceStateAttributes{
				ConvertedAmounts: regources.BalanceStateAttributeAmounts{
					Available: regources.Amount(0),
				},
				InitialAmounts: regources.BalanceStateAttributeAmounts{
					Available: regources.Amount(3200000),
				},
				IsConverted: false,
			},
		},
		{
			Attributes: regources.ConvertedBalanceStateAttributes{
				ConvertedAmounts: regources.BalanceStateAttributeAmounts{
					Available: regources.Amount(0),
				},
				InitialAmounts: regources.BalanceStateAttributeAmounts{
					Available: regources.Amount(200000),
				},
				IsConverted: false,
			},
		},
		{
			Attributes: regources.ConvertedBalanceStateAttributes{
				ConvertedAmounts: regources.BalanceStateAttributeAmounts{
					Available: regources.Amount(0),
				},
				InitialAmounts: regources.BalanceStateAttributeAmounts{
					Available: regources.Amount(0),
				},
				IsConverted: true,
			},
		},
		{
			Attributes: regources.ConvertedBalanceStateAttributes{
				ConvertedAmounts: regources.BalanceStateAttributeAmounts{
					Available: regources.Amount(0),
				},
				InitialAmounts: regources.BalanceStateAttributeAmounts{
					Available: regources.Amount(0),
				},
				IsConverted: false,
			},
		},
	}

	t.Run("correctly sort the state", func(t *testing.T) {
		sortedStates := SortConvertedStates(initialStates)
		assert.Equal(t, sortedStates, expectedStates)
	})
}
