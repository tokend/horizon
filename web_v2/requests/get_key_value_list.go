package requests

import (
	"gitlab.com/tokend/horizon/db2"
	"net/http"
)

type GetKeyValueList struct {
	*base
	PageParams *db2.OffsetPageParams
}

func NewGetKeyValueList(r *http.Request) (*GetKeyValueList, error) {
	b, err := newBase(r, baseOpts{})
	if err != nil {
		return nil, err
	}

	pageParams, err := b.getOffsetBasedPageParams()
	if err != nil {
		return nil, err
	}

	return &GetKeyValueList{
		base:       b,
		PageParams: pageParams,
	}, nil
}
