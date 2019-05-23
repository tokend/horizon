package operations

import (
	"gitlab.com/tokend/go/xdr"
	"gitlab.com/tokend/horizon/db2/history2"
	"gitlab.com/tokend/horizon/ingest2/internal"
)

type createKycRecoveryRequestOpHandler struct {
	effectsProvider
}

// Details returns details about create KYC request operation
func (h *createKycRecoveryRequestOpHandler) Details(op rawOperation, opRes xdr.OperationResultTr,
) (history2.OperationDetails, error) {
	createKYCRecoveryRequestOp := op.Body.MustCreateKycRecoveryRequestOp()
	createKYCRequestRes := opRes.MustCreateKycRecoveryRequestResult().MustSuccess()

	signersData := make([]history2.UpdateSignerDetails, 0, len(createKYCRecoveryRequestOp.SignersData))
	for _, signer := range createKYCRecoveryRequestOp.SignersData {
		signersData = append(signersData, history2.UpdateSignerDetails{
			Details:  internal.MarshalCustomDetails(signer.Details),
			RoleID:   uint64(signer.RoleId),
			Identity: uint32(signer.Identity),
			Weight:   uint32(signer.Weight),
		})
	}

	return history2.OperationDetails{
		Type: xdr.OperationTypeCreateKycRecoveryRequest,
		CreateKYCRecoveryRequest: &history2.CreateKYCRecoveryRequestDetails{
			TargetAccount:  createKYCRecoveryRequestOp.TargetAccount.Address(),
			SignersData:    signersData,
			CreatorDetails: internal.MarshalCustomDetails(createKYCRecoveryRequestOp.CreatorDetails),
			AllTasks:       (*uint32)(createKYCRecoveryRequestOp.AllTasks),
			RequestDetails: history2.RequestDetails{
				RequestID:   int64(createKYCRequestRes.RequestId),
				IsFulfilled: createKYCRequestRes.Fulfilled,
			},
		},
	}, nil
}

//ParticipantsEffects returns source participant effect and effect for account for which KYC is updated
func (h *createKycRecoveryRequestOpHandler) ParticipantsEffects(opBody xdr.OperationBody,
	opRes xdr.OperationResultTr, sourceAccountID xdr.AccountId, _ []xdr.LedgerEntryChange,
) ([]history2.ParticipantEffect, error) {
	createKycRecoveryRequestOp := opBody.MustCreateKycRecoveryRequestOp()

	updatedAccount := h.Participant(createKycRecoveryRequestOp.TargetAccount)

	source := h.Participant(sourceAccountID)
	if updatedAccount.AccountID == source.AccountID {
		return []history2.ParticipantEffect{source}, nil
	}

	return []history2.ParticipantEffect{source, updatedAccount}, nil
}
