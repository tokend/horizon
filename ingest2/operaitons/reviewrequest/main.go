package reviewrequest

import (
	"encoding/hex"

	"gitlab.com/tokend/go/amount"

	"gitlab.com/tokend/go/xdr"
	"gitlab.com/tokend/horizon/db2/history2"
)

type reviewRequestOpHandler struct {
	pubKeyProvider     publicKeyProvider // new interfaces
	balanceProvider    balanceProvider
	allRequestHandlers map[xdr.ReviewableRequestType]requestHandlerI
}

type publicKeyProvider interface {
	// GetAccountID returns int value which corresponds to xdr.AccountId
	GetAccountID(raw xdr.AccountId) int64
	// GetBalanceID returns int value which corresponds to xdr.BalanceId
	GetBalanceID(raw xdr.BalanceId) int64
}

type balanceProvider interface {
	// GetBalanceByID returns history balance struct for specific balance id
	GetBalanceByID(balanceID xdr.BalanceId) history2.Balance
}

func NewReviewRequestOpHandler(pubKeyProvider publicKeyProvider, balanceProvider balanceProvider) *reviewRequestOpHandler {
	return &reviewRequestOpHandler{
		pubKeyProvider:     pubKeyProvider,
		balanceProvider:    balanceProvider,
		allRequestHandlers: initializeReviewableRequestMap(balanceProvider, pubKeyProvider),
	}
}

func initializeReviewableRequestMap(balanceProvider balanceProvider, pubKeyProvider publicKeyProvider,
) map[xdr.ReviewableRequestType]requestHandlerI {
	effectHelper := effectHelper{
		balanceProvider: balanceProvider,
	}

	return map[xdr.ReviewableRequestType]requestHandlerI{
		xdr.ReviewableRequestTypeIssuanceCreate: &issuanceHandler{
			effectHelper: effectHelper,
		},
		xdr.ReviewableRequestTypeWithdraw: &withdrawHandler{
			effectHelper: effectHelper,
		},
		xdr.ReviewableRequestTypeAmlAlert: &amlAlertHandler{
			effectHelper: effectHelper,
		},
		xdr.ReviewableRequestTypeInvoice: &invoiceHandler{
			paymentHelper: paymentHelper{
				pubKeyProvider: pubKeyProvider,
			},
		},
		xdr.ReviewableRequestTypeAtomicSwap: &atomicSwapHandler{
			pubKeyProvider: pubKeyProvider,
		},
	}
}

// OperationDetails returns details about review request operation
func (h *reviewRequestOpHandler) OperationDetails(op rawOperation, opRes xdr.OperationResultTr,
) (history2.OperationDetails, error) {
	reviewRequestOp := op.Body.MustReviewRequestOp()
	reviewRequestOpRes := opRes.MustReviewRequestResult().MustSuccess()

	opDetails := history2.OperationDetails{
		Type: xdr.OperationTypeReviewRequest,
		ReviewRequest: &history2.ReviewRequestDetails{
			RequestID:      int64(reviewRequestOp.RequestId),
			RequestType:    reviewRequestOp.RequestDetails.RequestType,
			RequestHash:    hex.EncodeToString(reviewRequestOp.RequestHash[:]),
			Action:         reviewRequestOp.Action,
			Reason:         string(reviewRequestOp.Reason),
			IsFulfilled:    reviewRequestOpRes.Fulfilled,
			RequestDetails: reviewRequestOp.RequestDetails,
		},
	}

	aSwapExtended, ok := reviewRequestOpRes.TypeExt.GetASwapExtended()
	if !ok {
		return opDetails, nil
	}

	opDetails.ReviewRequest.AtomicSwapDetails = &aSwapExtended

	return opDetails, nil
}

