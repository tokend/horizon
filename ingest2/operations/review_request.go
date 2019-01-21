package operations

import (
	"encoding/hex"

	"gitlab.com/distributed_lab/logan/v3"
	"gitlab.com/distributed_lab/logan/v3/errors"
	"gitlab.com/tokend/go/xdr"
	"gitlab.com/tokend/horizon/db2/history2"
	regources "gitlab.com/tokend/regources/v2"
)

type reviewRequestOpHandler struct {
	pubKeyProvider     IDProvider
	balanceProvider    balanceProvider
	allRequestHandlers map[xdr.ReviewableRequestType]reviewRequestHandler
}

type reviewRequestHandler interface {
	ParticipantsEffects(op xdr.ReviewRequestOp, res xdr.ExtendedResult,
		request xdr.ReviewableRequestEntry, source history2.ParticipantEffect, ledgerChanges []xdr.LedgerEntryChange,
	) ([]history2.ParticipantEffect, error)
}

// newReviewRequestOpHandler creates new handler for review request operations
// with specific handlers for different types
func newReviewRequestOpHandler(pubKeyProvider IDProvider, balanceProvider balanceProvider,
) *reviewRequestOpHandler {
	// All reviewable requests must be explicitly handled here. If there is nothing to handle
	// use reviewableRequestHandlerStub
	return &reviewRequestOpHandler{
		pubKeyProvider:  pubKeyProvider,
		balanceProvider: balanceProvider,
		allRequestHandlers: map[xdr.ReviewableRequestType]reviewRequestHandler{
			xdr.ReviewableRequestTypeIssuanceCreate: &issuanceHandler{
				balanceProvider: balanceProvider,
			},
			xdr.ReviewableRequestTypeWithdraw: &withdrawHandler{
				balanceProvider: balanceProvider,
			},
			xdr.ReviewableRequestTypeAmlAlert: &amlAlertHandler{
				balanceProvider: balanceProvider,
			},
			xdr.ReviewableRequestTypeAtomicSwap: &atomicSwapHandler{
				pubKeyProvider: pubKeyProvider,
			},
			xdr.ReviewableRequestTypeAssetCreate:         &reviewableRequestHandlerStub{},
			xdr.ReviewableRequestTypeAssetUpdate:         &reviewableRequestHandlerStub{},
			xdr.ReviewableRequestTypePreIssuanceCreate:   &reviewableRequestHandlerStub{},
			xdr.ReviewableRequestTypeSale:                &reviewableRequestHandlerStub{},
			xdr.ReviewableRequestTypeLimitsUpdate:        &reviewableRequestHandlerStub{},
			xdr.ReviewableRequestTypeUpdateKyc:           &reviewableRequestHandlerStub{},
			xdr.ReviewableRequestTypeUpdateSaleDetails:   &reviewableRequestHandlerStub{},
			xdr.ReviewableRequestTypeInvoice:             &deprecatedReviewRequestHandler{},
			xdr.ReviewableRequestTypeContract:            &deprecatedReviewRequestHandler{},
			xdr.ReviewableRequestTypeCreateAtomicSwapBid: &reviewableRequestHandlerStub{},
		},
	}
}

// Details returns details about review request operation
func (h *reviewRequestOpHandler) Details(op rawOperation, opRes xdr.OperationResultTr,
) (regources.OperationDetails, error) {
	reviewRequestOp := op.Body.MustReviewRequestOp()
	reviewRequestOpRes := opRes.MustReviewRequestResult().MustSuccess()

	addedTasks := uint32(reviewRequestOp.ReviewDetails.TasksToAdd)
	removedTasks := uint32(reviewRequestOp.ReviewDetails.TasksToRemove)
	externalDetails := []byte(reviewRequestOp.ReviewDetails.ExternalDetails)

	opDetails := regources.OperationDetails{
		Type: xdr.OperationTypeReviewRequest,
		ReviewRequest: &regources.ReviewRequestDetails{
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

	return opDetails, nil
}

// ParticipantsEffects - returns source participant if request was not fulfilled
// finds specific handler otherwise
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
		return nil, errors.From(errors.New("failed to find handler for reviewable request"), logan.F{
			"request_type": request.Body.Type,
		})
	}

	return specificHandler.ParticipantsEffects(reviewRequestOp,
		reviewRequestRes, *request, source, ledgerChanges)
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
