package reviewablerequest

import (
	"gitlab.com/swarmfund/horizon/db2/history"
	"gitlab.com/swarmfund/horizon/resource/base"
	"gitlab.com/tokend/go/xdr"
	"gitlab.com/tokend/regources"
)

func PopulateAssetCreationRequest(histRequest history.AssetCreationRequest) (
	*regources.AssetCreationRequest,
	error,
) {
	return &regources.AssetCreationRequest{
		Code:                   histRequest.Asset,
		Policies:               base.FlagFromXdrAssetPolicy(histRequest.Policies, xdr.AssetPolicyAll),
		PreIssuedAssetSigner:   histRequest.PreIssuedAssetSigner,
		MaxIssuanceAmount:      histRequest.MaxIssuanceAmount,
		InitialPreissuedAmount: histRequest.InitialPreissuedAmount,
		Details:                histRequest.Details,
	}, nil

}
