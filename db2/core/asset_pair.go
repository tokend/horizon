package core

import (
	"gitlab.com/distributed_lab/logan/v3"
	"gitlab.com/distributed_lab/logan/v3/errors"
	"gitlab.com/tokend/go/amount"
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

//go:generate mockery -case underscore -name assetLoader -inpkg -testonly

// AssetLoader uses to load asset before converting, we need minimal amount
type assetLoader interface {
	LoadAsset(code string) (*Asset, error)
}

// ConvertToDestAsset - converts specified amount to dest asset using current price,
// returns false - if failed
func (pair AssetPair) ConvertToDestAsset(destCode string, amountToConvert int64, loader assetLoader,
) (int64, bool, error) {
	if pair.CurrentPrice == 0 {
		return 0, false, errors.New("Price is invalid")
	}

	destAsset, err := loader.LoadAsset(destCode)
	if err != nil {
		return 0, false, errors.From(errors.New("failed to select dest asset"), logan.F{
			"destCode": destCode,
		})
	}

	if destAsset == nil {
		return 0, false, errors.From(errors.New("asset not found"), logan.F{
			"destCode": destCode,
		})
	}

	switch destCode {
	case pair.QuoteAsset:
		result, isOverflow := amount.BigDivide(amountToConvert, pair.CurrentPrice,
			amount.One, amount.ROUND_UP, destAsset.GetMinimumAmount())
		return result, !isOverflow, nil
	case pair.BaseAsset:
		result, isOverflow := amount.BigDivide(amountToConvert, amount.One,
			pair.CurrentPrice, amount.ROUND_UP, destAsset.GetMinimumAmount())
		return result, !isOverflow, nil
	default:
		return 0, false, errors.From(errors.New("unexpected dest code"), logan.F{
			"base":        pair.BaseAsset,
			"quote":       pair.QuoteAsset,
			"actual dest": destCode,
		})
	}
}

func (pair AssetPair) IsSimilar(other AssetPair) bool {
	return (pair.BaseAsset == other.BaseAsset && pair.QuoteAsset == other.QuoteAsset) ||
		(pair.BaseAsset == other.QuoteAsset && pair.QuoteAsset == other.BaseAsset)
}

// ConvertFromSourceAsset - converts specified amount from source to another asset in pair using current price,
// returns false - if failed
func (pair AssetPair) ConvertFromSourceAsset(sourceCode string, amountToConvert int64, loader assetLoader,
) (int64, bool, error) {
	destCode := ""
	switch sourceCode {
	case pair.BaseAsset:
		destCode = pair.QuoteAsset
	case pair.QuoteAsset:
		destCode = pair.BaseAsset
	default:
		return 0, false, errors.From(errors.New("unexpected source code"), logan.F{
			"base":          pair.BaseAsset,
			"quote":         pair.QuoteAsset,
			"actual source": sourceCode,
		})
	}

	return pair.ConvertToDestAsset(destCode, amountToConvert, loader)
}

// Contains - returns true if base or quote equal to asset
func (pair AssetPair) Contains(asset string) bool {
	return pair.BaseAsset == asset || pair.QuoteAsset == asset
}

func (pair AssetPair) IsOverlaps(anotherPair AssetPair) bool {
	return pair.Contains(anotherPair.BaseAsset) || pair.Contains(anotherPair.QuoteAsset)
}
