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
			AccountId:      record.Account,
			FirstReserve:   record.FirstReserve,
			LpTokensAmount: record.TokensAmount,
			SecondReserve:  record.SecondReserve,
		},
		Relationships: regources.LiquidityPoolRelationships{
			FirstAsset:    NewAssetKey(record.FirstAssetCode).AsRelation(),
			SecondAsset:   NewBalanceKey(record.SecondAssetCode).AsRelation(),
			LpTokensAsset: NewAssetKey(record.TokenAsset).AsRelation(),
			FirstBalance:  NewAssetKey(record.FirstBalanceID).AsRelation(),
			SecondBalance: NewBalanceKey(record.SecondBalanceID).AsRelation(),
		},
	}
}
