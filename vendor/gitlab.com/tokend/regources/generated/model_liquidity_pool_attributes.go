/*
 * GENERATED. Do not modify. Your changes might be overwritten!
 */

package regources

type LiquidityPoolAttributes struct {
	// Liquidity pool's account ID
	AccountId string `json:"account_id"`
	// Code of first asset in liquidity pool
	FirstAssetCode string `json:"first_asset_code"`
	// Balance ID of first asset in liquidity pool
	FirstBalance string `json:"first_balance"`
	// Amount of first asset's reserve in liquidity pool
	FirstReserve Amount `json:"first_reserve"`
	// Supply of a liquidity pool tokens
	LpTokensAmount Amount `json:"lp_tokens_amount"`
	// Code of second asset in liquidity pool
	SecondAssetCode string `json:"second_asset_code"`
	// Balance ID of second asset in liquidity pool
	SecondBalance string `json:"second_balance"`
	// Amount of second asset's reserve in liquidity pool
	SecondReserve Amount `json:"second_reserve"`
}
