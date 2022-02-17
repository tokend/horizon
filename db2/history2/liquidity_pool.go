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
	FirstAsset      string           `db:"first_asset"`
	SecondAsset     string           `db:"second_asset"`
}
