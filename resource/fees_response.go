package resource

import "gitlab.com/tokend/go/xdr"

type FeesResponse struct {
	Fees map[xdr.AssetCode][]FeeEntry `json:"fees"`
}
