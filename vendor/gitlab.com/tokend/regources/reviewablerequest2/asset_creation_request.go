package reviewablerequest2

import "gitlab.com/tokend/regources/valueflag"

type AssetCreationRequest struct {
	Code                   string                 `json:"code"`
	Policies               []valueflag.Flag       `json:"policies"`
	PreIssuedAssetSigner   string                 `json:"pre_issued_asset_signer"`
	MaxIssuanceAmount      string                 `json:"max_issuance_amount"`
	InitialPreissuedAmount string                 `json:"initial_preissued_amount"`
	Details                map[string]interface{} `json:"details"`
}
