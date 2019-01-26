package regources

import "gitlab.com/tokend/go/xdr"

//ManageAssetAttrs - details of corresponding op
type ManageAsset struct {
	Key
	Attributes ManageAssetAttrs `json:"attributes"`
}

//ManageAssetAttrs - details of corresponding op
type ManageAssetAttrs struct {
	AssetCode         string                `json:"asset_code,omitempty"`
	RequestID         int64                 `json:"request_id"`
	Action            xdr.ManageAssetAction `json:"action"`
	Policies          *xdr.AssetPolicy      `json:"policies,omitempty"`
	Details           Details               `json:"details,omitempty"`
	PreissuedSigner   string                `json:"preissuance_signer,omitempty"`
	MaxIssuanceAmount Amount                `json:"max_issuance_amount,omitempty"`
}
