package operaitons

import (
	"encoding/json"

	"gitlab.com/distributed_lab/logan/v3/errors"
	"gitlab.com/tokend/go/xdr"
	"gitlab.com/tokend/horizon/db2/history2"
)

type operationHandler struct {
	allHandlers          map[xdr.OperationType]operationHandlerI
	opIDProvider         operationIDProvider
	partEffectIDProvider participantEffectIDProvider
	pubKeyProvider       publicKeyProvider
}

func newOperationHandler(mainProvider providerCluster) operationHandler {
	pubKeyProvider := mainProvider.GetPubKeyProvider()
	offerHelper := offerHelper{
		pubKeyProvider:        pubKeyProvider,
		ledgerChangesProvider: mainProvider.GetLedgerChangesProvider(),
	}
	paymentHelper := paymentHelper{
		pubKeyProvider: pubKeyProvider,
	}
	return operationHandler{
		allHandlers: map[xdr.OperationType]operationHandlerI{
			xdr.OperationTypeCreateAccount: &createAccountOpHandler{
				pubKeyProvider: pubKeyProvider,
			},
			xdr.OperationTypeManageAccount: &manageAccountOpHandler{
				pubKeyProvider: pubKeyProvider,
			},
			xdr.OperationTypeManageExternalSystemAccountIdPoolEntry: &manageExternalSystemPoolOpHandler{},
			xdr.OperationTypeManageBalance: &manageBalanceOpHandler{
				pubKeyProvider: pubKeyProvider,
			},
			xdr.OperationTypeManageKeyValue: &manageKeyValueOpHandler{},
			xdr.OperationTypeManageLimits: &manageLimitsOpHandler{
				pubKeyProvider: pubKeyProvider,
			},
			xdr.OperationTypeManageAsset:     &manageAssetOpHandler{},
			xdr.OperationTypeManageAssetPair: &manageAssetPairOpHadler{},
			xdr.OperationTypeManageOffer: &manageOfferOpHandler{
				pubKeyProvider:        pubKeyProvider,
				offerHelper:           offerHelper,
				ledgerChangesProvider: mainProvider.GetLedgerChangesProvider(),
				balanceProvider:       mainProvider.GetBalanceProvider(),
			},
			xdr.OperationTypeManageContract: &manageContractOpHandler{
				pubKeyProvider:  pubKeyProvider,
				requestProvider: mainProvider.GetRequestProvider(),
			},
			xdr.OperationTypeSetFees: &setFeeOpHandler{
				pubKeyProvider: pubKeyProvider,
			},
			xdr.OperationTypeCreatePreissuanceRequest: &createPreIssuanceRequestOpHandler{},
			xdr.OperationTypeCreateIssuanceRequest: &createIssuanceRequestOpHandler{
				pubKeyProvider: pubKeyProvider,
			},
			xdr.OperationTypeCreateSaleRequest: &createSaleRequestOpHandler{},
			xdr.OperationTypeCreateAswapBidRequest: &createAtomicSwapBidRequestOpHandler{
				balanceProvider: mainProvider.GetBalanceProvider(),
			},
			xdr.OperationTypeCreateAswapRequest: &createAtomicSwapRequestOpHandler{},
			xdr.OperationTypeCreateWithdrawalRequest: &createWithdrawRequestOpHandler{
				pubKeyProvider: pubKeyProvider,
			},
			xdr.OperationTypeCreateAmlAlert: &createAMLAlertReqeustOpHandler{
				balanceProvider: mainProvider.GetBalanceProvider(),
			},
			xdr.OperationTypeCreateManageLimitsRequest: &createManageLimitsRequestOpHandler{},
			xdr.OperationTypeManageInvoiceRequest:      &manageInvoiceRequestOpHandler{},
			xdr.OperationTypeManageContractRequest: &manageContractRequestOpHandler{
				pubKeyProvider: pubKeyProvider,
			},
			xdr.OperationTypeReviewRequest: &reviewRequestOpHandler{
				pubKeyProvider:        pubKeyProvider,
				ledgerChangesProvider: mainProvider.GetLedgerChangesProvider(),
			},
			xdr.OperationTypePaymentV2: &paymentOpHandler{
				pubKeyProvider: pubKeyProvider,
				paymentHelper:  paymentHelper,
			},
			xdr.OperationTypeCheckSaleState: &checkSaleStateOpHandler{
				pubKeyProvider:        pubKeyProvider,
				offerHelper:           offerHelper,
				ledgerChangesProvider: mainProvider.GetLedgerChangesProvider(),
				balanceProvider:       mainProvider.GetBalanceProvider(),
			},
			xdr.OperationTypeCancelAswapBid: &cancelAtomicSwapBidOpHandler{
				pubKeyProvider:        pubKeyProvider,
				ledgerChangesProvider: mainProvider.GetLedgerChangesProvider(),
			},
		},
		opIDProvider:         mainProvider.GetOperationIDProvider(),
		partEffectIDProvider: mainProvider.GetParticipantEffectIDProvider(),
		pubKeyProvider:       pubKeyProvider,
	}
}

