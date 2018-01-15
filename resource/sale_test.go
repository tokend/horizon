package resource

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"gitlab.com/swarmfund/go/amount"
	"gitlab.com/swarmfund/horizon/db2/core"
)

func TestSale_PopulateStat(t *testing.T) {
	sale := Sale{
		BaseAsset:  "QWE",
		QuoteAsset: "SUN",
	}
	assetPair := &core.AssetPair{
		BaseAsset:    "QWE",
		QuoteAsset:   "SUN",
		CurrentPrice: 2 * amount.One,
	}

	// test user 1 -> 2x 1 offer
	// test user 2 -> 1x 1 offer, 1x 2 balance
	// test user 3 -> 1x 1 offer, 1x 2 balance
	// test user 4 ->             1x 2 balance
	// 4 unique users,
	// 4 SUN at offers,
	//  3x2[QWE balance] x 2[price] = 16 SUN at balances
	offers := []core.Offer{
		{
			OwnerID: "test user 1",
			OrderBookEntry: core.OrderBookEntry{
				QuoteAmount: 1 * amount.One,
			},
		},
		{
			OwnerID: "test user 1",
			OrderBookEntry: core.OrderBookEntry{
				QuoteAmount: 1 * amount.One,
			},
		},
		{
			OwnerID: "test user 2",
			OrderBookEntry: core.OrderBookEntry{
				QuoteAmount: 1 * amount.One,
			},
		},
		{
			OwnerID: "test user 3",
			OrderBookEntry: core.OrderBookEntry{
				QuoteAmount: 1 * amount.One,
			},
		},
	}
	balances := []core.Balance{
		{
			AccountID: "test user 2",
			Amount:    2 * amount.One,
		},
		{
			AccountID: "test user 3",
			Amount:    2 * amount.One,
		},
		{
			AccountID: "test user 4",
			Amount:    2 * amount.One,
		},
	}

	err := sale.PopulateStat(offers, balances, assetPair)
	assert.Equal(t, nil, err)
	assert.Equal(t, 4, sale.Statistics.Investors)
	assert.Equal(t, amount.String(4*amount.One), sale.Statistics.AverageAmount)
}
