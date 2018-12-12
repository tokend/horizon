package operations

import (
	"gitlab.com/distributed_lab/logan/v3/errors"
	"gitlab.com/tokend/go/xdr"
	"gitlab.com/tokend/horizon/db2/history2"
	"gitlab.com/tokend/horizon/ingest2/operations/reviewrequest"
)

type operationHandler struct {
	allHandlers          map[xdr.OperationType]OperationHandler
	opIDProvider         operationIDProvider
	partEffectIDProvider participantEffectIDProvider
	pubKeyProvider       publicKeyProvider
}

// newOperationHandler returns new handler which can return
// details and participants effects of certain operation
func newOperationHandler(mainProvider providerCluster) operationHandler {
	pubKeyProvider := mainProvider.GetPubKeyProvider()
	offerHelper := offerHelper{
		pubKeyProvider: pubKeyProvider,
	}
	balanceProvider := mainProvider.GetBalanceProvider()
	return operationHandler{
		allHandlers: map[xdr.OperationType]OperationHandler{
			xdr.OperationTypeCreateAccount: &createAccountOpHandler{
				pubKeyProvider: pubKeyProvider,
			},
			xdr.OperationTypeManageAccount: &manageAccountOpHandler{
				pubKeyProvider: pubKeyProvider,
			},
			xdr.OperationTypeManageExternalSystemAccountIdPoolEntry: &manageExternalSystemPoolOpHandler{},
			xdr.OperationTypeBindExternalSystemAccountId:            &bindExternalSystemAccountOpHandler{},
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
				pubKeyProvider:  pubKeyProvider,
				offerHelper:     offerHelper,
				balanceProvider: balanceProvider,
			},
			xdr.OperationTypeManageContract: &manageContractOpHandler{
				pubKeyProvider:  pubKeyProvider,
				requestProvider: mainProvider.GetRequestProvider(),
			},
			xdr.OperationTypeSetFees: &setFeeOpHandler{
				pubKeyProvider: pubKeyProvider,
			},
			xdr.OperationTypeCreateKycRequest: &createKYCRequestOpHandler{
				pubKeyProvider: pubKeyProvider,
			},
			xdr.OperationTypeCreatePreissuanceRequest: &createPreIssuanceRequestOpHandler{},
			xdr.OperationTypeCreateIssuanceRequest: &createIssuanceRequestOpHandler{
				pubKeyProvider: pubKeyProvider,
			},
			xdr.OperationTypeCreateSaleRequest: &createSaleRequestOpHandler{},
			xdr.OperationTypeCreateAswapBidRequest: &createAtomicSwapBidRequestOpHandler{
				balanceProvider: balanceProvider,
			},
			xdr.OperationTypeCreateAswapRequest: &createAtomicSwapRequestOpHandler{},
			xdr.OperationTypeCreateWithdrawalRequest: &createWithdrawRequestOpHandler{
				pubKeyProvider: pubKeyProvider,
			},
			xdr.OperationTypeCreateAmlAlert: &createAMLAlertReqeustOpHandler{
				balanceProvider: balanceProvider,
			},
			xdr.OperationTypeCreateManageLimitsRequest: &createManageLimitsRequestOpHandler{},
			xdr.OperationTypeManageInvoiceRequest:      &manageInvoiceRequestOpHandler{},
			xdr.OperationTypeManageContractRequest: &manageContractRequestOpHandler{
				pubKeyProvider: pubKeyProvider,
			},
			xdr.OperationTypeReviewRequest: reviewrequest.NewReviewRequestOpHandler(pubKeyProvider, balanceProvider),
			xdr.OperationTypePaymentV2: &paymentOpHandler{
				pubKeyProvider: pubKeyProvider,
				paymentHelper: PaymentHelper{
					pubKeyProvider: pubKeyProvider,
				},
			},
			xdr.OperationTypeCheckSaleState: &checkSaleStateOpHandler{
				pubKeyProvider:  pubKeyProvider,
				offerHelper:     offerHelper,
				balanceProvider: balanceProvider,
			},
			xdr.OperationTypeCancelAswapBid: &cancelAtomicSwapBidOpHandler{
				pubKeyProvider: pubKeyProvider,
			},
		},
		opIDProvider:         mainProvider.GetOperationIDProvider(),
		partEffectIDProvider: mainProvider.GetParticipantEffectIDProvider(),
		pubKeyProvider:       pubKeyProvider,
	}
}

// ConvertOperation transforms xdr operation data to db suitable Operation and Participants Effects
func (h *operationHandler) ConvertOperation(op xdr.Operation, opRes xdr.OperationResultTr,
	txSource xdr.AccountId, ledgerChanges []xdr.LedgerEntryChange,
) (history2.Operation, []history2.ParticipantEffect, error) {
	handler, ok := h.allHandlers[op.Body.Type]
	if !ok {
		return history2.Operation{}, nil, errors.From(
			errors.New("no handler for such operation type"), map[string]interface{}{
				"operation type": op.Body.Type.String(),
			})
	}

	source := h.getOperationSource(op.SourceAccount, txSource)
	details, err := handler.OperationDetails(RawOperation{
		Body:   op.Body,
		Source: source,
	}, opRes)
	if err != nil {
		return history2.Operation{}, nil,
			errors.Wrap(err, "failed to get operation details", map[string]interface{}{
				"operation type": int32(op.Body.Type),
			})
	}

	participantsEffects, err := handler.ParticipantsEffects(op.Body, opRes,
		h.getBaseSourceParticipantEffect(op.SourceAccount, txSource), ledgerChanges)
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
		Source:           source.Address(),
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
