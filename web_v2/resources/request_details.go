package resources

import (
	"gitlab.com/distributed_lab/logan/v3"
	"gitlab.com/distributed_lab/logan/v3/errors"
	"gitlab.com/tokend/go/xdr"
	"gitlab.com/tokend/horizon/db2/history2"
	regources "gitlab.com/tokend/regources/generated"
)

// NewRequestDetails - returns new instance of reviewable request details
func NewRequestDetails(request history2.ReviewableRequest) regources.Resource {
	switch request.RequestType {
	case xdr.ReviewableRequestTypeUpdateAsset:
		return newAssetUpdateRequest(request.ID, *request.Details.UpdateAsset)
	case xdr.ReviewableRequestTypeCreatePreIssuance:
		return newPreIssuanceRequest(request.ID, *request.Details.CreatePreIssuance)
	case xdr.ReviewableRequestTypeCreateIssuance:
		return newIssuanceRequest(request.ID, *request.Details.CreateIssuance)
	case xdr.ReviewableRequestTypeCreateWithdraw:
		return newWithdrawalRequest(request.ID, *request.Details.CreateWithdraw)
	case xdr.ReviewableRequestTypeCreateSale:
		return newSaleRequest(request.ID, *request.Details.CreateSale)
	case xdr.ReviewableRequestTypeUpdateLimits:
		return newLimitsUpdateRequest(request.ID, *request.Details.UpdateLimits)
	case xdr.ReviewableRequestTypeCreateAmlAlert:
		return newAmlAlertRequest(request.ID, *request.Details.CreateAmlAlert)
	case xdr.ReviewableRequestTypeChangeRole:
		return newChangeRoleRequest(request.ID, *request.Details.ChangeRole)
	case xdr.ReviewableRequestTypeUpdateSaleDetails:
		return newUpdateSaleDetailsRequest(request.ID, *request.Details.UpdateSaleDetails)
	case xdr.ReviewableRequestTypeCreateAsset:
		return newAssetCreateRequest(request.ID, *request.Details.CreateAsset)
	case xdr.ReviewableRequestTypeCreateAtomicSwapAsk:
		return newAtomicSwapAskRequest(request.ID, *request.Details.CreateAtomicSwapAsk)
	case xdr.ReviewableRequestTypeCreateAtomicSwapBid:
		return newAtomicSwapBidRequest(request.ID, *request.Details.CreateAtomicSwapBid)
	case xdr.ReviewableRequestTypeCreatePoll:
		return newCreatePollRequest(request.ID, *request.Details.CreatePoll)
	case xdr.ReviewableRequestTypeKycRecovery:
		return newKYCRecoveryRequest(request.ID, *request.Details.KYCRecovery)
	case xdr.ReviewableRequestTypeManageOffer:
		return newManageOfferRequest(request.ID, *request.Details.ManageOffer)
	case xdr.ReviewableRequestTypeCreatePayment:
		return newCreatePaymentRequest(request.ID, *request.Details.CreatePayment)
	default:
		panic(errors.From(errors.New("unexpected operation type"), logan.F{
			"type": request.RequestType,
		}))
	}
	return nil
}

