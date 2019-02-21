package requests

import (
	"net/http"

	"gitlab.com/tokend/horizon/db2"
)

var includeTypeSignerRoleListAll = map[string]struct{}{
	IncludeTypeRoleRules: {},
}

//GetSignerRoleList - represents params to be specified for Get SignerRoles handler
type GetSignerRoleList struct {
	*base
	PageParams *db2.OffsetPageParams
}

// NewGetSignerRoleList returns the new instance of GetSignerRoleList request
func NewGetSignerRoleList(r *http.Request) (*GetSignerRoleList, error) {
	b, err := newBase(r, baseOpts{
		supportedIncludes: includeTypeSignerRoleListAll,
	})
	if err != nil {
		return nil, err
	}

	pageParams, err := b.getOffsetBasedPageParams()
	if err != nil {
		return nil, err
	}

	request := GetSignerRoleList{
		base:       b,
		PageParams: pageParams,
	}

	return &request, nil
}