// ParticipantsEffects can return different effects depended on request type
func (h *reviewRequestOpHandler) ParticipantsEffects(opBody xdr.OperationBody,
	opRes xdr.OperationResultTr, source history2.ParticipantEffect,
) ([]history2.ParticipantEffect, error) {
	reviewRequestOp := opBody.MustReviewRequestOp()
	reviewRequestRes := opRes.MustReviewRequestResult().MustSuccess()

	if !reviewRequestRes.Fulfilled {
		return []history2.ParticipantEffect{source}, nil
	}

	request := h.getReviewableRequestByID(int64(reviewRequestOp.RequestId))

	if request == nil {
		return []history2.ParticipantEffect{source}, nil
	}

	specificHandler, ok := h.allRequestHandlers[request.Body.Type]
	if !ok {
		return []history2.ParticipantEffect{source}, nil
	}

	return specificHandler.SpecificParticipantsEffects(reviewRequestOp,
		reviewRequestRes, *request, source)
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
	}

	return []history2.ParticipantEffect{{
		AccountID: balance.AccountID,
		BalanceID: &balance.ID,
		AssetCode: &balance.AssetCode,
		Effect:    effect,
	}, source}
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

type requestHandlerI interface {
	SpecificParticipantsEffects(op xdr.ReviewRequestOp, res xdr.ReviewRequestSuccessResult,
		request xdr.ReviewableRequestEntry, source history2.ParticipantEffect,
	) ([]history2.ParticipantEffect, error)
}

type issuanceHandler struct {
	effectHelper effectHelper
}

func (h *issuanceHandler) SpecificParticipantsEffects(op xdr.ReviewRequestOp,
	res xdr.ReviewRequestSuccessResult, request xdr.ReviewableRequestEntry, source history2.ParticipantEffect,
) ([]history2.ParticipantEffect, error) {
	if op.Action != xdr.ReviewRequestOpActionApprove {
		return []history2.ParticipantEffect{source}, nil
	}

	details := request.Body.MustIssuanceRequest()

	effect := history2.Effect{
		Type: history2.EffectTypeFunded,
		Funded: &history2.FundedEffect{
			Amount: amount.StringU(uint64(details.Amount)),
			FeePaid: history2.FeePaid{
				Fixed:             amount.StringU(uint64(details.Fee.Fixed)),
				CalculatedPercent: amount.StringU(uint64(details.Fee.Percent)),
			},
		},
	}

	return h.effectHelper.getParticipantEffectByBalanceID(details.Receiver, effect, source), nil
}

type withdrawHandler struct {
	effectHelper effectHelper
}

func (h *withdrawHandler) SpecificParticipantsEffects(op xdr.ReviewRequestOp,
	res xdr.ReviewRequestSuccessResult, request xdr.ReviewableRequestEntry, source history2.ParticipantEffect,
) ([]history2.ParticipantEffect, error) {
	details := request.Body.MustWithdrawalRequest()

	effect := history2.Effect{
		Type: history2.EffectTypeChargedFromLocked,
		ChargedFromLocked: &history2.ChargedFromLockedEffect{
			Amount: amount.StringU(uint64(details.Amount)),
			FeePaid: history2.FeePaid{
				Fixed:             amount.StringU(uint64(details.Fee.Fixed)),
				CalculatedPercent: amount.StringU(uint64(details.Fee.Percent)),
			},
		},
	}

	if op.Action != xdr.ReviewRequestOpActionApprove {
		effect = history2.Effect{
			Type: history2.EffectTypeUnlocked,
			Unlocked: &history2.UnlockedEffect{
				Amount: amount.StringU(uint64(details.Amount)),
				FeeUnlocked: history2.FeePaid{
					Fixed:             amount.StringU(uint64(details.Fee.Fixed)),
					CalculatedPercent: amount.StringU(uint64(details.Fee.Percent)),
				},
			},
		}
	}

	return h.effectHelper.getParticipantEffectByBalanceID(details.Balance, effect, source), nil
}

type amlAlertHandler struct {
	effectHelper effectHelper
}

