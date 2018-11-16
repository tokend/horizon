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
	pubKeyConverter      publicKeyConverter
}

func newOperationHandler(mainProvider providerCluster) operationHandler {
	return operationHandler{
		allHandlers: map[xdr.OperationType]operationHandlerI{
			xdr.OperationTypeCreateAccount: &createAccountOpHandler{
				pubKeyConverter: mainProvider.GetPubKeyConverter(),
			},
			xdr.OperationTypeManageKeyValue: &manageKeyValueOpHandler{},
			xdr.OperationTypeSetFees: &setFeeOpHandler{
				pubKeyConverter: mainProvider.GetPubKeyConverter(),
			},
			xdr.OperationTypeManageAccount: &manageAccountOpHandler{
				pubKeyConverter: mainProvider.GetPubKeyConverter(),
			},
			xdr.OperationTypeCreateWithdrawalRequest: &createWithdrawRequestOpHandler{
				pubKeyConverter: mainProvider.GetPubKeyConverter(),
			},
			xdr.OperationTypePayment: &paymentOpHandler{},
		},
		opIDProvider:         mainProvider.GetOperationIDProvider(),
		partEffectIDProvider: mainProvider.GetParticipantEffectIDProvider(),
		pubKeyConverter:      mainProvider.GetPubKeyConverter(),
	}
}

func (h *operationHandler) ConvertOperation(op xdr.Operation, opRes xdr.OperationResult, txSource xdr.AccountId) (history2.Operation, []history2.ParticipantEffect, error) {
	if op.Body.Type != opRes.MustTr().Type {
		panic("operation type mismatch")
	}

	handler := h.allHandlers[op.Body.Type]

	details, err := handler.OperationDetails(op.Body, opRes.MustTr())
	if err != nil {
		return history2.Operation{}, nil,
			errors.Wrap(err, "failed to get operation details", map[string]interface{}{
				"operation type": int32(op.Body.Type),
			})
	}

	participantsEffects, err := handler.ParticipantsEffects(op.Body,
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
		AccountID: h.pubKeyConverter.ConvertToInt64(xdr.PublicKey(source)),
	}
}

type providerCluster interface {
	GetOperationIDProvider() operationIDProvider
	GetParticipantEffectIDProvider() paricipantEffectIDProvider
	GetPubKeyConverter() publicKeyConverter
}

type operationIDProvider interface {
	GetOperationID() int64
}

type paricipantEffectIDProvider interface {
	GetNextParticipantEffectID() int64
}

type publicKeyConverter interface {
	ConvertToInt64(key xdr.PublicKey) int64
}

type operationHandlerI interface {
	OperationDetails(opBody xdr.OperationBody, opRes xdr.OperationResultTr) (history2.OperationDetails, error)
	ParticipantsEffects(opBody xdr.OperationBody, source history2.ParticipantEffect) ([]history2.ParticipantEffect, error)
}

type createAccountOpHandler struct {
	pubKeyConverter publicKeyConverter
}

func (h *createAccountOpHandler) OperationDetails(opBody xdr.OperationBody) (history2.OperationDetails, error) {
	op := opBody.MustCreateAccountOp()

	return history2.OperationDetails{
		Type: xdr.OperationTypeCreateAccount,
		CreateAccount: &history2.CreateAccountDetails{
			AccountType: op.AccountType,
		},
	}, nil
}

func (h *createAccountOpHandler) ParticipantsEffects(opBody xdr.OperationBody, source history2.ParticipantEffect) ([]history2.ParticipantEffect, error) {
	participants := []history2.ParticipantEffect{source}

	createAccountOp := opBody.MustCreateAccountOp()

	participants = append(participants, history2.ParticipantEffect{
		AccountID: h.pubKeyConverter.ConvertToInt64(xdr.PublicKey(createAccountOp.Destination)),
	})

	if createAccountOp.Referrer != nil {
		participants = append(participants, history2.ParticipantEffect{
			AccountID: h.pubKeyConverter.ConvertToInt64(xdr.PublicKey(*createAccountOp.Referrer)),
		})
	}

	return participants, nil
}

type manageKeyValueOpHandler struct {
	pubKeyConverter publicKeyConverter
}

func (h *manageKeyValueOpHandler) OperationDetails(opBody xdr.OperationBody) (history2.OperationDetails, error) {
	manageKVOp := opBody.MustManageKeyValueOp()

	var value *xdr.KeyValueEntryValue
	if manageKVOp.Action.Action == xdr.ManageKvActionPut {
		valueForPtr := manageKVOp.Action.MustValue().Value
		value = &valueForPtr
	}

	return history2.OperationDetails{
		Type: xdr.OperationTypeManageKeyValue,
		ManageKeyValue: &history2.ManageKeyValueDetails{
			Key:    string(manageKVOp.Key),
			Action: manageKVOp.Action.Action,
			Value:  value,
		},
	}, nil
}

func (h *manageKeyValueOpHandler) ParticipantsEffects(opBody xdr.OperationBody, source history2.ParticipantEffect) ([]history2.ParticipantEffect, error) {
	return []history2.ParticipantEffect{source}, nil
}

type setFeeOpHandler struct {
	pubKeyConverter publicKeyConverter
}

