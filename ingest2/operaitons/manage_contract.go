package operaitons

import (
	"gitlab.com/distributed_lab/logan/v3/errors"
	"gitlab.com/tokend/go/xdr"
	"gitlab.com/tokend/horizon/db2/history2"
)

type manageContractOpHandler struct {
}

func (h *manageContractOpHandler) OperationDetails(opBody xdr.OperationBody, opRes xdr.OperationResultTr,
) (history2.OperationDetails, error) {
	manageContractOp := opBody.MustManageContractOp()

	opDetails := history2.OperationDetails{
		Type: xdr.OperationTypeManageContract,
		ManageContract: &history2.ManageContractDetails{
			ContractID: int64(manageContractOp.ContractId),
			Action:     manageContractOp.Data.Action,
		},
	}

	switch opDetails.ManageContract.Action {
	case xdr.ManageContractActionAddDetails:
		opDetails.ManageContract.Details = customDetailsUnmarshal([]byte(manageContractOp.Data.MustDetails()))
	case xdr.ManageContractActionConfirmCompleted:
		isCompeted := opRes.MustManageContractResult().MustResponse().Data.MustIsCompleted()

		opDetails.ManageContract.IsCompleted = &isCompeted
	case xdr.ManageContractActionStartDispute:
		opDetails.ManageContract.Details = customDetailsUnmarshal(
			[]byte(manageContractOp.Data.MustDisputeReason()))
	case xdr.ManageContractActionResolveDispute:
		isRevert := manageContractOp.Data.MustIsRevert()

		opDetails.ManageContract.IsRevert = &isRevert
	default:
		return history2.OperationDetails{}, errors.New("unexpected manage contract actions")
	}

	return opDetails, nil
}

func (h *manageContractOpHandler) ParticipantsEffects(opBody xdr.OperationBody,
	opRes xdr.OperationResultTr, source history2.ParticipantEffect,
) ([]history2.ParticipantEffect, error) {
	return []history2.ParticipantEffect{source}, nil
}
