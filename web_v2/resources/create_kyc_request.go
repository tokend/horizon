package resources

import (
	"gitlab.com/tokend/horizon/db2/history2"
	"gitlab.com/tokend/regources/v2"
)

func newKeyRequest(id int64, details history2.CreateChangeRoleRequestDetails) *regources.CreateKYCRequest {
	return &regources.CreateKYCRequest{
		Key: regources.NewKeyInt64(id, regources.TypeCreateKycRequest),
		Attributes: regources.CreateChangeRoleRequestAttrs{
			DestinationAccount: details.DestinationAccount,
			AccountRoleToSet:   details.AccountRoleToSet,
			KYCData:            details.KYCData,
			AllTasks:           details.AllTasks,
			RequestDetails:     regources.Request(details.RequestDetails),
		},
	}
}
