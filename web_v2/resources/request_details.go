package resources

import (
	"gitlab.com/distributed_lab/logan/v3"
	"gitlab.com/distributed_lab/logan/v3/errors"
	"gitlab.com/tokend/go/xdr"
	"gitlab.com/tokend/horizon/db2/history2"
	"gitlab.com/tokend/regources/v2"
)

// NewRequestDetails - returns new instance of reviewable request details
func NewRequestDetails(request history2.ReviewableRequest) regources.Resource {
	switch request.RequestType {
	case xdr.ReviewableRequestTypeUpdateAsset:
		return newAssetUpdateRequest(request.ID, *request.Details.AssetUpdate)
	case xdr.ReviewableRequestTypeCreatePreIssuance:
		return newPreIssuanceRequest(request.ID, *request.Details.PreIssuanceCreate)
	case xdr.ReviewableRequestTypeCreateIssuance:
		return newIssuanceRequest(request.ID, *request.Details.IssuanceCreate)
	case xdr.ReviewableRequestTypeCreateWithdraw:
		return newWithdrawalRequest(request.ID, *request.Details.Withdraw)
	case xdr.ReviewableRequestTypeCreateSale:
		return newSaleRequest(request.ID, *request.Details.Sale)
	case xdr.ReviewableRequestTypeUpdateLimits:
		return newLimitsUpdateRequest(request.ID, *request.Details.LimitsUpdate)
	case xdr.ReviewableRequestTypeCreateAmlAlert:
		return newAmlAlertRequest(request.ID, *request.Details.AmlAlert)
	case xdr.ReviewableRequestTypeChangeRole:
		return newChangeRoleRequest(request.ID, *request.Details.ChangeRole)
	case xdr.ReviewableRequestTypeUpdateSaleDetails:
		return newUpdateSaleDetailsRequest(request.ID, *request.Details.UpdateSaleDetails)
	case xdr.ReviewableRequestTypeCreateAsset:
		return newAssetCreateRequest(request.ID, *request.Details.AssetCreation)
	case xdr.ReviewableRequestTypeCreateAtomicSwapBid:
		return newAtomicSwapBidRequest(request.ID, *request.Details.AtomicSwapBidCreation)
	case xdr.ReviewableRequestTypeCreateAtomicSwap:
		return newAtomicSwapRequest(request.ID, *request.Details.AtomicSwap)
	default:
		panic(errors.From(errors.New("unexpected operation type"), logan.F{
			"type": request.RequestType,
		}))
	}
	return nil
}

