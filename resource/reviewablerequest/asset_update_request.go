package reviewablerequest

import (
	"gitlab.com/tokend/go/xdr"
	"gitlab.com/swarmfund/horizon/db2/history"
	"gitlab.com/swarmfund/horizon/resource/base"
)

type AssetUpdateRequest struct {
	Code     string                 `json:"code"`
	Policies []base.Flag            `json:"policies"`
	Details  map[string]interface{} `json:"details"`
}

func (r *AssetUpdateRequest) Populate(histRequest history.AssetUpdateRequest) error {
	r.Code = histRequest.Asset
	r.Policies = base.FlagFromXdrAssetPolicy(histRequest.Policies, xdr.AssetPolicyAll)
	r.Details = histRequest.Details
	return nil
}
