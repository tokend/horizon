package operaitons

import (
	"gitlab.com/tokend/go/xdr"
	"gitlab.com/tokend/horizon/db2/history2"
)

type checkSaleStateOpHandler struct {
}

func (h *checkSaleStateOpHandler) OperationDetails(opBody xdr.OperationBody,
	opRes xdr.OperationResultTr,
) (history2.OperationDetails, error) {

	return history2.OperationDetails{
		Type: xdr.OperationTypeCheckSaleState,
		CheckSaleState: &history2.CheckSaleStateDetails{
			SaleID: int64(opBody.MustCheckSaleStateOp().SaleId),
			Effect: opRes.MustCheckSaleStateResult().MustSuccess().Effect.Effect,
		},
	}, nil
}

func (h *checkSaleStateOpHandler) ParticipantsEffects(opBody xdr.OperationBody,
	opRes xdr.OperationResultTr, source history2.ParticipantEffect,
) ([]history2.ParticipantEffect, error) {
	return []history2.ParticipantEffect{source}, nil
}