func newAmlAlertRequest(id int64, details history2.AmlAlertRequest) *regources.CreateAmlAlertRequest {
	return &regources.CreateAmlAlertRequest{
		Key: regources.NewKeyInt64(id, regources.TypeRequestDetailsAMLAlert),
		Attributes: regources.CreateAmlAlertRequestAttrs{
			Amount:         details.Amount,
			CreatorDetails: details.CreatorDetails,
		},
		Relationships: regources.CreateAmlAlertRequestRelations{
			Balance: NewBalanceKey(details.BalanceID).AsRelation(),
		},
	}
}
func newAssetCreateRequest(id int64, details history2.AssetCreationRequest) *regources.CreateAssetRequest {
	return &regources.CreateAssetRequest{
		Key: regources.NewKeyInt64(id, regources.TypeRequestDetailsAssetCreate),
		Attributes: regources.CreateAssetRequestAttrs{
			Asset:                  details.Asset,
			Policies:               details.Policies,
			PreIssuanceAssetSigner: details.PreIssuedAssetSigner,
			MaxIssuanceAmount:      details.MaxIssuanceAmount,
			InitialPreissuedAmount: details.InitialPreissuedAmount,
			CreatorDetails:         details.CreatorDetails,
		},
	}
}
func newAssetUpdateRequest(id int64, details history2.AssetUpdateRequest) *regources.UpdateAssetRequest {
	return &regources.UpdateAssetRequest{
		Key: regources.NewKeyInt64(id, regources.TypeRequestDetailsAssetUpdate),
		Attributes: regources.AssetUpdateRequestAttrs{
			Policies:       details.Policies,
			CreatorDetails: details.CreatorDetails,
		},
		Relationships: regources.UpdateAssetRequestRelations{
			Asset: NewAssetKey(details.Asset).AsRelation(),
		},
	}
}
func newAtomicSwapRequest(id int64, details history2.AtomicSwap) *regources.CreateAtomicSwapRequest {
	return &regources.CreateAtomicSwapRequest{
		Key: regources.NewKeyInt64(id, regources.TypeRequestDetailsAtomicSwap),
		Attributes: regources.CreateAtomicSwapRequestAttrs{
			BaseAmount: regources.Amount(details.BaseAmount),
		},
		Relationships: regources.CreateAtomicSwapRequestRelations{
			Bid:        regources.NewKeyInt64(int64(details.BidID), regources.TypeAswapBid).AsRelation(),
			QuoteAsset: newQuoteAssetKey(details.QuoteAsset).AsRelation(),
		},
	}
}
func newAtomicSwapBidRequest(id int64, details history2.AtomicSwapBidCreation) *regources.CreateAtomicSwapBidRequest {
	quoteAssets := &regources.RelationCollection{
		Data: make([]regources.Key, 0, len(details.QuoteAssets)),
	}
	for _, quoteAsset := range details.QuoteAssets {
		quoteAssets.Data = append(quoteAssets.Data, newQuoteAssetKey(quoteAsset.Asset))
	}

	return &regources.CreateAtomicSwapBidRequest{
		Key: regources.NewKeyInt64(id, regources.TypeRequestDetailsAtomicSwapBid),
		Attributes: regources.CreateAtomicSwapBidRequestAttrs{
			BaseAmount:     regources.Amount(details.BaseAmount),
			CreatorDetails: details.CreatorDetails,
		},
		Relationships: regources.CreateAtomicSwapBidRequestRelations{
			BaseBalance: NewBalanceKey(details.BaseBalance).AsRelation(),
			QuoteAssets: quoteAssets,
		},
	}
}
func newIssuanceRequest(id int64, details history2.IssuanceRequest) *regources.CreateIssuanceRequest {
	return &regources.CreateIssuanceRequest{
		Key: regources.NewKeyInt64(id, regources.TypeRequestDetailsIssuance),
		Attributes: regources.CreateIssuanceRequestAttrs{
			Amount:         details.Amount,
			CreatorDetails: details.CreatorDetails,
		},
		Relationships: regources.CreateIssuanceRequestRelations{
			Asset:    NewAssetKey(details.Asset).AsRelation(),
			Receiver: NewBalanceKey(details.Receiver).AsRelation(),
		},
	}
}
func newLimitsUpdateRequest(id int64, details history2.LimitsUpdateRequest) *regources.UpdateLimitsRequest {
	return &regources.UpdateLimitsRequest{
		Key: regources.NewKeyInt64(id, regources.TypeRequestDetailsLimitsUpdate),
		Attributes: regources.UpdateLimitsRequestAttrs{
			DocumentHash:   details.DocumentHash,
			CreatorDetails: details.CreatorDetails,
		},
	}
}
func newPreIssuanceRequest(id int64, details history2.PreIssuanceRequest) *regources.CreatePreIssuanceRequest {
	return &regources.CreatePreIssuanceRequest{
		Key: regources.NewKeyInt64(id, regources.TypeRequestDetailsPreIssuance),
		Attributes: regources.CreatePreIssuanceRequestAttrs{
			Amount:    details.Amount,
			Signature: details.Signature,
			Reference: details.Reference,
		},
		Relationships: regources.CreatePreIssuanceRequestRelations{
			Asset: NewAssetKey(details.Asset).AsRelation(),
		},
	}
}
func newSaleRequest(id int64, details history2.SaleRequest) *regources.CreateSaleRequest {
	quoteAssets := &regources.RelationCollection{
		Data: make([]regources.Key, 0, len(details.QuoteAssets)),
	}
	for _, quoteAsset := range details.QuoteAssets {
		quoteAssets.Data = append(quoteAssets.Data, newQuoteAssetKey(quoteAsset.Asset))
	}

	return &regources.CreateSaleRequest{
		Key: regources.NewKeyInt64(id, regources.TypeRequestDetailsSale),
		Attributes: regources.CreateSaleRequestAttrs{
			BaseAssetForHardCap: details.BaseAssetForHardCap,
			StartTime:           details.StartTime,
			EndTime:             details.EndTime,
			SaleType:            details.SaleType,
			CreatorDetails:      details.CreatorDetails,
		},
		Relationships: regources.CreateSaleRequestRelations{
			QuoteAssets:       quoteAssets,
			BaseAsset:         NewAssetKey(details.BaseAsset).AsRelation(),
			DefaultQuoteAsset: newQuoteAssetKey(details.DefaultQuoteAsset).AsRelation(),
		},
	}
}
func newChangeRoleRequest(id int64, details history2.ChangeRoleRequest) *regources.ChangeRoleRequest {
	return &regources.ChangeRoleRequest{
		Key: regources.NewKeyInt64(id, regources.TypeRequestDetailsChangeRole),
		Attributes: regources.ChangeRoleRequestAttrs{
			AccountRoleToSet: details.AccountRoleToSet,
			KYCData:          details.KYCData,
			SequenceNumber:   details.SequenceNumber,
		},
		Relationships: regources.ChangeRoleRequestRelations{
			AccountToUpdateRole: NewAccountKey(details.DestinationAccount).AsRelation(),
		},
	}
}
func newUpdateSaleDetailsRequest(id int64, details history2.UpdateSaleDetailsRequest) *regources.UpdateSaleDetailsRequest {
	return &regources.UpdateSaleDetailsRequest{
		Key: regources.NewKeyInt64(id, regources.TypeRequestDetailsUpdateSaleDetails),
		Attributes: regources.UpdateSaleDetailsRequestAttrs{
			CreatorDetails: details.CreatorDetails,
		},
		Relationships: regources.UpdateSaleDetailsRequestRelations{
			Sale: NewSaleKey(int64(details.SaleID)).AsRelation(),
		},
	}
}
func newWithdrawalRequest(id int64, details history2.WithdrawalRequest) *regources.CreateWithdrawalRequest {
	return &regources.CreateWithdrawalRequest{
		Key: regources.NewKeyInt64(id, regources.TypeRequestDetailsWithdrawal),
		Attributes: regources.CreateWithdrawalRequestAttrs{
			Fee: regources.FeeStr{
				Fixed:             details.FixedFee,
				CalculatedPercent: details.PercentFee,
			},
			Amount:          details.Amount,
			CreatorDetails:  details.CreatorDetails,
			ReviewerDetails: details.ReviewerDetails,
		},
		Relationships: regources.CreateWithdrawalRequestRelations{
			Balance: NewBalanceKey(details.BalanceID).AsRelation(),
		},
	}
}
