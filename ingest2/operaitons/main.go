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
	partEffectIDProvider paricipantEffectIDProvider
	pubKeyProvider       publicKeyProvider
}

func newOperationHandler(mainProvider providerCluster) operationHandler {
	pubKeyProvider := mainProvider.GetPubKeyProvider()
	return operationHandler{
		allHandlers: map[xdr.OperationType]operationHandlerI{
			xdr.OperationTypeCreateAccount: &createAccountOpHandler{
				pubKeyProvider: pubKeyProvider,
			},
			xdr.OperationTypeManageAccount: &manageAccountOpHandler{
				pubKeyProvider: pubKeyProvider,
			},
			xdr.OperationTypeManageBalance: &manageBalanceOpHandler{
				pubKeyProvider: pubKeyProvider,
			},
			xdr.OperationTypeManageKeyValue: &manageKeyValueOpHandler{},
			xdr.OperationTypeManageLimits: &manageLimitsOpHandler{
				pubKeyProvider: pubKeyProvider,
			},
			xdr.OperationTypeManageAsset: &manageAssetOpHandler{
				pubKeyProvider: pubKeyProvider,
			},
			xdr.OperationTypeManageAssetPair: &manageAssetPairOpHadler{},
			xdr.OperationTypeManageOffer: &manageOfferOpHandler{
				pubKeyProvider: pubKeyProvider,
			},
			xdr.OperationTypeManageContract: &manageContractOpHandler{},
			xdr.OperationTypeSetFees: &setFeeOpHandler{
				pubKeyProvider: pubKeyProvider,
			},
			xdr.OperationTypeCreatePreissuanceRequest: &createPreIssuanceRequestOpHandler{},
			xdr.OperationTypeCreateIssuanceRequest: &createIssuanceRequestOpHandler{
				pubKeyProvider: pubKeyProvider,
			},
			xdr.OperationTypeCreateWithdrawalRequest: &createWithdrawRequestOpHandler{
				pubKeyProvider: pubKeyProvider,
			},
			xdr.OperationTypeCreateManageLimitsRequest: &createManageLimitsRequestOpHandler{},
			xdr.OperationTypeManageInvoiceRequest: &manageInvoiceRequestOpHandler{
				pubKeyProvider: pubKeyProvider,
			},
			xdr.OperationTypeManageContractRequest: &manageContractRequestOpHandler{
				pubKeyProvider: pubKeyProvider,
			},
			xdr.OperationTypeReviewRequest: &reviewRequestOpHandler{
				pubKeyProvider: pubKeyProvider,
			},
			//xdr.OperationTypePayment: &paymentOpHandler{},
		},
		opIDProvider:         mainProvider.GetOperationIDProvider(),
		partEffectIDProvider: mainProvider.GetParticipantEffectIDProvider(),
		pubKeyProvider:       pubKeyProvider,
	}
}

func (h *operationHandler) ConvertOperation(op xdr.Operation, opRes xdr.OperationResultTr, txSource xdr.AccountId) (history2.Operation, []history2.ParticipantEffect, error) {
	if op.Body.Type != opRes.Type {
		errors.New("operation type mismatch")
	}

	handler := h.allHandlers[op.Body.Type]
	if handler == nil {
		return history2.Operation{}, nil, errors.New("no handler for such operation type")
	}

	details, err := handler.OperationDetails(op.Body, opRes)
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

func (h *operationHandler) getBaseSourceParticipantEffect(opSource *xdr.AccountId,
	txSource xdr.AccountId,
) history2.ParticipantEffect {
	source := txSource
	if opSource != nil {
		source = *opSource
	}

	return history2.ParticipantEffect{
		AccountID: h.pubKeyProvider.GetAccountID(source),
	}
}

type providerCluster interface {
	GetOperationIDProvider() operationIDProvider
	GetParticipantEffectIDProvider() paricipantEffectIDProvider
	GetPubKeyProvider() publicKeyProvider
}

type operationIDProvider interface {
	GetOperationID() int64
}

type paricipantEffectIDProvider interface {
	GetNextParticipantEffectID() int64
}

type publicKeyProvider interface {
	GetAccountID(raw xdr.AccountId) int64
	GetBalanceID(raw xdr.BalanceId) int64
}

type operationHandlerI interface {
	OperationDetails(opBody xdr.OperationBody, opRes xdr.OperationResultTr) (history2.OperationDetails, error)
	ParticipantsEffects(opBody xdr.OperationBody, opRes xdr.OperationResultTr, source history2.ParticipantEffect) ([]history2.ParticipantEffect, error)
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

// TODO set option operation handler

type paymentOpHandler struct {
}

func (h *paymentOpHandler) OperationDetails(opBody xdr.OperationBody) (history2.OperationDetails, error) {
	op := opBody.MustCreateAccountOp()

	return history2.OperationDetails{
		Type:    xdr.OperationTypeCreateAccount,
		Payment: &history2.PaymentDetails{},
	}, nil
}

func (h *paymentOpHandler) GetParticipantsEffects(opBody xdr.OperationBody) ([]history2.ParticipantEffect, error) {
	var participants []history2.ParticipantEffect
	var converter history2.PubKeyConverter

	op := opBody.MustPaymentOp()
	participants = append(participants, history2.ParticipantEffect{
		AccountID:   converter.ConvertToInt64(xdr.PublicKey(op.Destination)),
		OperationID: 0,
	})

	if op.Referrer != nil {
		participants = append(participants, history2.ParticipantEffect{
			AccountID: converter.ConvertToInt64(xdr.PublicKey(*op.Referrer)),
		})
	}

	return participants, nil
}
