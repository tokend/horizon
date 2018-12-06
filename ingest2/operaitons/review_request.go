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
	allRequestHandlers    map[xdr.ReviewableRequestType]requestHandlerI
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

// TODO handle participants effects based on operation result

func (h *reviewRequestOpHandler) ParticipantsEffects(opBody xdr.OperationBody,
	opRes xdr.OperationResultTr, source history2.ParticipantEffect,
) ([]history2.ParticipantEffect, error) {
	reviewRequestOp := opBody.MustReviewRequestOp()
	reviewRequestRes := opRes.MustReviewRequestResult().MustSuccess()

	request := h.getReviewableRequestByID(int64(reviewRequestOp.RequestId))

	if request == nil {
		return []history2.ParticipantEffect{source}, nil
	}

	specificHandler, ok := h.allRequestHandlers[request.Body.Type]
	if !ok {
		return []history2.ParticipantEffect{source}, nil
	}

	return specificHandler.SpecificParticipantsEffects(reviewRequestOp,
		reviewRequestRes, *request, source), nil

	// maybe do map with specific handlers
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
	case xdr.ReviewableRequestTypeAtomicSwap:
		atomicSwapExtendedResult := reviewRequestRes.Ext.MustExtendedResult().TypeExt.MustASwapExtended()
		if !reviewRequestRes.Ext.MustExtendedResult().Fulfilled {
			break
		}

		ownerBalanceID := h.pubKeyProvider.GetBalanceID(atomicSwapExtendedResult.BidOwnerBaseBalanceId)

		participants = append(participants, history2.ParticipantEffect{
			AccountID: h.pubKeyProvider.GetAccountID(atomicSwapExtendedResult.BidOwnerId),
			BalanceID: &ownerBalanceID,
			AssetCode: &atomicSwapExtendedResult.BaseAsset,
			Effect: history2.Effect{
				Type: history2.EffectTypeChargedFromLocked,
				AtomicSwap: &history2.AtomicSwapEffect{
					Amount: amount.StringU(uint64(atomicSwapExtendedResult.BaseAmount)),
				},
			},
		})

		purchaserBaseBalanceId := h.pubKeyProvider.GetBalanceID(atomicSwapExtendedResult.PurchaserBaseBalanceId)

		participants = append(participants, history2.ParticipantEffect{
			AccountID: h.pubKeyProvider.GetAccountID(atomicSwapExtendedResult.PurchaserId),
			BalanceID: &purchaserBaseBalanceId,
			AssetCode: &atomicSwapExtendedResult.BaseAsset,
			Effect: history2.Effect{
				Type: history2.EffectTypeFunded,
				AtomicSwap: &history2.AtomicSwapEffect{
					Amount: amount.StringU(uint64(atomicSwapExtendedResult.BaseAmount)),
				},
			},
		})

		bid := h.getAtomicSwapBid(atomicSwapExtendedResult.BidId)
		if bid == nil {
			break
		}

		bidIsSoldOut := (bid.Amount == 0) && (bid.LockedAmount == 0)
		bidIsRemoved := bidIsSoldOut || (bid.IsCancelled && bid.LockedAmount == 0)
		if bidIsRemoved && (bid.Amount != 0) {
			participants = append(participants, history2.ParticipantEffect{
				AccountID: h.pubKeyProvider.GetAccountID(atomicSwapExtendedResult.BidOwnerId),
				BalanceID: &ownerBalanceID,
				AssetCode: &atomicSwapExtendedResult.BaseAsset,
				Effect: history2.Effect{
					Type: history2.EffectTypeUnlocked,
					AtomicSwap: &history2.AtomicSwapEffect{
						Amount: amount.StringU(uint64(bid.Amount)),
					},
				},
			})
		}
	}
}

type effectHelper struct {
	balanceProvider balanceProvider
}