func newAmlAlertRequest(id int64, details history2.CreateAmlAlertRequest) *regources.CreateAmlAlertRequest {
	return &regources.CreateAmlAlertRequest{
		Key: regources.NewKeyInt64(id, regources.REQUEST_DETAILS_AML_ALERT),
		Attributes: regources.CreateAmlAlertRequestAttributes{
			Amount:         details.Amount,
			CreatorDetails: details.CreatorDetails,
		},
		Relationships: regources.CreateAmlAlertRequestRelationships{
			Balance: NewBalanceKey(details.BalanceID).AsRelation(),
		},
	}
}
func newAssetCreateRequest(id int64, details history2.CreateAssetRequest) *regources.CreateAssetRequest {
	return &regources.CreateAssetRequest{
		Key: regources.NewKeyInt64(id, regources.REQUEST_DETAILS_ASSET_CREATE),
		Attributes: regources.CreateAssetRequestAttributes{
			Asset:                  details.Asset,
			Policies:               details.Policies,
			PreIssuanceAssetSigner: details.PreIssuedAssetSigner,
			MaxIssuanceAmount:      details.MaxIssuanceAmount,
			InitialPreissuedAmount: details.InitialPreissuedAmount,
			CreatorDetails:         details.CreatorDetails,
			Type:                   details.Type,
			TrailingDigitsCount:    details.TrailingDigitsCount,
		},
	}
}
func newAssetUpdateRequest(id int64, details history2.UpdateAssetRequest) *regources.UpdateAssetRequest {
	return &regources.UpdateAssetRequest{
		Key: regources.NewKeyInt64(id, regources.REQUEST_DETAILS_ASSET_UPDATE),
		Attributes: regources.UpdateAssetRequestAttributes{
			Policies:       details.Policies,
			CreatorDetails: details.CreatorDetails,
		},
		Relationships: regources.UpdateAssetRequestRelationships{
			Asset: NewAssetKey(details.Asset).AsRelation(),
		},
	}
}
func newAtomicSwapBidRequest(id int64, details history2.CreateAtomicSwapBidRequest) *regources.CreateAtomicSwapBidRequest {
	return &regources.CreateAtomicSwapBidRequest{
		Key: regources.NewKeyInt64(id, regources.REQUEST_DETAILS_ATOMIC_SWAP_BID),
		Attributes: regources.CreateAtomicSwapBidRequestAttributes{
			BaseAmount:     regources.Amount(details.BaseAmount),
			CreatorDetails: details.CreatorDetails,
		},
		Relationships: regources.CreateAtomicSwapBidRequestRelationships{
			Ask:        regources.NewKeyInt64(int64(details.AskID), regources.ATOMIC_SWAP_ASK).AsRelation(),
			QuoteAsset: NewAtomicSwapAskQuoteAssetKey(details.QuoteAsset, details.AskID).AsRelation(),
		},
	}
}
func newAtomicSwapAskRequest(id int64, details history2.CreateAtomicSwapAskRequest) *regources.CreateAtomicSwapAskRequest {
	quoteAssets := &regources.RelationCollection{
		Data: make([]regources.Key, 0, len(details.QuoteAssets)),
	}
	for _, quoteAsset := range details.QuoteAssets {
		quoteAssets.Data = append(quoteAssets.Data, newQuoteAssetKey(quoteAsset.Asset))
	}

	return &regources.CreateAtomicSwapAskRequest{
		Key: regources.NewKeyInt64(id, regources.REQUEST_DETAILS_ATOMIC_SWAP_ASK),
		Attributes: regources.CreateAtomicSwapAskRequestAttributes{
			BaseAmount:     regources.Amount(details.BaseAmount),
			CreatorDetails: details.CreatorDetails,
		},
		Relationships: regources.CreateAtomicSwapAskRequestRelationships{
			BaseBalance: NewBalanceKey(details.BaseBalance).AsRelation(),
			QuoteAssets: quoteAssets,
		},
	}
}
func newIssuanceRequest(id int64, details history2.CreateIssuanceRequest) *regources.CreateIssuanceRequest {
	return &regources.CreateIssuanceRequest{
		Key: regources.NewKeyInt64(id, regources.REQUEST_DETAILS_ISSUANCE),
		Attributes: regources.CreateIssuanceRequestAttributes{
			Amount:         details.Amount,
			CreatorDetails: details.CreatorDetails,
		},
		Relationships: regources.CreateIssuanceRequestRelationships{
			Asset:    NewAssetKey(details.Asset).AsRelation(),
			Receiver: NewBalanceKey(details.Receiver).AsRelation(),
		},
	}
}
func newLimitsUpdateRequest(id int64, details history2.UpdateLimitsRequest) *regources.UpdateLimitsRequest {
	return &regources.UpdateLimitsRequest{
		Key: regources.NewKeyInt64(id, regources.REQUEST_DETAILS_LIMITS_UPDATE),
		Attributes: regources.UpdateLimitsRequestAttributes{
			CreatorDetails: details.CreatorDetails,
		},
	}
}
func newPreIssuanceRequest(id int64, details history2.CreatePreIssuanceRequest) *regources.CreatePreIssuanceRequest {
	return &regources.CreatePreIssuanceRequest{
		Key: regources.NewKeyInt64(id, regources.REQUEST_DETAILS_PRE_ISSUANCE),
		Attributes: regources.CreatePreIssuanceRequestAttributes{
			Amount:    details.Amount,
			Signature: details.Signature,
			Reference: details.Reference,
		},
		Relationships: regources.CreatePreIssuanceRequestRelationships{
			Asset: NewAssetKey(details.Asset).AsRelation(),
		},
	}
}
func newSaleRequest(id int64, details history2.CreateSaleRequest) *regources.CreateSaleRequest {
	quoteAssets := &regources.RelationCollection{
		Data: make([]regources.Key, 0, len(details.QuoteAssets)),
	}
	for _, quoteAsset := range details.QuoteAssets {
		quoteAssets.Data = append(quoteAssets.Data, newQuoteAssetKey(quoteAsset.Asset))
	}

	return &regources.CreateSaleRequest{
		Key: regources.NewKeyInt64(id, regources.REQUEST_DETAILS_SALE),
		Attributes: regources.CreateSaleRequestAttributes{
			AccessDefinitionType: details.AccessDefinitionType,
			BaseAssetForHardCap:  details.BaseAssetForHardCap,
			SoftCap:              details.SoftCap,
			HardCap:              details.HardCap,
			StartTime:            details.StartTime,
			EndTime:              details.EndTime,
			SaleType:             details.SaleType,
			CreatorDetails:       details.CreatorDetails,
		},
		Relationships: regources.CreateSaleRequestRelationships{
			QuoteAssets:       quoteAssets,
			BaseAsset:         NewAssetKey(details.BaseAsset).AsRelation(),
			DefaultQuoteAsset: newQuoteAssetKey(details.DefaultQuoteAsset).AsRelation(),
		},
	}
}
func newChangeRoleRequest(id int64, details history2.ChangeRoleRequest) *regources.ChangeRoleRequest {
	return &regources.ChangeRoleRequest{
		Key: regources.NewKeyInt64(id, regources.REQUEST_DETAILS_CHANGE_ROLE),
		Attributes: regources.ChangeRoleRequestAttributes{
			AccountRoleToSet: details.AccountRoleToSet,
			CreatorDetails:   details.CreatorDetails,
			SequenceNumber:   details.SequenceNumber,
		},
		Relationships: regources.ChangeRoleRequestRelationships{
			AccountToUpdateRole: NewAccountKey(details.DestinationAccount).AsRelation(),
		},
	}
}
func newUpdateSaleDetailsRequest(id int64, details history2.UpdateSaleDetailsRequest) *regources.UpdateSaleDetailsRequest {
	return &regources.UpdateSaleDetailsRequest{
		Key: regources.NewKeyInt64(id, regources.REQUEST_DETAILS_UPDATE_SALE_DETAILS),
		Attributes: regources.UpdateSaleDetailsRequestAttributes{
			CreatorDetails: details.CreatorDetails,
		},
		Relationships: regources.UpdateSaleDetailsRequestRelationships{
			Sale: NewSaleKey(int64(details.SaleID)).AsRelation(),
		},
	}
}
func newWithdrawalRequest(id int64, details history2.CreateWithdrawalRequest) *regources.CreateWithdrawRequest {
	return &regources.CreateWithdrawRequest{
		Key: regources.NewKeyInt64(id, regources.REQUEST_DETAILS_WITHDRAWAL),
		Attributes: regources.CreateWithdrawRequestAttributes{
			Fee:            details.Fee,
			Amount:         details.Amount,
			CreatorDetails: details.CreatorDetails,
		},
		Relationships: regources.CreateWithdrawRequestRelationships{
			Balance: NewBalanceKey(details.BalanceID).AsRelation(),
			Asset:   NewAssetKey(details.Asset).AsRelation(),
		},
	}
}

