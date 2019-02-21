package requests

import (
	"net/http"

	"gitlab.com/tokend/horizon/db2"
)

//GetAccountRuleList - represents params to be specified for Get AccountRules handler
type GetAccountRuleList struct {
	*base
	PageParams *db2.OffsetPageParams
}

// NewGetAccountRuleList returns the new instance of GetAccountRuleList request
func NewGetAccountRuleList(r *http.Request) (*GetAccountRuleList, error) {
	b, err := newBase(r, baseOpts{})
	if err != nil {
		return nil, err
	}

	pageParams, err := b.getOffsetBasedPageParams()
	if err != nil {
		return nil, err
	}

	request := GetAccountRuleList{
		base:       b,
		PageParams: pageParams,
	}

	return &request, nil
}
