package operations

import (
	"gitlab.com/tokend/go/xdr"
	"gitlab.com/tokend/horizon/db2/history2"
	"gitlab.com/tokend/regources/v2"
)

type createKYCRequestOpHandler struct {
	pubKeyProvider IDProvider
}

// Details returns details about create KYC request operation
func (h *createKYCRequestOpHandler) Details(op rawOperation, opRes xdr.OperationResultTr,
) (regources.OperationDetails, error) {
	createKYCRequestOp := op.Body.MustCreateUpdateKycRequestOp()
	createKYCRequestRes := opRes.MustCreateUpdateKycRequestResult().MustSuccess()

	var allTasks *uint32
	if createKYCRequestOp.AllTasks != nil {
		allTasksInt := uint32(*createKYCRequestOp.AllTasks)
		allTasks = &allTasksInt
	}

	return regources.OperationDetails{
		Type: xdr.OperationTypeCreateKycRequest,
		CreateKYCRequest: &regources.CreateKYCRequestDetails{
			AccountAddressToUpdateKYC: createKYCRequestOp.UpdateKycRequestData.AccountToUpdateKyc.Address(),
			AccountTypeToSet:          createKYCRequestOp.UpdateKycRequestData.AccountTypeToSet,
			KYCData:                   []byte(createKYCRequestOp.UpdateKycRequestData.KycData),
			AllTasks:                  allTasks,
			RequestDetails: regources.RequestDetails{
				RequestID:   int64(createKYCRequestRes.RequestId),
				IsFulfilled: createKYCRequestRes.Fulfilled,
			},
		},
	}, nil
}

//ParticipantsEffects returns source participant effect and effect for account for which KYC is updated
func (h *createKYCRequestOpHandler) ParticipantsEffects(opBody xdr.OperationBody,
	opRes xdr.OperationResultTr, source history2.ParticipantEffect, _ []xdr.LedgerEntryChange,
) ([]history2.ParticipantEffect, error) {
	createKYCRequestOp := opBody.MustCreateUpdateKycRequestOp().UpdateKycRequestData

	accountIDToUpdateKYC := h.pubKeyProvider.MustAccountID(createKYCRequestOp.AccountToUpdateKyc)

	if accountIDToUpdateKYC == source.AccountID {
		return []history2.ParticipantEffect{source}, nil
	}

	return []history2.ParticipantEffect{source, {
		AccountID: accountIDToUpdateKYC,
	}}, nil
}
