package reviewablerequest

import (
	"gitlab.com/swarmfund/horizon/db2/history"
	"gitlab.com/swarmfund/horizon/resource/base"
	"gitlab.com/swarmfund/go/xdr"
	"encoding/json"
	"gitlab.com/distributed_lab/logan/v3/errors"
)

type AssetCreationRequest struct {
	Code                 string      `json:"code"`
	Description          string      `json:"description"`
	ExternalResourceLink string      `json:"external_resource_link"`
	Policies             []base.Flag `json:"policies"`
	Name                 string      `json:"name"`
	PreIssuedAssetSigner string      `json:"pre_issued_asset_signer"`
	MaxIssuanceAmount    string      `json:"max_issuance_amount"`
}

func (r *AssetCreationRequest) Populate(histRequest history.AssetCreationRequest) {
	r.Code = histRequest.Code
	r.Description = histRequest.Description
	r.ExternalResourceLink = histRequest.ExternalResourceLink
	r.Policies = base.FlagFromXdrAssetPolicy(histRequest.Policies, xdr.AssetPolicyAll)
	r.Name = histRequest.Name
	r.PreIssuedAssetSigner = histRequest.PreIssuedAssetSigner
	r.MaxIssuanceAmount = histRequest.MaxIssuanceAmount
}

func (r *AssetCreationRequest) PopulateFromRawJsonHistory(rawJson []byte) error {
	var histRequest history.AssetCreationRequest
	err := json.Unmarshal(rawJson, &histRequest)
	if err != nil {
		return errors.Wrap(err, "failed to unmarshal history.AssetCreationRequest")
	}

	r.Populate(histRequest)
	return nil
}
