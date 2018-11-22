package operaitons

import (
	"gitlab.com/tokend/go/amount"
	"gitlab.com/tokend/go/xdr"
	"gitlab.com/tokend/horizon/db2/history2"
)

type createPreIssuanceRequestOpHandler struct {
}

func (h *createPreIssuanceRequestOpHandler) OperationDetails(opBody xdr.OperationBody,
	opRes xdr.OperationResultTr,
) (history2.OperationDetails, error) {
	preissuanceRequest := opBody.MustCreatePreIssuanceRequest().Request
	successResult := opRes.MustCreatePreIssuanceRequestResult().MustSuccess()

	return history2.OperationDetails{
		Type: xdr.OperationTypeCreatePreissuanceRequest,
		CreatePreIssuanceRequest: &history2.CreatePreIssuanceRequestDetails{
			AssetCode:   preissuanceRequest.Asset,
			Amount:      amount.StringU(uint64(preissuanceRequest.Amount)),
			RequestID:   int64(successResult.RequestId),
			IsFulfilled: successResult.Fulfilled,
		},
	}, nil
}

func (h *createPreIssuanceRequestOpHandler) ParticipantsEffects(opBody xdr.OperationBody,
	opRes xdr.OperationResultTr, source history2.ParticipantEffect,
) ([]history2.ParticipantEffect, error) {
	return []history2.ParticipantEffect{source}, nil
}
