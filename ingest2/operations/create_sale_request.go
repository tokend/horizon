package operations

import (
	"gitlab.com/tokend/go/xdr"
	"gitlab.com/tokend/horizon/db2/history2"
	"gitlab.com/tokend/horizon/ingest2/internal"
	"gitlab.com/tokend/regources/v2"
)

type createSaleRequestOpHandler struct {
}

// Details returns details about create sale request operation
func (h *createSaleRequestOpHandler) Details(op rawOperation, opRes xdr.OperationResultTr,
) (regources.OperationDetails, error) {
	createSaleRequest := op.Body.MustCreateSaleCreationRequestOp().Request

	return regources.OperationDetails{
		Type: xdr.OperationTypeCreateSaleRequest,
		CreateSaleRequest: &regources.CreateSaleRequestDetails{
			RequestID:         int64(opRes.MustCreateSaleCreationRequestResult().MustSuccess().RequestId),
			BaseAsset:         string(createSaleRequest.BaseAsset),
			DefaultQuoteAsset: string(createSaleRequest.DefaultQuoteAsset),
			StartTime:         internal.ToTime(createSaleRequest.StartTime),
			EndTime:           internal.ToTime(createSaleRequest.EndTime),
			HardCap:           regources.Amount(createSaleRequest.HardCap),
			SoftCap:           regources.Amount(createSaleRequest.SoftCap),
			Details:           []byte(createSaleRequest.Details),
		},
	}, nil
}

//ParticipantsEffects returns source participant effect
func (h *createSaleRequestOpHandler) ParticipantsEffects(opBody xdr.OperationBody,
	opRes xdr.OperationResultTr, source history2.ParticipantEffect, _ []xdr.LedgerEntryChange,
) ([]history2.ParticipantEffect, error) {
	return []history2.ParticipantEffect{source}, nil
}
