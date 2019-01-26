package resources

import (
	"gitlab.com/tokend/horizon/db2/history2"
	"gitlab.com/tokend/regources/v2"
)

func newCreateIssuanceOpDetails(id int64, details history2.CreateIssuanceRequestDetails) *regources.CreateIssuanceRequest {
	return &regources.CreateIssuanceRequest{
		Key: regources.NewKeyInt64(id, regources.TypeCreateIssuanceRequest),
		Attributes: regources.CreateIssuanceRequestAttrs{
			Fee:       details.Fee,
			Reference: details.Reference,
			Amount:    details.Amount,
			Asset:     details.Asset,
			ReceiverAccountAddress: details.ReceiverAccountAddress,
			ReceiverBalanceAddress: details.ReceiverBalanceAddress,
			ExternalDetails:        details.ExternalDetails,
			AllTasks:               details.AllTasks,
			RequestDetails:         regources.Request(details.RequestDetails),
		},
	}
}
