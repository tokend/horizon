package operaitons

import (
	"gitlab.com/tokend/go/amount"
	"gitlab.com/tokend/go/xdr"
	"gitlab.com/tokend/horizon/db2/history2"
	"gitlab.com/tokend/horizon/utf8"
)

type createIssuanceRequestOpHandler struct {
	pubKeyProvider publicKeyProvider
}

func (h *createIssuanceRequestOpHandler) OperationDetails(opBody xdr.OperationBody,
	opRes xdr.OperationResultTr,
) (history2.OperationDetails, error) {
	createIssuanceRequestOp := opBody.MustCreateIssuanceRequestOp()
	issuanceRequest := createIssuanceRequestOp.Request

	var allTasks *int64
	rawAllTasks, ok := createIssuanceRequestOp.Ext.GetAllTasks()
	if ok && rawAllTasks != nil {
		allTasksInt := int64(*rawAllTasks)
		allTasks = &allTasksInt
	}

	return history2.OperationDetails{
		Type: xdr.OperationTypeCreateIssuanceRequest,
		CreateIssuanceRequest: &history2.CreateIssuanceRequestDetails{
			FixedFee:        amount.StringU(uint64(issuanceRequest.Fee.Fixed)),
			PercentFee:      amount.StringU(uint64(issuanceRequest.Fee.Percent)),
			Reference:       utf8.Scrub(string(createIssuanceRequestOp.Reference)),
			Amount:          amount.StringU(uint64(issuanceRequest.Amount)),
			Asset:           issuanceRequest.Asset,
			BalanceID:       h.pubKeyProvider.GetBalanceID(issuanceRequest.Receiver),
			ExternalDetails: customDetailsUnmarshal([]byte(issuanceRequest.ExternalDetails)),
			AllTasks:        allTasks,
		},
	}, nil
}

func (h *createIssuanceRequestOpHandler) ParticipantsEffects(opBody xdr.OperationBody,
	opRes xdr.OperationResultTr, source history2.ParticipantEffect,
) ([]history2.ParticipantEffect, error) {

}