func (h *effectHelper) getParticipantEffectByBalanceID(balanceID xdr.BalanceId,
	effect history2.Effect, source history2.ParticipantEffect,
) []history2.ParticipantEffect {
	balance := h.balanceProvider.GetBalanceByID(balanceID)
	if balance.AccountID == source.AccountID {
		source.BalanceID = &balance.ID
		source.AssetCode = &balance.AssetCode
		source.Effect = effect
		return []history2.ParticipantEffect{source}
	} else {
		return []history2.ParticipantEffect{{
			AccountID: balance.AccountID,
			BalanceID: &balance.ID,
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

func (h *reviewRequestOpHandler) getAtomicSwapBid(bidID xdr.Uint64) *xdr.AtomicSwapBidEntry {
	ledgerChanges := h.ledgerChangesProvider.GetLedgerChanges()

	for _, change := range ledgerChanges {
		if change.Type != xdr.LedgerEntryChangeTypeUpdated {
			continue
		}

		if change.MustUpdated().Data.Type != xdr.LedgerEntryTypeAtomicSwapBid {
			continue
		}

		atomicSwapBid := change.MustUpdated().Data.MustAtomicSwapBid()

		if atomicSwapBid.BidId == bidID {
			return &atomicSwapBid
		}
	}

	return nil
}

type requestHandlerI interface {
	SpecificParticipantsEffects(op xdr.ReviewRequestOp, res xdr.ReviewRequestResultSuccess,
		request xdr.ReviewableRequestEntry, source history2.ParticipantEffect,
	) []history2.ParticipantEffect
}

type issuanceHandler struct{}

func (h *issuanceHandler) SpecificParticipantsEffects(op xdr.ReviewRequestOp, res xdr.ReviewRequestResultSuccess) []history2.ParticipantEffect {

}

type withdrawHandler struct {
}

func (h *withdrawHandler) SpecificParticipantsEffects(op xdr.ReviewRequestOp, res xdr.ReviewRequestResultSuccess) []history2.ParticipantEffect {

}

type amlAlertHandler struct {
	effectHelper effectHelper
}

func (h *amlAlertHandler) SpecificParticipantsEffects(op xdr.ReviewRequestOp,
	res xdr.ReviewRequestResultSuccess, request xdr.ReviewableRequestEntry, source history2.ParticipantEffect,
) []history2.ParticipantEffect {
	details := request.Body.MustAmlAlertRequest()

	effect := history2.Effect{
		Type: history2.EffectTypeChargedFromLocked,
		ChargedFromLocked: &history2.ChargedFromLockedEffect{
			Amount: amount.StringU(uint64(details.Amount)),
		},
	}

	if op.Action != xdr.ReviewRequestOpActionApprove {
		effect.Type = history2.EffectTypeUnlocked
		effect.ChargedFromLocked = nil
		effect.Unlocked = &history2.UnlockedEffect{
			Amount: amount.StringU(uint64(details.Amount)),
		}
	}

	return h.effectHelper.getParticipantEffectByBalanceID(details.BalanceId, effect, source)
}

type invoiceHandler struct{}

func (h *invoiceHandler) SpecificParticipantsEffects(op xdr.ReviewRequestOp, res xdr.ReviewRequestResultSuccess) []history2.ParticipantEffect {
}

type atomicSwapHandler struct{}

func (h *atomicSwapHandler) SpecificParticipantsEffects(op xdr.ReviewRequestOp, res xdr.ReviewRequestResultSuccess) []history2.ParticipantEffect {
}

/*func (h *reviewRequestOpHandler) getNotApprovedRequestParticipnatsEffects(op xdr.ReviewRequestOp) []history2.ParticipantEffect {
	var participants []history2.ParticipantEffect

	switch op.RequestDetails.RequestType {
	case xdr.ReviewableRequestTypeAmlAlert:
		op.RequestDetails.MustAmlAlertDetails().

		participants = append(participants, history2.ParticipantEffect{
			AccountID: op
		})
	case xdr.ReviewableRequestTypeWithdraw:

	case xdr.ReviewableRequestTypeCreateAtomicSwapBid:
	case xdr.ReviewableRequestTypeAtomicSwap:
	}
}*/
