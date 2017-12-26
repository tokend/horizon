package reviewablerequest

import (
	"encoding/json"
	"gitlab.com/distributed_lab/logan/v3/errors"
	"gitlab.com/swarmfund/go/xdr"
	"gitlab.com/swarmfund/horizon/db2/history"
	"gitlab.com/swarmfund/horizon/resource/base"
)

type AssetCreationRequest struct {
	Code                 string                 `json:"code"`
	Policies             []base.Flag            `json:"policies"`
	PreIssuedAssetSigner string                 `json:"pre_issued_asset_signer"`
	MaxIssuanceAmount    string                 `json:"max_issuance_amount"`
	Details              map[string]interface{} `json:"details"`
}

func (r *AssetCreationRequest) Populate(histRequest history.AssetCreationRequest) {
	r.Code = histRequest.Asset
	r.Policies = base.FlagFromXdrAssetPolicy(histRequest.Policies, xdr.AssetPolicyAll)
	r.PreIssuedAssetSigner = histRequest.PreIssuedAssetSigner
	r.MaxIssuanceAmount = histRequest.MaxIssuanceAmount
	r.Details = histRequest.Details
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