// ConvertOperation transform xdr operation data to db suitable Operation and Participants Effects
func (h *operationHandler) ConvertOperation(op xdr.Operation, opRes xdr.OperationResultTr, txSource xdr.AccountId,
) (history2.Operation, []history2.ParticipantEffect, error) {
	handler, ok := h.allHandlers[op.Body.Type]
	if !ok {
		return history2.Operation{}, nil, errors.From(
			errors.New("no handler for such operation type"), map[string]interface{}{
				"operation type": op.Body.Type.String(),
			})
	}

	details, err := handler.OperationDetails(rawOperation{
		Body:   op.Body,
		Source: h.getOperationSource(op.SourceAccount, txSource),
	}, opRes)
	if err != nil {
		return history2.Operation{}, nil,
			errors.Wrap(err, "failed to get operation details", map[string]interface{}{
				"operation type": int32(op.Body.Type),
			})
	}

	participantsEffects, err := handler.ParticipantsEffects(op.Body, opRes,
		h.getBaseSourceParticipantEffect(op.SourceAccount, txSource))
	if err != nil {
		return history2.Operation{}, nil,
			errors.Wrap(err, "failed to get participants effects", map[string]interface{}{
				"operation type": int32(op.Body.Type),
			})
	}

	operationID := h.opIDProvider.GetOperationID()
	for i := range participantsEffects {
		participantsEffects[i].OperationID = operationID
		participantsEffects[i].ID = h.partEffectIDProvider.GetNextParticipantEffectID()
	}

	return history2.Operation{
		ID:               operationID,
		OperationDetails: details,
		Type:             op.Body.Type,
	}, participantsEffects, nil
}

func (h *operationHandler) getOperationSource(opSource *xdr.AccountId,
	txSource xdr.AccountId,
) xdr.AccountId {
	source := txSource
	if opSource != nil {
		source = *opSource
	}

	return source
}

func (h *operationHandler) getBaseSourceParticipantEffect(opSource *xdr.AccountId,
	txSource xdr.AccountId,
) history2.ParticipantEffect {
	return history2.ParticipantEffect{
		AccountID: h.pubKeyProvider.GetAccountID(h.getOperationSource(opSource, txSource)),
	}
}

type providerCluster interface {
	GetOperationIDProvider() operationIDProvider
	GetParticipantEffectIDProvider() participantEffectIDProvider
	GetPubKeyProvider() publicKeyProvider
	GetBalanceProvider() balanceProvider
	GetLedgerChangesProvider() ledgerChangesProvider
	GetRequestProvider() requestProvider
}

type operationIDProvider interface {
	GetOperationID() int64
}

type participantEffectIDProvider interface {
	GetNextParticipantEffectID() int64
}

type publicKeyProvider interface {
	GetAccountID(raw xdr.AccountId) int64
	GetBalanceID(raw xdr.BalanceId) int64
}

type balanceProvider interface {
	GetBalanceByID(balanceID xdr.BalanceId) history2.Balance
}

type ledgerChangesProvider interface {
	GetLedgerChanges() xdr.LedgerEntryChanges
}

type requestProvider interface {
	GetInvoiceRequestsByContractID(contractID int64) []xdr.ReviewableRequestEntry
}

type operationHandlerI interface {
	OperationDetails(op rawOperation, opRes xdr.OperationResultTr) (history2.OperationDetails, error)
	ParticipantsEffects(opBody xdr.OperationBody, opRes xdr.OperationResultTr, source history2.ParticipantEffect) ([]history2.ParticipantEffect, error)
}

type rawOperation struct {
	Source xdr.AccountId
	Body   xdr.OperationBody
}

func customDetailsUnmarshal(rawDetails []byte) map[string]interface{} {
	var result map[string]interface{}
	err := json.Unmarshal([]byte(rawDetails), &result)
	if err != nil {
		result = make(map[string]interface{})
		result["data"] = string(rawDetails)
		result["error"] = err.Error()
	}

	return result
}

func initializeRequeviwableRequestMap() map[xdr.ReviewableRequestType]requestHandlerI {
	return map[xdr.ReviewableRequestType]requestHandlerI{
		xdr.ReviewableRequestTypeIssuanceCreate: &issuanceHandler{},
		xdr.ReviewableRequestTypeWithdraw:       &withdrawHandler{},
		xdr.ReviewableRequestTypeAmlAlert:       &amlAlertHandler{},
		xdr.ReviewableRequestTypeInvoice:        &invoiceHandler{},
		xdr.ReviewableRequestTypeAtomicSwap:     &atomicSwapHandler{},
	}
}

// TODO set option operation handler
