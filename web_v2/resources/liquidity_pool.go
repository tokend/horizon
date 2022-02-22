package resources

import (
	"gitlab.com/tokend/horizon/db2/history2"

	regources "gitlab.com/tokend/regources/generated"
)

// NewLiquidityPoolKey - creates new Key for liquidity pool
func NewLiquidityPoolKey(id int64) regources.Key {
	return regources.NewKeyInt64(id, regources.LIQUIDITY_POOLS)
}

// NewLiquidityPool - creates a new instance of regources.LiquidityPool
func NewLiquidityPool(record history2.LiquidityPool) regources.LiquidityPool {
	return regources.LiquidityPool{
		Key: NewLiquidityPoolKey(record.ID),
		Attributes: regources.LiquidityPoolAttributes{
			AccountId:       record.Account,
			FirstAssetCode:  record.FirstAsset,
			FirstReserve:    record.FirstReserve,
			LpTokensAmount:  record.TokensAmount,
			LpTokensAsset:   record.TokenAsset,
			SecondAssetCode: record.SecondAsset,
			SecondReserve:   record.SecondReserve,
			FirstBalance:    record.FirstBalanceID,
			SecondBalance:   record.SecondBalanceID,
		},
	}
}
