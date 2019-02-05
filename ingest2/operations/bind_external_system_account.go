package operations

import (
	"gitlab.com/tokend/go/xdr"
	"gitlab.com/tokend/horizon/db2/history2"
)

type bindExternalSystemAccountOpHandler struct {
}

// CreatorDetails returns details about bind external system account operation
func (h *bindExternalSystemAccountOpHandler) Details(op rawOperation,
	opRes xdr.OperationResultTr,
) (history2.OperationDetails, error) {
	bindExternalSystemAccountOp := op.Body.MustBindExternalSystemAccountIdOp()

	return history2.OperationDetails{
		Type: xdr.OperationTypeBindExternalSystemAccountId,
		BindExternalSystemAccount: &history2.BindExternalSystemAccountDetails{
			ExternalSystemType: int32(bindExternalSystemAccountOp.ExternalSystemType),
		},
	}, nil
}

// ParticipantsEffects returns only source without effects
func (h *bindExternalSystemAccountOpHandler) ParticipantsEffects(opBody xdr.OperationBody,
	opRes xdr.OperationResultTr, source history2.ParticipantEffect, _ []xdr.LedgerEntryChange,
) ([]history2.ParticipantEffect, error) {
	return []history2.ParticipantEffect{source}, nil
}
