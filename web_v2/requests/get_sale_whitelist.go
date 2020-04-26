package requests

import (
	"gitlab.com/tokend/horizon/bridge"
	"net/http"
)

const (
	// FilterTypeSaleWhitelistAddress - defines if we need to filter entries by account address they are applied to
	FilterTypeSaleWhitelistAddress = "address"
)

// GetSaleWhitelist - represents params to be specified by user for getSaleWhitelist handler
type GetSaleWhitelist struct {
	*base
	SaleID  uint64
	Filters struct {
		Address string `fig:"address"`
	}
	PageParams *bridge.CursorPageParams
}

// NewGetSaleWhitelist returns new instance of GetSaleWhitelist
func NewGetSaleWhitelist(r *http.Request) (*GetSaleWhitelist, error) {
	b, err := newBase(r, baseOpts{
		supportedFilters: map[string]struct{}{
			FilterTypeSaleWhitelistAddress: {},
		},
	})
	if err != nil {
		return nil, err
	}

	id, err := b.getUint64ID()
	if err != nil {
		return nil, err
	}
	pageParams, err := b.getCursorBasedPageParams()
	if err != nil {
		return nil, err
	}

	request := GetSaleWhitelist{
		base:       b,
		SaleID:     id,
		PageParams: pageParams,
	}

	err = b.populateFilters(&request.Filters)
	if err != nil {
		return nil, err
	}

	return &request, nil
}