func newCreatePollRequest(id int64, details history2.CreatePollRequest) *regources.CreatePollRequest {
	return &regources.CreatePollRequest{
		Key: regources.NewKeyInt64(id, regources.REQUEST_DETAILS_CREATE_POLL),
		Attributes: regources.CreatePollRequestAttributes{
			PollData: regources.PollData{
				Type: details.PollData.Type,
			},
			StartTime:                details.StartTime,
			EndTime:                  details.EndTime,
			NumberOfChoices:          details.NumberOfChoices,
			PermissionType:           details.PermissionType,
			VoteConfirmationRequired: details.VoteConfirmationRequired,
			CreatorDetails:           details.CreatorDetails,
		},
		Relationships: regources.CreatePollRequestRelationships{
			ResultProvider: NewAccountKey(details.ResultProviderID).AsRelation(),
		},
	}
}

func newKYCRecoveryRequest(id int64, details history2.KYCRecoveryRequest) *regources.KycRecoveryRequest {
	signersData := make([]regources.UpdateSignerData, 0, len(details.SignersData))
	for _, signer := range details.SignersData {
		signersData = append(signersData, regources.UpdateSignerData{
			Details:  signer.Details,
			RoleId:   signer.RoleID,
			Identity: signer.Identity,
			Weight:   signer.Weight,
		})
	}
	return &regources.KycRecoveryRequest{

		Key: regources.NewKeyInt64(id, regources.REQUEST_DETAILS_KYC_RECOVERY),
		Attributes: regources.KycRecoveryRequestAttributes{
			SignersData:    signersData,
			CreatorDetails: details.CreatorDetails,
			SequenceNumber: details.SequenceNumber,
		},
		Relationships: regources.KycRecoveryRequestRelationships{
			TargetAccount: NewAccountKey(details.TargetAccount).AsRelation(),
		},
	}
}

func newManageOfferRequest(id int64, details history2.ManageOfferRequest) *regources.ManageOfferRequest {
	return &regources.ManageOfferRequest{
		Key: regources.NewKeyInt64(id, regources.REQUEST_DETAILS_MANAGE_OFFER),
		Attributes: regources.ManageOfferRequestAttributes{
			BaseAmount:  details.Amount,
			Fee:         details.Fee,
			IsBuy:       details.IsBuy,
			OfferId:     details.OfferID,
			OrderBookId: details.OrderBookID,
			Price:       details.Price,
		},
	}
}

func newCreatePaymentRequest(id int64, details history2.CreatePaymentRequest) *regources.CreatePaymentRequest {
	return &regources.CreatePaymentRequest{
		Key: regources.NewKeyInt64(id, regources.REQUEST_DETAILS_CREATE_PAYMENT),
		Attributes: regources.CreatePaymentRequestAttributes{
			Amount:                  details.Amount,
			SourceFee:               details.SourceFee,
			DestinationFee:          details.DestinationFee,
			SourcePayForDestination: details.SourcePayForDestination,
			Reference:               details.Reference,
			Subject:                 details.Subject,
		},
		Relationships: regources.CreatePaymentRequestRelationships{
			BalanceFrom: NewBalanceKey(details.BalanceFrom).AsRelation(),
		},
	}
}
