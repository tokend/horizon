package requests

import (
	"net/http"

	"gitlab.com/tokend/horizon/db2"
)

//GetSignerRuleList - represents params to be specified for Get SignerRules handler
type GetSignerRuleList struct {
	*base
	PageParams *db2.OffsetPageParams
}

// NewGetSignerRuleList returns the new instance of GetSignerRuleList request
func NewGetSignerRuleList(r *http.Request) (*GetSignerRuleList, error) {
	b, err := newBase(r, baseOpts{})
	if err != nil {
		return nil, err
	}

	pageParams, err := b.getOffsetBasedPageParams()
	if err != nil {
		return nil, err
	}

	request := GetSignerRuleList{
		base:       b,
		PageParams: pageParams,
	}

	return &request, nil
}
