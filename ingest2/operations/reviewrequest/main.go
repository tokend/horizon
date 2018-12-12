package reviewrequest

import (
	"encoding/hex"

	"gitlab.com/distributed_lab/logan/v3/errors"

	"gitlab.com/tokend/horizon/ingest2/operations"

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

// NewReviewRequestOpHandler creates new handler for review request operations
// with specific handlers for different types
func NewReviewRequestOpHandler(pubKeyProvider publicKeyProvider, balanceProvider balanceProvider,
) operations.OperationHandler {
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
			paymentHelper: operations.NewPaymentHelper(pubKeyProvider),
		},
		xdr.ReviewableRequestTypeAtomicSwap: &atomicSwapHandler{
			pubKeyProvider: pubKeyProvider,
		},
	}
}

// OperationDetails returns details about review request operation
func (h *reviewRequestOpHandler) OperationDetails(op operations.RawOperation, opRes xdr.OperationResultTr,
) (history2.OperationDetails, error) {
	reviewRequestOp := op.Body.MustReviewRequestOp()
	reviewRequestOpRes := opRes.MustReviewRequestResult().MustSuccess()

	var addedTasks uint32
	var removedTasks uint32
	var externalDetails string
	if details, ok := reviewRequestOp.Ext.GetReviewDetails(); ok {
		addedTasks = uint32(details.TasksToAdd)
		removedTasks = uint32(details.TasksToRemove)
		externalDetails = string(details.ExternalDetails)
	}

	opDetails := history2.OperationDetails{
		Type: xdr.OperationTypeReviewRequest,
		ReviewRequest: &history2.ReviewRequestDetails{
			RequestID:       int64(reviewRequestOp.RequestId),
			RequestType:     reviewRequestOp.RequestDetails.RequestType,
			RequestHash:     hex.EncodeToString(reviewRequestOp.RequestHash[:]),
			Action:          reviewRequestOp.Action,
			Reason:          string(reviewRequestOp.Reason),
			IsFulfilled:     reviewRequestOpRes.Fulfilled,
			RequestDetails:  reviewRequestOp.RequestDetails,
			AddedTasks:      addedTasks,
			RemovedTasks:    removedTasks,
			ExternalDetails: externalDetails,
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
	opRes xdr.OperationResultTr, source history2.ParticipantEffect, ledgerChanges []xdr.LedgerEntryChange,
) ([]history2.ParticipantEffect, error) {
	reviewRequestOp := opBody.MustReviewRequestOp()
	reviewRequestRes := opRes.MustReviewRequestResult().MustSuccess()

	if !reviewRequestRes.Fulfilled {
		return []history2.ParticipantEffect{source}, nil
	}

	request := h.getReviewableRequestByID(int64(reviewRequestOp.RequestId), ledgerChanges)

	if request == nil {
		return nil, errors.From(
			errors.New("expected request to be in ledger changes"), map[string]interface{}{
				"request_id": int64(reviewRequestOp.RequestId),
			})
	}

	specificHandler, ok := h.allRequestHandlers[request.Body.Type]
	if !ok {
		return []history2.ParticipantEffect{source}, nil
	}

	return specificHandler.specificParticipantsEffects(reviewRequestOp,
		reviewRequestRes, *request, source, ledgerChanges)
}

type effectHelper struct {
	balanceProvider balanceProvider
}

func (h *effectHelper) populateEffectWithBalanceDetails(balanceID xdr.BalanceId,
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

// Tries to get latest version of reviewable request by ID (updated first, created or state otherwise)
func (h *reviewRequestOpHandler) getReviewableRequestByID(id int64, ledgerChanges []xdr.LedgerEntryChange,
) *xdr.ReviewableRequestEntry {
	var bestResult *xdr.ReviewableRequestEntry

	for _, change := range ledgerChanges {
		var ledgerEntryData xdr.LedgerEntryData

		switch change.Type {
		case xdr.LedgerEntryChangeTypeCreated:
			ledgerEntryData = change.MustCreated().Data
		case xdr.LedgerEntryChangeTypeUpdated:
			ledgerEntryData = change.MustUpdated().Data
		case xdr.LedgerEntryChangeTypeState:
			ledgerEntryData = change.MustState().Data
		default:
			continue
		}

		if ledgerEntryData.Type != xdr.LedgerEntryTypeReviewableRequest {
			continue
		}

		request := ledgerEntryData.MustReviewableRequest()
		if int64(request.RequestId) != id {
			continue
		}

		bestResult = &request
		// we have found latest version of the request, so can return
		if change.Type == xdr.LedgerEntryChangeTypeUpdated {
			return bestResult
		}
	}

	return bestResult
}

type requestHandlerI interface {
	specificParticipantsEffects(op xdr.ReviewRequestOp, res xdr.ReviewRequestSuccessResult,
		request xdr.ReviewableRequestEntry, source history2.ParticipantEffect, ledgerChanges []xdr.LedgerEntryChange,
	) ([]history2.ParticipantEffect, error)
}
