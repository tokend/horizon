package reviewablerequest

import (
	"gitlab.com/swarmfund/horizon/resource/base"
	"gitlab.com/swarmfund/horizon/db2/history"
	"gitlab.com/swarmfund/go/xdr"
	"encoding/json"
	"gitlab.com/distributed_lab/logan/v3/errors"
)

type AssetUpdateRequest struct {
	Code                 string `json:"code"`
	Description          string `json:"description"`
	ExternalResourceLink string `json:"external_resource_link"`
	Policies             []base.Flag `json:"policies"`
}

func (r *AssetUpdateRequest) Populate(histRequest history.AssetUpdateRequest) {
	r.Code = histRequest.Code
	r.Description = histRequest.Description
	r.ExternalResourceLink = histRequest.ExternalResourceLink
	r.Policies = base.FlagFromXdrAssetPolicy(histRequest.Policies, xdr.AssetPolicyAll)
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
