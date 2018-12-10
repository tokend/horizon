package operaitons

import (
	"gitlab.com/distributed_lab/logan/v3/errors"
	"gitlab.com/tokend/go/xdr"
	"gitlab.com/tokend/horizon/db2/history2"
)

type manageContractRequestOpHandler struct {
	pubKeyProvider publicKeyProvider
}

// OperationDetails returns details about manage contract request operation
func (h *manageContractRequestOpHandler) OperationDetails(op rawOperation,
	opRes xdr.OperationResultTr,
) (history2.OperationDetails, error) {
	manageContractRequestOp := op.Body.MustManageContractRequestOp()

	opDetails := history2.OperationDetails{
		Type: xdr.OperationTypeManageContractRequest,
		ManageContractRequest: &history2.ManageContractRequestDetails{
			Action: manageContractRequestOp.Details.Action,
		},
	}

	switch opDetails.ManageContractRequest.Action {
	case xdr.ManageContractRequestActionCreate:
		creationDetails := manageContractRequestOp.Details.MustContractRequest()

		opDetails.ManageContractRequest.Create = &history2.CreateContractRequestDetails{
			Customer:  creationDetails.Customer.Address(),
			Escrow:    creationDetails.Escrow.Address(),
			Details:   customDetailsUnmarshal([]byte(creationDetails.Details)),
			StartTime: int64(creationDetails.StartTime),
			EndTime:   int64(creationDetails.EndTime),
			RequestID: int64(opRes.MustManageContractRequestResult().MustSuccess().Details.MustResponse().RequestId),
		}
	case xdr.ManageContractRequestActionRemove:
		opDetails.ManageContractRequest.Remove = &history2.RemoveContractReqeustDetails{
			RequestID: int64(manageContractRequestOp.Details.MustRequestId()),
		}
	default:
		return history2.OperationDetails{}, errors.New("unexpected manage contract request action")
	}

	return opDetails, nil
}

func (h *manageContractRequestOpHandler) ParticipantsEffects(opBody xdr.OperationBody,
	opRes xdr.OperationResultTr, source history2.ParticipantEffect,
) ([]history2.ParticipantEffect, error) {
	creationDetails := opBody.MustManageContractRequestOp().Details.MustContractRequest()

	participants := []history2.ParticipantEffect{source}

	participants = append(participants, history2.ParticipantEffect{
		AccountID: h.pubKeyProvider.GetAccountID(creationDetails.Customer),
	})

	participants = append(participants, history2.ParticipantEffect{
		AccountID: h.pubKeyProvider.GetAccountID(creationDetails.Escrow),
	})

	return participants, nil
}
