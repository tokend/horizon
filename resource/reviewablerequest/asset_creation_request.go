package reviewablerequest

import (
	"gitlab.com/swarmfund/horizon/db2/history"
	"gitlab.com/swarmfund/horizon/resource/base"
	"gitlab.com/tokend/go/xdr"
	"gitlab.com/tokend/regources"
)

type AssetCreationRequest struct {
	Code                   string                 `json:"code"`
	Policies               []regources.Flag       `json:"policies"`
	PreIssuedAssetSigner   string                 `json:"pre_issued_asset_signer"`
	MaxIssuanceAmount      string                 `json:"max_issuance_amount"`
	InitialPreissuedAmount string                 `json:"initial_preissued_amount"`
	Details                map[string]interface{} `json:"details"`
}

func (r *AssetCreationRequest) Populate(histRequest history.AssetCreationRequest) error {
	r.Code = histRequest.Asset
	r.Policies = base.FlagFromXdrAssetPolicy(histRequest.Policies, xdr.AssetPolicyAll)
	r.PreIssuedAssetSigner = histRequest.PreIssuedAssetSigner
	r.MaxIssuanceAmount = histRequest.MaxIssuanceAmount
	r.InitialPreissuedAmount = histRequest.InitialPreissuedAmount
	r.Details = histRequest.Details
	return nil
}
