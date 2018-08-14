package reviewablerequest

import (
	"gitlab.com/swarmfund/horizon/db2/history"
	"gitlab.com/swarmfund/horizon/resource/base"
	"gitlab.com/tokend/go/xdr"
	"gitlab.com/tokend/regources/reviewablerequest2"
)

func PopulateAssetUpdateRequest(histRequest history.AssetUpdateRequest) (
	*reviewablerequest2.AssetUpdateRequest, error,
) {
	return &reviewablerequest2.AssetUpdateRequest{
		Code:     histRequest.Asset,
		Policies: base.FlagFromXdrAssetPolicy(histRequest.Policies, xdr.AssetPolicyAll),
		Details:  histRequest.Details,
	}, nil
}
