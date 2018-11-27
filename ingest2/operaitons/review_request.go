package operaitons

import (
	"encoding/hex"

	"gitlab.com/tokend/go/amount"

	"gitlab.com/tokend/go/xdr"
	"gitlab.com/tokend/horizon/db2/history2"
)

type reviewRequestOpHandler struct {
	pubKeyProvider        publicKeyProvider
	ledgerChangesProvider ledgerChangesProvider
	balanceProvider       balanceProvider
	paymentHelper         paymentHelper
}

func (h *reviewRequestOpHandler) OperationDetails(op rawOperation, opRes xdr.OperationResultTr,
) (history2.OperationDetails, error) {
	reviewRequestOp := op.Body.MustReviewRequestOp()
	reviewRequestOpRes := opRes.MustReviewRequestResult().MustSuccess().Ext

	opDetails := history2.OperationDetails{
		Type: xdr.OperationTypeReviewRequest,
		ReviewRequest: &history2.ReviewRequestDetails{
			RequestID:      int64(reviewRequestOp.RequestId),
			RequestType:    reviewRequestOp.RequestDetails.RequestType,
			RequestHash:    hex.EncodeToString(reviewRequestOp.RequestHash[:]),
			Action:         reviewRequestOp.Action,
			Reason:         string(reviewRequestOp.Reason),
			IsFulfilled:    true,
			RequestDetails: reviewRequestOp.RequestDetails,
		},
	}

	extRes, ok := reviewRequestOpRes.GetExtendedResult()
	if !ok {
		return opDetails, nil
	}

	opDetails.ReviewRequest.IsFulfilled = extRes.Fulfilled

	aSwapExtended, ok := extRes.TypeExt.GetASwapExtended()
	if !ok {
		return opDetails, nil
	}

	opDetails.ReviewRequest.AtomicSwapDetails = &aSwapExtended

	return opDetails, nil
}

func (h *reviewRequestOpHandler) ParticipantsEffects(opBody xdr.OperationBody,
	opRes xdr.OperationResultTr, source history2.ParticipantEffect,
) ([]history2.ParticipantEffect, error) {
	reviewRequestOp := opBody.MustReviewRequestOp()

	request := h.getReviewableRequestByID(int64(reviewRequestOp.RequestId))

	if request == nil {
		return []history2.ParticipantEffect{source}, nil
	}

	var participants []history2.ParticipantEffect

	switch reviewRequestOp.RequestDetails.RequestType {
	case xdr.ReviewableRequestTypeIssuanceCreate:
		details := request.Body.MustIssuanceRequest()

		effect := history2.Effect{
			Type: history2.EffectTypeFunded,
			Issuance: &history2.IssuanceEffect{
				Amount: amount.StringU(uint64(details.Amount)),
			},
		}

		participants = h.getParticipantEffectByBalanceID(details.Receiver, effect, source)
	case xdr.ReviewableRequestTypeWithdraw:
		details := request.Body.MustWithdrawalRequest()

		effect := history2.Effect{
			Type: history2.EffectTypeChargedFromLocked,
			Withdraw: &history2.WithdrawEffect{
				Amount: amount.StringU(uint64(details.Amount)),
			},
		}

		participants = h.getParticipantEffectByBalanceID(details.Balance, effect, source)
	case xdr.ReviewableRequestTypeAmlAlert:
		details := request.Body.MustAmlAlertRequest()

		effect := history2.Effect{
			Type: history2.EffectTypeChargedFromLocked,
			AMLAlert: &history2.AMLAlertEffect{
				Amount: amount.StringU(uint64(details.Amount)),
			},
		}

		participants = h.getParticipantEffectByBalanceID(details.BalanceId, effect, source)
	case xdr.ReviewableRequestTypeInvoice:
		paymentOp := reviewRequestOp.RequestDetails.MustBillPay().PaymentDetails
		paymentRes := opRes.MustReviewRequestResult().MustSuccess().Ext.MustPaymentV2Response()

		effect := history2.EffectTypeFunded
		if request.Body.MustInvoiceRequest().ContractId != nil {
			effect = history2.EffectTypeFundedToLocked
		}

		participants = h.paymentHelper.getParticipantsEffects(paymentOp, paymentRes, source, effect)
	}

	if source.AccountID == h.pubKeyProvider.GetAccountID(request.Requestor) {
		return participants, nil
	}

	if request.Body.Type != xdr.ReviewableRequestTypeAtomicSwap {
		return append(participants, history2.ParticipantEffect{
			AccountID: h.pubKeyProvider.GetAccountID(request.Requestor),
		}), nil
	}

	extendedResult, ok := opRes.MustReviewRequestResult().MustSuccess().Ext.GetExtendedResult()
	if !ok {
		return participants, nil
	}

	atomicSwapExtendedResult, ok := extendedResult.TypeExt.GetASwapExtended()
	if !ok {
		return participants, nil
	}

	ownerBalanceID := h.pubKeyProvider.GetBalanceID(atomicSwapExtendedResult.BidOwnerBaseBalanceId)

	participants = append(participants, history2.ParticipantEffect{
		AccountID: h.pubKeyProvider.GetAccountID(atomicSwapExtendedResult.BidOwnerId),
		BalanceID: &ownerBalanceID,
		AssetCode: &atomicSwapExtendedResult.baseAsset,
	})

	purchaserBaseBalanceId := h.pubKeyProvider.GetBalanceID(atomicSwapExtendedResult.PurchaserBaseBalanceId)

	participants = append(participants, history2.ParticipantEffect{
		AccountID: h.pubKeyProvider.GetAccountID(atomicSwapExtendedResult.PurchaserId),
		BalanceID: &purchaserBaseBalanceId,
		AssetCode: &atomicSwapExtendedResult.baseAsset,
	})

	return participants, nil
}

func (h *reviewRequestOpHandler) getParticipantEffectByBalanceID(balanceID xdr.BalanceId,
	effect history2.Effect, source history2.ParticipantEffect,
) []history2.ParticipantEffect {
	balance := h.balanceProvider.GetBalanceByID(balanceID)
	if balance.AccountID == source.AccountID {
		source.BalanceID = &balance.BalanceID
		source.AssetCode = &balance.AssetCode
		source.Effect = effect
		return []history2.ParticipantEffect{source}
	} else {
		return []history2.ParticipantEffect{{
			AccountID: balance.AccountID,
			BalanceID: &balance.BalanceID,
			AssetCode: &balance.AssetCode,
			Effect:    effect,
		}, source}
	}
}

func (h *reviewRequestOpHandler) getReviewableRequestByID(id int64) *xdr.ReviewableRequestEntry {
	ledgerChanges := h.ledgerChangesProvider.GetLedgerChanges()

	tryFindUpdated := false
	var bestResult *xdr.ReviewableRequestEntry

	for _, change := range ledgerChanges {
		var reviewableRequest *xdr.ReviewableRequestEntry

		switch change.Type {
		case xdr.LedgerEntryChangeTypeCreated:
			reviewableRequest = change.MustCreated().Data.ReviewableRequest
		case xdr.LedgerEntryChangeTypeUpdated:
			tryFindUpdated = false
			reviewableRequest = change.MustUpdated().Data.ReviewableRequest
		case xdr.LedgerEntryChangeTypeState:
			tryFindUpdated = true
			reviewableRequest = change.MustState().Data.ReviewableRequest
		default:
			continue
		}

		if reviewableRequest == nil {
			continue
		}

		if int64(reviewableRequest.RequestId) == id {
			if !tryFindUpdated {
				return reviewableRequest
			}

			bestResult = reviewableRequest
		}
	}

	return bestResult
}
