package reviewablerequest

import (
	"encoding/json"
	"gitlab.com/distributed_lab/logan/v3/errors"
	"gitlab.com/swarmfund/go/xdr"
	"gitlab.com/swarmfund/horizon/db2/history"
	"gitlab.com/swarmfund/horizon/resource/base"
)

type AssetUpdateRequest struct {
	Code     string                 `json:"code"`
	Policies []base.Flag            `json:"policies"`
	Details  map[string]interface{} `json:"details"`
}

func (r *AssetUpdateRequest) Populate(histRequest history.AssetUpdateRequest) {
	r.Code = histRequest.Asset
	r.Policies = base.FlagFromXdrAssetPolicy(histRequest.Policies, xdr.AssetPolicyAll)
	r.Details = histRequest.Details
}

func (r *AssetUpdateRequest) PopulateFromRawJsonHistory(rawJson []byte) error {
	var histRequest history.AssetUpdateRequest
	err := json.Unmarshal(rawJson, &histRequest)
	if err != nil {
		return errors.Wrap(err, "failed to unmarshal history.AssetUpdateRequest")
	}

	r.Populate(histRequest)
	return nil
}
