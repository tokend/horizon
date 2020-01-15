package requests

import "net/http"

type GetInvestments struct {
	*base
	AssetCode      string
	AccountAddress string
}

func NewGetInvestments(r *http.Request) (*GetInvestments, error) {
	b, err := newBase(r, baseOpts{})
	if err != nil {
		return nil, err
	}

	accountAddress := b.getString("id")
	assetCode := b.getString("asset_code")

	request := &GetInvestments{
		base:           b,
		AccountAddress: accountAddress,
		AssetCode:      assetCode,
	}

	return request, nil
}
