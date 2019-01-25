package operations

import (
	"gitlab.com/tokend/go/xdr"
	"gitlab.com/tokend/horizon/db2/history2"
	"gitlab.com/tokend/regources/v2"
)

type createPreIssuanceRequestOpHandler struct {
}

// Details returns details about create pre issuance request operation
func (h *createPreIssuanceRequestOpHandler) Details(op rawOperation,
	opRes xdr.OperationResultTr,
) (history2.OperationDetails, error) {
	preissuanceRequest := op.Body.MustCreatePreIssuanceRequest().Request
	successResult := opRes.MustCreatePreIssuanceRequestResult().MustSuccess()

	return history2.OperationDetails{
		Type: xdr.OperationTypeCreatePreissuanceRequest,
		CreatePreIssuanceRequest: &history2.CreatePreIssuanceRequestDetails{
			AssetCode:   string(preissuanceRequest.Asset),
			Amount:      regources.Amount(preissuanceRequest.Amount),
			RequestID:   int64(successResult.RequestId),
			IsFulfilled: successResult.Fulfilled,
		},
	}, nil
}

//ParticipantsEffects returns source participant effect
func (h *createPreIssuanceRequestOpHandler) ParticipantsEffects(opBody xdr.OperationBody,
	opRes xdr.OperationResultTr, source history2.ParticipantEffect, _ []xdr.LedgerEntryChange,
) ([]history2.ParticipantEffect, error) {
	return []history2.ParticipantEffect{source}, nil
}
