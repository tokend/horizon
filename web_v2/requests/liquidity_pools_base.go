package requests

const (
	IncludeTypeLiquidityPoolListAssets = "assets"

	FilterTypeLiquidityPoolListAsset   = "asset"
	FilterTypeLiquidityPoolListLPToken = "lp_token"
)

var includeTypeLiquidityPoolListAll = map[string]struct{}{
	IncludeTypeLiquidityPoolListAssets: {},
}

var filterTypeLiquidityPoolListAll = map[string]struct{}{
	FilterTypeLiquidityPoolListAsset:   {},
	FilterTypeLiquidityPoolListLPToken: {},
}

type LiquidityPoolsBase struct {
	*base
	Filters struct {
		Asset   string `filter:"asset"`
		LPToken string `filter:"lp_token"`
	}
	Includes struct {
		Assets bool `include:"assets"`
	}
}
