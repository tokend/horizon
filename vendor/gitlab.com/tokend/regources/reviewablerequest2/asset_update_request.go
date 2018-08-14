package reviewablerequest2

import "gitlab.com/tokend/regources/valueflag"

type AssetUpdateRequest struct {
	Code     string                 `json:"code"`
	Policies []valueflag.Flag       `json:"policies"`
	Details  map[string]interface{} `json:"details"`
}
