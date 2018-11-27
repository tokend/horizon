package operaitons

import (
	"gitlab.com/distributed_lab/logan/v3/errors"
	"gitlab.com/tokend/go/amount"
	"gitlab.com/tokend/go/xdr"
	"gitlab.com/tokend/horizon/db2/history2"
)

type manageContractOpHandler struct {
	requestProvider requestProvider
	pubKeyProvider  publicKeyProvider
}

func (h *manageContractOpHandler) OperationDetails(op rawOperation, opRes xdr.OperationResultTr,
) (history2.OperationDetails, error) {
	manageContractOp := op.Body.MustManageContractOp()

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
	manageContractOp := opBody.MustManageContractOp()
	manageContractRes := opRes.MustManageContractResult().MustResponse()

	isCompleted, ok := manageContractRes.Data.GetIsCompleted()
	if !ok && (manageContractOp.Data.Action != xdr.ManageContractActionResolveDispute) || (ok && !isCompleted) {
		return []history2.ParticipantEffect{source}, nil
	}

	request := h.requestProvider.GetInvoiceRequestsByContractID(int64(manageContractOp.ContractId))

	if isRevert, ok := manageContractOp.Data.GetIsRevert(); ok && isRevert {
		return h.getRevertedInvoiceParticipantsEffects(request), nil
	} else {
		return h.getConfirmedInvoicesParticipantsEffects(request), nil
	}
}

func (h *manageContractOpHandler) getRevertedInvoiceParticipantsEffects(requests []xdr.ReviewableRequestEntry) []history2.ParticipantEffect {
	var participants []history2.ParticipantEffect

	for _, request := range requests {
		invoiceDetails := request.Body.MustInvoiceRequest()

		if !invoiceDetails.IsApproved {
			continue
		}

		receiverBalanceID := h.pubKeyProvider.GetBalanceID(invoiceDetails.ReceiverBalance)

		participants = append(participants, history2.ParticipantEffect{
			AccountID: h.pubKeyProvider.GetAccountID(request.Requestor),
			BalanceID: &receiverBalanceID,
			AssetCode: &invoiceDetails.Asset,
			Effect: history2.Effect{
				Type: history2.EffectTypeChargedFromLocked,
				Payment: &history2.PaymentEffect{
					Amount: amount.StringU(uint64(invoiceDetails.Amount)),
				},
			},
		})

		senderBalanceID := h.pubKeyProvider.GetBalanceID(invoiceDetails.SenderBalance)

		participants = append(participants, history2.ParticipantEffect{
			AccountID: h.pubKeyProvider.GetAccountID(request.Reviewer),
			BalanceID: &senderBalanceID,
			AssetCode: &invoiceDetails.Asset,
			Effect: history2.Effect{
				Type: history2.EffectTypeFunded,
				Payment: &history2.PaymentEffect{
					Amount: amount.StringU(uint64(invoiceDetails.Amount)),
				},
			},
		})
	}

	return participants
}

func (h *manageContractOpHandler) getConfirmedInvoicesParticipantsEffects(requests []xdr.ReviewableRequestEntry) []history2.ParticipantEffect {
	var participants []history2.ParticipantEffect

	for _, request := range requests {
		invoiceDetails := request.Body.MustInvoiceRequest()

		if !invoiceDetails.IsApproved {
			continue
		}

		balanceID := h.pubKeyProvider.GetBalanceID(invoiceDetails.ReceiverBalance)

		participants = append(participants, history2.ParticipantEffect{
			AccountID: h.pubKeyProvider.GetAccountID(request.Requestor),
			BalanceID: &balanceID,
			AssetCode: &invoiceDetails.Asset,
			Effect: history2.Effect{
				Type: history2.EffectTypeUnlocked,
				Payment: &history2.PaymentEffect{
					Amount: amount.StringU(uint64(invoiceDetails.Amount)),
				},
			},
		})
	}

	return participants
}
