package resources

import "gitlab.com/tokend/regources/v2"

//NewRequestKey - creates new instance of request key
func NewRequestKey(requestId int64) regources.Key {
	return regources.NewKeyInt64(requestId, regources.TypeRequests)
}
