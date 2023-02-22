package requests

const (
	IncludeTypeLiquidityPoolListAssets = "assets"

	FilterTypeLiquidityPoolListAsset          = "asset"
	FilterTypeLiquidityPoolListLPToken        = "lp_token"
	FilterTypeLiquidityPoolListBalancesOwner  = "balances_owner"
	FilterTypeLiquidityPoolListExcludedAssets = "excluded_assets"
)

var includeTypeLiquidityPoolListAll = map[string]struct{}{
	IncludeTypeLiquidityPoolListAssets: {},
}

var filterTypeLiquidityPoolListAll = map[string]struct{}{
	FilterTypeLiquidityPoolListAsset:          {},
	FilterTypeLiquidityPoolListLPToken:        {},
	FilterTypeLiquidityPoolListBalancesOwner:  {},
	FilterTypeLiquidityPoolListExcludedAssets: {},
}

type LiquidityPoolsBase struct {
	*base
	Filters struct {
		Asset          string   `filter:"asset"`
		LPToken        string   `filter:"lp_token"`
		BalancesOwner  string   `filter:"balances_owner"`
		ExcludedAssets []string `filter:"excluded_assets"`
	}
	Includes struct {
		Assets bool `include:"assets"`
	}
}