func (h *amlAlertHandler) SpecificParticipantsEffects(op xdr.ReviewRequestOp,
	res xdr.ReviewRequestSuccessResult, request xdr.ReviewableRequestEntry, source history2.ParticipantEffect,
) ([]history2.ParticipantEffect, error) {
	details := request.Body.MustAmlAlertRequest()

	effect := history2.Effect{
		Type: history2.EffectTypeChargedFromLocked,
		ChargedFromLocked: &history2.ChargedFromLockedEffect{
			Amount: amount.StringU(uint64(details.Amount)),
		},
	}

	if op.Action != xdr.ReviewRequestOpActionApprove {
		effect = history2.Effect{
			Type: history2.EffectTypeUnlocked,
			Unlocked: &history2.UnlockedEffect{
				Amount: amount.StringU(uint64(details.Amount)),
			},
		}
	}

	return h.effectHelper.getParticipantEffectByBalanceID(details.BalanceId, effect, source), nil
}

type invoiceHandler struct {
	paymentHelper paymentHelper
}

func (h *invoiceHandler) SpecificParticipantsEffects(op xdr.ReviewRequestOp,
	res xdr.ReviewRequestSuccessResult, request xdr.ReviewableRequestEntry, source history2.ParticipantEffect,
) ([]history2.ParticipantEffect, error) {
	paymentOp := op.RequestDetails.MustBillPay().PaymentDetails
	paymentRes := res.TypeExt.MustInvoiceExtended().PaymentV2Response

	effect := history2.EffectTypeFunded
	if request.Body.MustInvoiceRequest().ContractId != nil {
		effect = history2.EffectTypeFundedToLocked
	}

	return h.paymentHelper.getParticipantsEffects(paymentOp, paymentRes, source, effect)
}

type atomicSwapHandler struct {
	pubKeyProvider        publicKeyProvider
	ledgerChangesProvider ledgerChangesProvider
}

func (h *atomicSwapHandler) getAtomicSwapBid(bidID xdr.Uint64) *xdr.AtomicSwapBidEntry {
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

func (h *atomicSwapHandler) SpecificParticipantsEffects(op xdr.ReviewRequestOp,
	res xdr.ReviewRequestSuccessResult, request xdr.ReviewableRequestEntry, source history2.ParticipantEffect,
) ([]history2.ParticipantEffect, error) {
	atomicSwapExtendedResult := res.TypeExt.MustASwapExtended()

	ownerBalanceID := h.pubKeyProvider.GetBalanceID(atomicSwapExtendedResult.BidOwnerBaseBalanceId)

	participants := []history2.ParticipantEffect{{
		AccountID: h.pubKeyProvider.GetAccountID(atomicSwapExtendedResult.BidOwnerId),
		BalanceID: &ownerBalanceID,
		AssetCode: &atomicSwapExtendedResult.BaseAsset,
		Effect: history2.Effect{
			Type: history2.EffectTypeChargedFromLocked,
			ChargedFromLocked: &history2.ChargedFromLockedEffect{
				Amount: amount.StringU(uint64(atomicSwapExtendedResult.BaseAmount)),
			},
		},
	}}

	purchaserBaseBalanceID := h.pubKeyProvider.GetBalanceID(atomicSwapExtendedResult.PurchaserBaseBalanceId)

	participants = append(participants, history2.ParticipantEffect{
		AccountID: h.pubKeyProvider.GetAccountID(atomicSwapExtendedResult.PurchaserId),
		BalanceID: &purchaserBaseBalanceID,
		AssetCode: &atomicSwapExtendedResult.BaseAsset,
		Effect: history2.Effect{
			Type: history2.EffectTypeFunded,
			Funded: &history2.FundedEffect{
				Amount: amount.StringU(uint64(atomicSwapExtendedResult.BaseAmount)),
			},
		},
	})

	bid := h.getAtomicSwapBid(atomicSwapExtendedResult.BidId)
	if bid == nil {
		return participants, nil
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
				Unlocked: &history2.UnlockedEffect{
					Amount: amount.StringU(uint64(bid.Amount)),
				},
			},
		})
	}

	return participants, nil
}
