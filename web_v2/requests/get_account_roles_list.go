package requests

import (
	"gitlab.com/distributed_lab/kit/pgdb"
	"net/http"
)

var includeTypeAccountRoleListAll = map[string]struct{}{
	IncludeTypeRoleRules: {},
}

//GetAccountRoleList - represents params to be specified for Get AccountRoles handler
type GetAccountRoleList struct {
	*base
	PageParams *pgdb.OffsetPageParams
}

// NewGetAccountRoleList returns the new instance of GetAccountRoleList request
func NewGetAccountRoleList(r *http.Request) (*GetAccountRoleList, error) {
	b, err := newBase(r, baseOpts{
		supportedIncludes: includeTypeAccountRoleListAll,
	})
	if err != nil {
		return nil, err
	}

	pageParams, err := b.getOffsetBasedPageParams()
	if err != nil {
		return nil, err
	}

	request := GetAccountRoleList{
		base:       b,
		PageParams: pageParams,
	}

	return &request, nil
}