func (h *setFeeOpHandler) OperationDetails(opBody xdr.OperationBody) (history2.OperationDetails, error) {
	setFeeOp := opBody.MustSetFeesOp()

	if setFeeOp.IsDelete || setFeeOp.Fee == nil {
		return history2.OperationDetails{
			Type:   xdr.OperationTypeSetFees,
			SetFee: &history2.SetFeeDetails{},
		}, nil
	}

	fee := *setFeeOp.Fee

	var accountIDPtr *int64
	if fee.AccountId != nil {
		accountIDInt := h.pubKeyConverter.ConvertToInt64(xdr.PublicKey(*fee.AccountId))
		accountIDPtr = &accountIDInt
	}

	return history2.OperationDetails{
		Type: xdr.OperationTypeSetFees,
		SetFee: &history2.SetFeeDetails{
			AssetCode:   fee.Asset,
			FixedFee:    int64(fee.FixedFee),
			PercentFee:  int64(fee.PercentFee),
			FeeType:     fee.FeeType,
			AccountID:   accountIDPtr,
			AccountType: fee.AccountType,
			Subtype:     int64(fee.Subtype),
			LowerBound:  int64(fee.LowerBound),
			UpperBound:  int64(fee.UpperBound),
		},
	}, nil
}

func (h *setFeeOpHandler) ParticipantsEffects(opBody xdr.OperationBody, source history2.ParticipantEffect) ([]history2.ParticipantEffect, error) {
	participants := []history2.ParticipantEffect{source}

	setFeeOp := opBody.MustSetFeesOp()
	if (setFeeOp.Fee != nil) && (setFeeOp.Fee.AccountId != nil) {
		participants = append(participants, history2.ParticipantEffect{
			AccountID: h.pubKeyConverter.ConvertToInt64(xdr.PublicKey(*setFeeOp.Fee.AccountId)),
		})
	}

	return participants, nil
}

type manageAccountOpHandler struct {
	pubKeyConverter publicKeyConverter
}

func (h *manageAccountOpHandler) OperationDetails(opBody xdr.OperationBody) (history2.OperationDetails, error) {
	manageAccountOp := opBody.MustManageAccountOp()

	return history2.OperationDetails{
		Type: xdr.OperationTypeManageAccount,
		ManageAccount: &history2.ManageAccountDetails{
			AccountID:            h.pubKeyConverter.ConvertToInt64(xdr.PublicKey(manageAccountOp.Account)),
			BlockReasonsToAdd:    int32(manageAccountOp.BlockReasonsToAdd),
			BlockReasonsToRemove: int32(manageAccountOp.BlockReasonsToRemove),
		},
	}, nil
}

func (h *manageAccountOpHandler) ParticipantsEffects(opBody xdr.OperationBody, source history2.ParticipantEffect) ([]history2.ParticipantEffect, error) {
	participants := []history2.ParticipantEffect{source}

	participants = append(participants, history2.ParticipantEffect{
		AccountID: h.pubKeyConverter.ConvertToInt64(xdr.PublicKey(
			opBody.MustManageAccountOp().Account)),
	})

	return participants, nil
}

type createWithdrawRequestOpHandler struct {
	pubKeyConverter publicKeyConverter
}

func (h *createWithdrawRequestOpHandler) OperationDetails(opBody xdr.OperationBody,
) (history2.OperationDetails, error) {
	withdrawRequest := opBody.MustCreateWithdrawalRequestOp().Request

	var externalDetails map[string]interface{}
	json.Unmarshal([]byte(withdrawRequest.ExternalDetails), externalDetails)

	destinationAsset := xdr.AssetCode("")
	destinationAmount := int64(0)
	if autoConvDet, ok := withdrawRequest.Details.GetAutoConversion(); ok {
		destinationAsset = autoConvDet.DestAsset
		destinationAmount = int64(autoConvDet.ExpectedAmount)
	}

	return history2.OperationDetails{
		Type: xdr.OperationTypeCreateWithdrawalRequest,
		CreateWithdrawRequest: &history2.CreateWithdrawRequestDetails{
			BalanceID:         h.pubKeyConverter.ConvertToInt64(xdr.PublicKey(withdrawRequest.Balance)),
			Amount:            int64(withdrawRequest.Amount),
			FixedFee:          int64(withdrawRequest.Fee.Fixed),
			PercentFee:        int64(withdrawRequest.Fee.Percent),
			ExternalDetails:   externalDetails,
			DestinationAsset:  destinationAsset,
			DestinationAmount: destinationAmount,
		},
	}, nil
}

func (h *createWithdrawRequestOpHandler) ParticipantsEffects(opBody xdr.OperationBody,
	source history2.ParticipantEffect,
) ([]history2.ParticipantEffect, error) {
	withdrawRequest := opBody.MustCreateWithdrawalRequestOp().Request
	balanceIDInt := h.pubKeyConverter.ConvertToInt64(xdr.PublicKey(withdrawRequest.Balance))
	amount := int64(withdrawRequest.Amount)

	source.BalanceID = &balanceIDInt
	source.Effect.Type = history2.EffectTypeWithdraw
	source.Effect.WithdrawAmount = &amount

	return []history2.ParticipantEffect{source}, nil
}

type manageBalanceOpHandler struct {
	pubKeyConverter publicKeyConverter
}

func (h *manageBalanceOpHandler) OperationDetails(opBody xdr.OperationBody, opRes xdr.OperationResultTr) (history2.OperationDetails, error) {

}

func (h *manageBalanceOpHandler) ParticipantsEffects(opBody xdr.OperationBody, source history2.ParticipantEffect) ([]history2.ParticipantEffect, error) {

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
