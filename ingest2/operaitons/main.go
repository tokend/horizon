package operaitons

import (
	"encoding/hex"
	"encoding/json"

	"gitlab.com/tokend/horizon/utf8"

	"gitlab.com/tokend/go/amount"

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
			xdr.OperationTypeManageAccount: &manageAccountOpHandler{
				pubKeyConverter: mainProvider.GetPubKeyConverter(),
			},
			xdr.OperationTypeManageBalance: &manageBalanceOpHandler{
				pubKeyConverter: mainProvider.GetPubKeyConverter(),
			},
			xdr.OperationTypeManageKeyValue: &manageKeyValueOpHandler{},
			xdr.OperationTypeSetFees: &setFeeOpHandler{
				pubKeyConverter: mainProvider.GetPubKeyConverter(),
			},
			xdr.OperationTypeCreateWithdrawalRequest: &createWithdrawRequestOpHandler{
				pubKeyConverter: mainProvider.GetPubKeyConverter(),
			},
			xdr.OperationTypeManageLimits: &manageLimitsOpHandler{
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
	if handler == nil {
		return history2.Operation{}, nil, errors.New("no handler for such operation type")
	}

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
	ParticipantsEffects(opBody xdr.OperationBody, opRes xdr.OperationResultTr, source history2.ParticipantEffect) ([]history2.ParticipantEffect, error)
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
	err := json.Unmarshal([]byte(withdrawRequest.ExternalDetails), externalDetails)
	if err != nil {
		externalDetails = make(map[string]interface{})
		externalDetails["invalid_json_data"] = string(withdrawRequest.ExternalDetails)
	}

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
	manageBalanceOp := opBody.MustManageBalanceOp()

	return history2.OperationDetails{
		Type: xdr.OperationTypeManageBalance,
		ManageBalance: &history2.ManageBalanceDetails{
			DestinationAccount: h.pubKeyConverter.ConvertToInt64(xdr.PublicKey(
				manageBalanceOp.Destination)),
			Action: manageBalanceOp.Action,
			Asset:  manageBalanceOp.Asset,
			BalanceID: h.pubKeyConverter.ConvertToInt64(xdr.PublicKey(
				opRes.MustManageBalanceResult().MustSuccess().BalanceId)),
		},
	}, nil
}

func (h *manageBalanceOpHandler) ParticipantsEffects(opBody xdr.OperationBody, opRes xdr.OperationResultTr, source history2.ParticipantEffect) ([]history2.ParticipantEffect, error) {
	participants := []history2.ParticipantEffect{source}

	manageBalanceOp := opBody.MustManageBalanceOp()

	destinationAccount := h.pubKeyConverter.ConvertToInt64(xdr.PublicKey(manageBalanceOp.Destination))
	destinationBalance := h.pubKeyConverter.ConvertToInt64(xdr.PublicKey(
		opRes.MustManageBalanceResult().MustSuccess().BalanceId))
	if source.AccountID != destinationAccount {
		participants = append(participants, history2.ParticipantEffect{
			AccountID: destinationAccount,
			BalanceID: &destinationBalance,
			AssetCode: &manageBalanceOp.Asset,
			// maybe add effect - creation
		})
	} else {
		participants[0].BalanceID = &destinationBalance
		participants[0].AssetCode = &manageBalanceOp.Asset
	}

	return participants, nil
}

type manageLimitsOpHandler struct {
	pubKeyConverter publicKeyConverter
}

func (h *manageLimitsOpHandler) OperationDetails(opBody xdr.OperationBody, opRes xdr.OperationResultTr) (history2.OperationDetails, error) {
	manageLimitsOp := opBody.MustManageLimitsOp()

	opDetails := history2.OperationDetails{
		Type: xdr.OperationTypeManageLimits,
		ManageLimits: &history2.ManageLimitsDetails{
			Action: manageLimitsOp.Details.Action,
		},
	}

	switch opDetails.ManageLimits.Action {
	case xdr.ManageLimitsActionCreate:
		creationDetails := manageLimitsOp.Details.MustLimitsCreateDetails()

		var accountID *int64
		if creationDetails.AccountId != nil {
			accountIDInt := h.pubKeyConverter.ConvertToInt64(xdr.PublicKey(
				*creationDetails.AccountId))
			accountID = &accountIDInt
		}

		opDetails.ManageLimits.Creation = &history2.ManageLimitsCreationDetails{
			AccountID:       accountID,
			AccountType:     creationDetails.AccountType,
			StatsOpType:     creationDetails.StatsOpType,
			AssetCode:       creationDetails.AssetCode,
			IsConvertNeeded: creationDetails.IsConvertNeeded,
			DailyOut:        uint64(creationDetails.DailyOut),
			WeeklyOut:       uint64(creationDetails.WeeklyOut),
			MonthlyOut:      uint64(creationDetails.MonthlyOut),
			AnnualOut:       uint64(creationDetails.AnnualOut),
		}
	case xdr.ManageLimitsActionRemove:
		opDetails.ManageLimits.Removal = &history2.ManageLimitsRemovalDetails{
			ID: int64(manageLimitsOp.Details.MustId()),
		}
	default:
		return history2.OperationDetails{}, errors.New("unexpected manage limits action")
	}

	return opDetails, nil
}

func (h *manageLimitsOpHandler) ParticipantsEffects(opBody xdr.OperationBody,
	opRes xdr.OperationResultTr, source history2.ParticipantEffect,
) ([]history2.ParticipantEffect, error) {
	participants := []history2.ParticipantEffect{source}

	manageLimitsOp := opBody.MustManageLimitsOp()

	if manageLimitsOp.Details.Action != xdr.ManageLimitsActionCreate {
		return participants, nil
	}

	creationDetails := manageLimitsOp.Details.MustLimitsCreateDetails()

	if creationDetails.AccountId == nil {
		return participants, nil
	}

	accountID := h.pubKeyConverter.ConvertToInt64(xdr.PublicKey(
		*creationDetails.AccountId))

	if source.AccountID == accountID {
		return participants, nil
	}

	participants = append(participants, history2.ParticipantEffect{
		AccountID: accountID,
		AssetCode: &creationDetails.AssetCode,
	})

	return participants, nil
}

type createManageLimitsRequestOpHandler struct {
}

func (h *createManageLimitsRequestOpHandler) OperationDetails(opBody xdr.OperationBody,
	opRes xdr.OperationResultTr,
) (history2.OperationDetails, error) {
	createManageLimitsRequestOp := opBody.MustCreateManageLimitsRequestOp()

	data := make(map[string]interface{})
	rawData, ok := createManageLimitsRequestOp.ManageLimitsRequest.Ext.GetDetails()
	if !ok {
		data["invalid_json_details"] = "empty create manage limits json details"
	}

	err := json.Unmarshal([]byte(rawData), &data)
	if err != nil {
		data["invalid_json_data"] = string(rawData)
	}

	var requestID int64
	if rawID, ok := createManageLimitsRequestOp.Ext.GetRequestId(); ok {
		requestID = int64(rawID)
	}

	return history2.OperationDetails{
		Type: xdr.OperationTypeCreateManageLimitsRequest,
		CreateManageLimitsRequest: &history2.CreateManageLimitsRequestDetails{
			Data:      data,
			RequestID: requestID,
		},
	}, nil
}

func (h *createManageLimitsRequestOpHandler) ParticipantsEffects(opBody xdr.OperationBody,
	opRes xdr.OperationResultTr, source history2.ParticipantEffect,
) ([]history2.ParticipantEffect, error) {
	return []history2.ParticipantEffect{source}, nil
}

type manageAssetPairOpHadler struct {
}

func (h *manageAssetPairOpHadler) OperationDetails(opBody xdr.OperationBody, opRes xdr.OperationResultTr,
) (history2.OperationDetails, error) {
	manageAssetPairOp := opBody.MustManageAssetPairOp()

	return history2.OperationDetails{
		Type: xdr.OperationTypeManageAssetPair,
		ManageAssetPair: &history2.ManageAssetPairDetails{
			BaseAsset:               manageAssetPairOp.Base,
			QuoteAsset:              manageAssetPairOp.Quote,
			PhysicalPrice:           int64(manageAssetPairOp.PhysicalPrice),
			PhysicalPriceCorrection: int64(manageAssetPairOp.PhysicalPriceCorrection),
			MaxPriceStep:            int64(manageAssetPairOp.MaxPriceStep),
			PoliciesI:               int32(manageAssetPairOp.Policies),
		},
	}, nil
}

func (h *manageAssetPairOpHadler) ParticipantsEffects(opBody xdr.OperationBody,
	opRes xdr.OperationResultTr, source history2.ParticipantEffect,
) ([]history2.ParticipantEffect, error) {
	return []history2.ParticipantEffect{source}, nil
}

type manageOfferOpHandler struct {
	pubKeyConverter publicKeyConverter
}

func (h *manageOfferOpHandler) OperationDetails(opBody xdr.OperationBody, opRes xdr.OperationResultTr,
) (history2.OperationDetails, error) {
	manageOfferOp := opBody.MustManageOfferOp()
	manageOfferOpRes := opRes.MustManageOfferResult().MustSuccess()

	offerID := int64(manageOfferOp.OfferId)
	isDeleted := manageOfferOpRes.Offer.Effect == xdr.ManageOfferEffectDeleted
	if !isDeleted {
		offerID = int64(manageOfferOpRes.Offer.MustOffer().OfferId)
	}

	return history2.OperationDetails{
		Type: xdr.OperationTypeManageOffer,
		ManageOffer: &history2.ManageOfferDetails{
			OfferID:     offerID,
			OrderBookID: int64(manageOfferOp.OrderBookId),
			BaseAsset:   manageOfferOpRes.BaseAsset,
			QuoteAsset:  manageOfferOpRes.QuoteAsset,
			Amount:      amount.String(int64(manageOfferOp.Amount)),
			Price:       amount.String(int64(manageOfferOp.Price)),
			IsBuy:       manageOfferOp.IsBuy,
			Fee:         amount.String(int64(manageOfferOp.Fee)),
			IsDeleted:   isDeleted,
		},
	}, nil
}

func (h *manageOfferOpHandler) ParticipantsEffects(opBody xdr.OperationBody,
	opRes xdr.OperationResultTr, baseSource history2.ParticipantEffect,
) ([]history2.ParticipantEffect, error) {
	manageOfferOp := opBody.MustManageOfferOp()
	manageOfferOpRes := opRes.MustManageOfferResult().MustSuccess()

	participants := h.participantEffects(
		baseSource, manageOfferOp.BaseBalance, manageOfferOp.QuoteBalance,
		manageOfferOpRes, int64(manageOfferOp.Amount))

	for _, claimedOffer := range manageOfferOpRes.OffersClaimed {
		participantBase := history2.ParticipantEffect{
			AccountID: h.pubKeyConverter.ConvertToInt64(xdr.PublicKey(claimedOffer.BAccountId)),
		}

		participants = append(participants,
			h.participantEffects(
				participantBase, claimedOffer.BaseBalance, claimedOffer.QuoteBalance,
				manageOfferOpRes, int64(manageOfferOp.Amount))...)
	}

	return participants, nil
}

func (h *manageOfferOpHandler) participantEffects(participantBase history2.ParticipantEffect,
	baseBalance, quoteBalance xdr.BalanceId, manageOfferOpRes xdr.ManageOfferSuccessResult, baseAmount int64,
) []history2.ParticipantEffect {
	participantBaseBalanceID := h.pubKeyConverter.ConvertToInt64(xdr.PublicKey(baseBalance))
	participantQuoteBalanceID := h.pubKeyConverter.ConvertToInt64(xdr.PublicKey(quoteBalance))

	participantBase.Effect = history2.Effect{
		Type: history2.EffectTypeOffer,
		Offer: &history2.OfferEffect{
			BaseBalanceID:  participantBaseBalanceID,
			QuoteBalanceID: participantQuoteBalanceID,
			BaseAmount:     amount.String(baseAmount),
			BaseAsset:      manageOfferOpRes.BaseAsset,
			QuoteAsset:     manageOfferOpRes.QuoteAsset,
		},
	}

	participantQuote := participantBase
	participantBase.BalanceID = &participantBaseBalanceID
	participantBase.AssetCode = &manageOfferOpRes.BaseAsset
	participantQuote.BalanceID = &participantQuoteBalanceID
	participantQuote.AssetCode = &manageOfferOpRes.QuoteAsset

	return []history2.ParticipantEffect{participantBase, participantQuote}
}

type manageInvoiceRequestOpHandler struct {
	pubKeyConverter publicKeyConverter
}

func (h *manageInvoiceRequestOpHandler) OperationDetails(opBody xdr.OperationBody,
	opRes xdr.OperationResultTr,
) (history2.OperationDetails, error) {
	manageInvoiceRequestOp := opBody.MustManageInvoiceRequestOp()
	manageInvoiceRequestOpRes := opRes.MustManageInvoiceRequestResult().MustSuccess()

	opDetails := history2.OperationDetails{
		Type: xdr.OperationTypeManageInvoiceRequest,
		ManageInvoiceRequest: &history2.ManageInvoiceRequestDetails{
			Action: manageInvoiceRequestOp.Details.Action,
		},
	}

	switch manageInvoiceRequestOp.Details.Action {
	case xdr.ManageInvoiceRequestActionCreate:
		creationDetails := manageInvoiceRequestOp.Details.MustInvoiceRequest()

		var contractID *int64
		if creationDetails.ContractId != nil {
			contractIDInt := int64(*creationDetails.ContractId)
			contractID = &contractIDInt
		}

		var details map[string]interface{}
		err := json.Unmarshal([]byte(creationDetails.Details), details)
		if err != nil {
			details = make(map[string]interface{})
			details["invalid_invoice_details"] = creationDetails.Details
		}

		opDetails.ManageInvoiceRequest.Create = &history2.CreateInvoiceRequestDetails{
			Amount:     amount.StringU(uint64(creationDetails.Amount)),
			Sender:     h.pubKeyConverter.ConvertToInt64(xdr.PublicKey(creationDetails.Sender)),
			RequestID:  int64(manageInvoiceRequestOpRes.Details.MustResponse().RequestId),
			Asset:      creationDetails.Asset,
			ContractID: contractID,
			Details:    details,
		}
	case xdr.ManageInvoiceRequestActionRemove:
		opDetails.ManageInvoiceRequest.Remove = &history2.RemoveInvoiceRequestDetails{
			RequestID: int64(manageInvoiceRequestOp.Details.MustRequestId()),
		}
	}

	return opDetails, nil
}

type manageContractRequestOpHandler struct {
	pubKeyConverter publicKeyConverter
}

func (h *manageContractRequestOpHandler) OperationDetails(opBody xdr.OperationBody,
	opRes xdr.OperationResultTr,
) (history2.OperationDetails, error) {
	manageContractRequestOp := opBody.MustManageContractRequestOp()

	opDetails := history2.OperationDetails{
		Type: xdr.OperationTypeManageContractRequest,
		ManageContractRequest: &history2.ManageContractRequestDetails{
			Action: manageContractRequestOp.Details.Action,
		},
	}

	switch opDetails.ManageContractRequest.Action {
	case xdr.ManageContractRequestActionCreate:
		creationDetails := manageContractRequestOp.Details.MustContractRequest()

		var contractDetails map[string]interface{}
		err := json.Unmarshal([]byte(creationDetails.Details), &contractDetails)
		if err != nil {
			contractDetails = make(map[string]interface{})
			contractDetails["invalid_json_details"] = string(creationDetails.Details)
		}

		opDetails.ManageContractRequest.Create = &history2.CreateContractRequestDetails{
			Customer:  h.pubKeyConverter.ConvertToInt64(xdr.PublicKey(creationDetails.Customer)),
			Escrow:    h.pubKeyConverter.ConvertToInt64(xdr.PublicKey(creationDetails.Escrow)),
			Details:   contractDetails,
			StartTime: int64(creationDetails.StartTime),
			EndTime:   int64(creationDetails.EndTime),
			RequestID: int64(opRes.MustManageContractRequestResult().MustSuccess().Details.MustResponse().RequestId),
		}
	case xdr.ManageContractRequestActionRemove:
		opDetails.ManageContractRequest.Remove = &history2.RemoveContractReqeustDetails{
			RequestID: int64(manageContractRequestOp.Details.MustRequestId()),
		}
	default:
		return history2.OperationDetails{}, errors.New("unexpected manage contract request action")
	}

	return opDetails, nil
}

func (h *manageContractRequestOpHandler) ParticipantsEffects(opBody xdr.OperationBody,
	opRes xdr.OperationResultTr, source history2.ParticipantEffect,
) ([]history2.ParticipantEffect, error) {
	creationDetails := opBody.MustManageContractRequestOp().Details.MustContractRequest()

	participants := []history2.ParticipantEffect{source}

	participants = append(participants, history2.ParticipantEffect{
		AccountID: h.pubKeyConverter.ConvertToInt64(xdr.PublicKey(creationDetails.Customer)),
	})

	participants = append(participants, history2.ParticipantEffect{
		AccountID: h.pubKeyConverter.ConvertToInt64(xdr.PublicKey(creationDetails.Escrow)),
	})

	return participants, nil
}

type manageContractOpHandler struct {
}

func (h *manageContractOpHandler) OperationDetails(opBody xdr.OperationBody,
	opRes xdr.OperationResultTr,
) (history2.OperationDetails, error) {
	manageContractOp := opBody.MustManageContractOp()

	opDetails := history2.OperationDetails{
		Type: xdr.OperationTypeManageContract,
		ManageContract: &history2.ManageContractDetails{
			ContractID: int64(manageContractOp.ContractId),
			Action:     manageContractOp.Data.Action,
		},
	}

	switch opDetails.ManageContract.Action {
	case xdr.ManageContractActionAddDetails:
		var newDetails map[string]interface{}

		err := json.Unmarshal([]byte(manageContractOp.Data.MustDetails()), &newDetails)
		if err != nil {
			newDetails = make(map[string]interface{})
			newDetails["invalid_json_data"] = string(manageContractOp.Data.MustDetails())
		}

		opDetails.ManageContract.Details = newDetails
	case xdr.ManageContractActionConfirmCompleted:
		isCompeted := opRes.MustManageContractResult().MustResponse().Data.MustIsCompleted()

		opDetails.ManageContract.IsCompleted = &isCompeted
	case xdr.ManageContractActionStartDispute:
		var disputeReason map[string]interface{}

		err := json.Unmarshal([]byte(manageContractOp.Data.MustDisputeReason()), &disputeReason)
		if err != nil {
			disputeReason = make(map[string]interface{})
			disputeReason["invalid_json_data"] = string(manageContractOp.Data.MustDisputeReason())
		}

		opDetails.ManageContract.Details = disputeReason
	case xdr.ManageContractActionResolveDispute:
		isRevert := manageContractOp.Data.MustIsRevert()

		opDetails.ManageContract.IsRevert = &isRevert
	default:
		return history2.OperationDetails{}, errors.New("unexpected manage contract actions")
	}

	return opDetails, nil
}

func (h *manageContractOpHandler) ParticipantsEffects(opBody xdr.OperationBody,
	opRes xdr.OperationResultTr, source history2.ParticipantEffect,
) ([]history2.ParticipantEffect, error) {
	return []history2.ParticipantEffect{source}, nil
}

type reviewRequestOpHandler struct {
	pubKeyConverter publicKeyConverter
}

func (h *reviewRequestOpHandler) OperationDetails(opBody xdr.OperationBody, opRes xdr.OperationResultTr) (history2.OperationDetails, error) {
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
	participants := []history2.ParticipantEffect{source}

	extendedResult, ok := opRes.MustReviewRequestResult().MustSuccess().Ext.GetExtendedResult()
	if !ok {
		return participants, nil
	}

	atomicSwapExtendedResult, ok := extendedResult.TypeExt.GetASwapExtended()
	if !ok {
		return participants, nil
	}

	ownerBalanceID := h.pubKeyConverter.ConvertToInt64(xdr.PublicKey(atomicSwapExtendedResult.BidOwnerBaseBalanceId))

	participants = append(participants, history2.ParticipantEffect{
		AccountID: h.pubKeyConverter.ConvertToInt64(xdr.PublicKey(atomicSwapExtendedResult.BidOwnerId)),
		BalanceID: &ownerBalanceID,
		AssetCode: &atomicSwapExtendedResult.baseAsset,
	})

	purchaserBaseBalanceId := h.pubKeyConverter.ConvertToInt64(xdr.PublicKey(atomicSwapExtendedResult.PurchaserBaseBalanceId))

	participants = append(participants, history2.ParticipantEffect{
		AccountID: h.pubKeyConverter.ConvertToInt64(xdr.PublicKey(atomicSwapExtendedResult.PurchaserId)),
		BalanceID: &purchaserBaseBalanceId,
		AssetCode: &atomicSwapExtendedResult.baseAsset,
	})

	return participants, nil
}

type manageAssetOpHandler struct {
	pubKeyConverter publicKeyConverter
}

func (h *manageAssetOpHandler) OperationDetails(opBody xdr.OperationBody, opRes xdr.OperationResultTr,
) (history2.OperationDetails, error) {
	manageAssetOp := opBody.MustManageAssetOp()

	opDetails := history2.OperationDetails{
		Type: xdr.OperationTypeManageAsset,
		ManageAsset: &history2.ManageAssetDetails{
			RequestID: int64(manageAssetOp.RequestId),
			Action:    manageAssetOp.Request.Action,
		},
	}

	if manageAssetOp.RequestId == 0 {
		opDetails.ManageAsset.RequestID = int64(opRes.MustManageAssetResult().MustSuccess().RequestId)
	}

	switch opDetails.ManageAsset.Action {
	case xdr.ManageAssetActionCreateAssetCreationRequest:
		creationDetails := manageAssetOp.Request.MustCreateAsset()

		var details map[string]interface{}
		err := json.Unmarshal([]byte(creationDetails.Details), &details)
		if err != nil {
			details = make(map[string]interface{})
			details["invalid_json_details"] = string(creationDetails.Details)
			// error  - error from unmarshal,
			// data - rawDetails
			// custom unmarshl details,
		}

		policies := int32(creationDetails.Policies)
		preissuedSigner := h.pubKeyConverter.ConvertToInt64(xdr.PublicKey(creationDetails.PreissuedAssetSigner))

		opDetails.ManageAsset.AssetCode = creationDetails.Code
		opDetails.ManageAsset.Details = details
		opDetails.ManageAsset.Policies = &policies
		opDetails.ManageAsset.PreissuedSigner = &preissuedSigner
		opDetails.ManageAsset.MaxIssuanceAmount = amount.StringU(uint64(creationDetails.MaxIssuanceAmount))
	case xdr.ManageAssetActionCreateAssetUpdateRequest:
		updateDetails := manageAssetOp.Request.MustUpdateAsset()

		var details map[string]interface{}
		err := json.Unmarshal([]byte(updateDetails.Details), &details)
		if err != nil {
			details = make(map[string]interface{})
			details["invalid_json_details"] = string(updateDetails.Details)
		}

		policies := int32(updateDetails.Policies)

		opDetails.ManageAsset.AssetCode = updateDetails.Code
		opDetails.ManageAsset.Details = details
		opDetails.ManageAsset.Policies = &policies
	case xdr.ManageAssetActionCancelAssetRequest:
	case xdr.ManageAssetActionChangePreissuedAssetSigner:
		data := manageAssetOp.Request.MustChangePreissuedSigner()
		preissuedSigner := h.pubKeyConverter.ConvertToInt64(xdr.PublicKey(data.AccountId))

		opDetails.ManageAsset.AssetCode = data.Code
		opDetails.ManageAsset.PreissuedSigner = &preissuedSigner
	case xdr.ManageAssetActionUpdateMaxIssuance:
		data := manageAssetOp.Request.MustUpdateMaxIssuance()

		opDetails.ManageAsset.AssetCode = data.AssetCode
		opDetails.ManageAsset.MaxIssuanceAmount = amount.StringU(uint64(data.MaxIssuanceAmount))
	default:
		return history2.OperationDetails{}, errors.New("unexpected manage asset action")
	}

	return opDetails, nil
}

func (h *manageAssetOpHandler) ParticipantsEffects(opBody xdr.OperationBody,
	opRes xdr.OperationResultTr, source history2.ParticipantEffect,
) ([]history2.ParticipantEffect, error) {
	manageAssetOp := opBody.MustManageAssetOp()

	participants := []history2.ParticipantEffect{source}

	switch manageAssetOp.Request.Action {
	case xdr.ManageAssetActionCreateAssetCreationRequest:
		creationDetails := manageAssetOp.Request.MustCreateAsset()

		participants = append(participants, history2.ParticipantEffect{
			AccountID: h.pubKeyConverter.ConvertToInt64(xdr.PublicKey(creationDetails.PreissuedAssetSigner)),
			AssetCode: &creationDetails.Code,
		})
	case xdr.ManageAssetActionChangePreissuedAssetSigner:
		updateDetails := manageAssetOp.Request.MustChangePreissuedSigner()

		participants = append(participants, history2.ParticipantEffect{
			AccountID: h.pubKeyConverter.ConvertToInt64(xdr.PublicKey(updateDetails.AccountId)),
			AssetCode: &updateDetails.Code,
		})
	}

	return participants, nil
}

type createPreissuanceRequestOpHandler struct {
}

func (h *createPreissuanceRequestOpHandler) OperationDetails(opBody xdr.OperationBody,
	opRes xdr.OperationResultTr,
) (history2.OperationDetails, error) {
	preissuanceRequest := opBody.MustCreatePreIssuanceRequest().Request
	successResult := opRes.MustCreatePreIssuanceRequestResult().MustSuccess()

	return history2.OperationDetails{
		Type: xdr.OperationTypeCreatePreissuanceRequest,
		CreatePreIssuanceRequest: &history2.CreatePreIssuanceRequestDetails{
			AssetCode:   preissuanceRequest.Asset,
			Amount:      amount.StringU(uint64(preissuanceRequest.Amount)),
			RequestID:   int64(successResult.RequestId),
			IsFulfilled: successResult.Fulfilled,
		},
	}, nil
}

func (h *createPreissuanceRequestOpHandler) ParticipantsEffects(opBody xdr.OperationBody,
	opRes xdr.OperationResultTr, source history2.ParticipantEffect,
) ([]history2.ParticipantEffect, error) {
	return []history2.ParticipantEffect{source}, nil
}

type createIssuanceRequestOpHandler struct {
	pubKeyConverter publicKeyConverter
}

func (h *createIssuanceRequestOpHandler) OperationDetails(opBody xdr.OperationBody, opRes xdr.OperationResultTr) (history2.OperationDetails, error) {
	createIssuanceRequestOp := opBody.MustCreateIssuanceRequestOp()
	issuanceRequest := createIssuanceRequestOp.Request

	var externalDetails map[string]interface{}
	err := json.Unmarshal([]byte(issuanceRequest.ExternalDetails), &externalDetails)
	if err != nil {
		externalDetails = make(map[string]interface{})
		externalDetails["invalid_json_details"] = string(issuanceRequest.ExternalDetails)
	}

	var allTasks *int64
	rawAllTasks, ok := createIssuanceRequestOp.Ext.GetAllTasks()
	if ok && rawAllTasks != nil {
		allTasksInt := int64(*rawAllTasks)
		allTasks = &allTasksInt
	}

	return history2.OperationDetails{
		Type: xdr.OperationTypeCreateIssuanceRequest,
		CreateIssuanceRequest: &history2.CreateIssuanceRequestDetails{
			FixedFee:        amount.StringU(uint64(issuanceRequest.Fee.Fixed)),
			PercentFee:      amount.StringU(uint64(issuanceRequest.Fee.Percent)),
			Reference:       utf8.Scrub(string(createIssuanceRequestOp.Reference)),
			Amount:          amount.StringU(uint64(issuanceRequest.Amount)),
			Asset:           issuanceRequest.Asset,
			BalanceID:       h.pubKeyConverter.ConvertToInt64(xdr.PublicKey(issuanceRequest.Receiver)),
			ExternalDetails: externalDetails,
			AllTasks:        allTasks,
		},
	}, nil
}

func (h *createIssuanceRequestOpHandler) ParticipantsEffects(opBody xdr.OperationBody,
	opRes xdr.OperationResultTr, source history2.ParticipantEffect,
) ([]history2.ParticipantEffect, error) {

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
