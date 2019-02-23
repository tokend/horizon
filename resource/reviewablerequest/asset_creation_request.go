package reviewablerequest

import (
	"gitlab.com/tokend/go/amount"
	"gitlab.com/tokend/go/xdr"
	"gitlab.com/tokend/horizon/db2/history"
	"gitlab.com/tokend/horizon/resource/base"
	"gitlab.com/tokend/regources"
)

func PopulateAssetCreationRequest(histRequest history.AssetCreationRequest) (
	*regources.AssetCreationRequest,
	error,
) {
	maxIssuanceAmount := amount.MustParse(histRequest.MaxIssuanceAmount)
	initialPreissuedAmount := amount.MustParse(histRequest.InitialPreissuedAmount)
	return &regources.AssetCreationRequest{
		Code:                   histRequest.Asset,
		Type:                   histRequest.Type,
		Policies:               base.FlagFromXdrAssetPolicy(histRequest.Policies, xdr.AssetPolicyAll),
		PreIssuedAssetSigner:   histRequest.PreIssuedAssetSigner,
		MaxIssuanceAmount:      regources.Amount(maxIssuanceAmount),
		InitialPreissuedAmount: regources.Amount(initialPreissuedAmount),
		Details:                histRequest.Details,
	}, nil

}
