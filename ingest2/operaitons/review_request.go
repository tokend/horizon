package operaitons

import (
	"encoding/hex"

	"gitlab.com/tokend/go/xdr"
	"gitlab.com/tokend/horizon/db2/history2"
)

type reviewRequestOpHandler struct {
	pubKeyProvider        publicKeyProvider
	ledgerChangesProvider ledgerChangesProvider
}

func (h *reviewRequestOpHandler) OperationDetails(opBody xdr.OperationBody,
	opRes xdr.OperationResultTr,
) (history2.OperationDetails, error) {
	reviewRequestOp := opBody.MustReviewRequestOp()
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

	participants := []history2.ParticipantEffect{source}

	request := h.getReviewableRequestByID(int64(reviewRequestOp.RequestId))
	if request == nil || source.AccountID == h.pubKeyProvider.GetAccountID(request.Requestor) {
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
