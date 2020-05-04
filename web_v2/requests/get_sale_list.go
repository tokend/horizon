package requests

import (
	"gitlab.com/tokend/horizon/db2"
	"net/http"
)

const (
	// FilterTypeSaleListOwner - defines if we need to filter response by participant
	FilterTypeSaleListParticipant = "participant"
)

var filterTypeSaleListAllWithParticipant = map[string]struct{}{
	FilterTypeSaleListOwner:        {},
	FilterTypeSaleListBaseAsset:    {},
	FilterTypeSaleListMaxEndTime:   {},
	FilterTypeSaleListMaxStartTime: {},
	FilterTypeSaleListMinStartTime: {},
	FilterTypeSaleListMinEndTime:   {},
	FilterTypeSaleListState:        {},
	FilterTypeSaleListSaleType:     {},
	FilterTypeSaleListMinHardCap:   {},
	FilterTypeSaleListMinSoftCap:   {},
	FilterTypeSaleListMaxHardCap:   {},
	FilterTypeSaleListMaxSoftCap:   {},
	FilterTypeSaleListParticipant:  {},
}

// GetSaleList - represents params to be specified by user for getSaleList handler
type GetSaleList struct {
	SalesBase
	SpecialFilters struct {
		Participant string `json:"participant"`
	}
	PageParams *db2.OffsetPageParams
}

// NewGetSaleList returns new instance of GetSaleList request
func NewGetSaleList(r *http.Request) (*GetSaleList, error) {
	b, err := newBase(r, baseOpts{
		supportedIncludes: includeTypeSaleListAll,
		supportedFilters:  filterTypeSaleListAllWithParticipant,
	})
	if err != nil {
		return nil, err
	}

	pageParams, err := b.getOffsetBasedPageParams()
	if err != nil {
		return nil, err
	}

	request := GetSaleList{
		SalesBase: SalesBase{
			base: b,
		},
		PageParams: pageParams,
	}

	err = b.populateFilters(&request.Filters)
	if err != nil {
		return nil, err
	}

	err = b.populateFilters(&request.SpecialFilters)
	if err != nil {
		return nil, err
	}

	return &request, nil
}

func (g GetSaleList) GetLoganFields() map[string]interface{} {
	return map[string]interface{}{
		"participant": g.SpecialFilters.Participant,
	}
}
