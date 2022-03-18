package history2

import regources "gitlab.com/tokend/regources/generated"

type LiquidityPool struct {
	ID              int64            `db:"id"`
	Account         string           `db:"account"`
	TokenAsset      string           `db:"token_asset"`
	FirstBalanceID  string           `db:"first_balance"`
	SecondBalanceID string           `db:"second_balance"`
	TokensAmount    regources.Amount `db:"tokens_amount"`
	FirstReserve    regources.Amount `db:"first_reserve"`
	SecondReserve   regources.Amount `db:"second_reserve"`
	FirstAssetCode  string           `db:"first_asset_code"`
	SecondAssetCode string           `db:"second_asset_code"`

	FirstAsset   *Asset `db:"first_asset"`
	SecondAsset  *Asset `db:"second_asset"`
	LPTokenAsset *Asset `db:"lp_tokens_asset"`
}
