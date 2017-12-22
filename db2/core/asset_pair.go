package core

import (
	"github.com/go-errors/errors"
	"gitlab.com/swarmfund/go/amount"
)

type AssetPair struct {
	BaseAsset               string `db:"base"`
	QuoteAsset              string `db:"quote"`
	CurrentPrice            int64  `db:"current_price"`
	PhysicalPrice           int64  `db:"physical_price"`
	PhysicalPriceCorrection int64  `db:"physical_price_correction"`
	MaxPriceStep            int64  `db:"max_price_step"`
	Policies                int32  `db:"policies"`
}

// ConvertToDestAsset - converts specified amount to dest asset using current price,
// returns false - if failed
func (pair AssetPair) ConvertToDestAsset(destCode string, amountToConvert int64) (int64, bool, error) {
	if pair.CurrentPrice == 0 {
		return 0, false, errors.New("Price is invalid")
	}

	if pair.QuoteAsset == destCode {
		result, isOverflow := amount.BigDivide(amountToConvert, pair.CurrentPrice, amount.One, amount.ROUND_DOWN)
		return result, !isOverflow, nil
	}

	result, isOverflow := amount.BigDivide(amountToConvert, amount.One, pair.CurrentPrice, amount.ROUND_DOWN)
	return result, !isOverflow, nil
}