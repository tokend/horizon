package requests

import (
	"net/http"
)

const (
	// IncludeTypeConvertedBalancesStates - defines if converted balances states should be included in the response
	IncludeTypeConvertedBalancesStates = "states"
	// IncludeTypeConvertedBalancesBalance - defines if balances states should be included in the response
	IncludeTypeConvertedBalancesBalance = "balance"
	// IncludeTypeConvertedBalancesBalanceState - defines if balances states should be included in the response
	IncludeTypeConvertedBalancesBalanceState = "balance.state"
	// IncludeTypeConvertedBalancesBalanceAsset - defines if balances assets should be included in the response
	IncludeTypeConvertedBalancesBalanceAsset = "balance.asset"
)

var includeTypeConvertedBalancesAll = map[string]struct{}{
	IncludeTypeConvertedBalancesStates:       {},
	IncludeTypeConvertedBalancesBalance:      {},
	IncludeTypeConvertedBalancesBalanceState: {},
	IncludeTypeConvertedBalancesBalanceAsset: {},
}

// GetConvertedBalances - represents params to be specified by user for GetConvertedBalances handler
type GetConvertedBalances struct {
	*base
	AssetCode      string
	AccountAddress string
}

// NewGetConvertedBalances returns new instance of GetConvertedBalances request
func NewGetConvertedBalances(r *http.Request) (*GetConvertedBalances, error) {
	b, err := newBase(r, baseOpts{
		supportedIncludes: includeTypeConvertedBalancesAll,
	})
	if err != nil {
		return nil, err
	}

	accountAddress := b.getString("id")
	assetCode := b.getString("asset_code")

	return &GetConvertedBalances{
		base:           b,
		AccountAddress: accountAddress,
		AssetCode:      assetCode,
	}, nil
}
