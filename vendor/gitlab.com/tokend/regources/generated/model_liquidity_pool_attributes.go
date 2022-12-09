/*
 * GENERATED. Do not modify. Your changes might be overwritten!
 */

package regources

type LiquidityPoolAttributes struct {
	// Liquidity pool's account ID
	AccountId string `json:"account_id"`
	// Amount of first asset's reserve in liquidity pool
	FirstReserve Amount `json:"first_reserve"`
	// Supply of a liquidity pool tokens
	LpTokensAmount Amount `json:"lp_tokens_amount"`
	// Amount of second asset's reserve in liquidity pool
	SecondReserve Amount `json:"second_reserve"`
}
