package operaitons

import (
	"gitlab.com/tokend/go/xdr"
	"gitlab.com/tokend/horizon/db2/history2"
)

type createKYCRequestOpHandler struct {
	pubKeyProvider publicKeyProvider
}

func (h *createKYCRequestOpHandler) OperationDetails(op rawOperation, opRes xdr.OperationResultTr,
) (history2.OperationDetails, error) {
	createKYCRequestOp := op.Body.MustCreateUpdateKycRequestOp().UpdateKycRequestData
	createKYCRequestRes := opRes.MustCreateUpdateKycRequestResult().MustSuccess()

	var allTasks *uint32
	if createKYCRequestOp.AllTasks != nil {
		allTasksInt := uint32(*createKYCRequestOp.AllTasks)
		allTasks = &allTasksInt
	}

	return history2.OperationDetails{
		Type: xdr.OperationTypeCreateKycRequest,
		CreateKYCRequest: &history2.CreateKYCRequestDetails{
			AccountIDToUpdateKYC: createKYCRequestOp.AccountToUpdateKyc.Address(),
			AccountTypeToSet:     createKYCRequestOp.AccountTypeToSet,
			KYCData:              customDetailsUnmarshal([]byte(createKYCRequestOp.KycData)),
			AllTasks:             allTasks,
			RequestDetails: history2.RequestDetails{
				RequestID:   int64(createKYCRequestRes.RequestId),
				IsFulfilled: createKYCRequestRes.Fulfilled,
			},
		},
	}, nil
}

func (h *createKYCRequestOpHandler) ParticipantsEffects(opBody xdr.OperationBody,
	opRes xdr.OperationResultTr, source history2.ParticipantEffect,
) ([]history2.ParticipantEffect, error) {
	createKYCRequestOp := opBody.MustCreateUpdateKycRequestOp().UpdateKycRequestData

	accountIDToUpdateKYC := h.pubKeyProvider.GetAccountID(createKYCRequestOp.AccountToUpdateKyc)

	if accountIDToUpdateKYC == source.AccountID {
		return []history2.ParticipantEffect{source}, nil
	}

	return []history2.ParticipantEffect{source, {
		AccountID: accountIDToUpdateKYC,
	}}, nil
}
