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
	case xdr.ReviewableRequestTypeAssetUpdate:
		return newAssetUpdateRequest(request.ID, *request.Details.AssetUpdate)
	case xdr.ReviewableRequestTypePreIssuanceCreate:
		return newPreIssuanceRequest(request.ID, *request.Details.PreIssuanceCreate)
	case xdr.ReviewableRequestTypeIssuanceCreate:
		return newIssuanceRequest(request.ID, *request.Details.IssuanceCreate)
	case xdr.ReviewableRequestTypeWithdraw:
		return newWithdrawalRequest(request.ID, *request.Details.Withdraw)
	case xdr.ReviewableRequestTypeSale:
		return newSaleRequest(request.ID, *request.Details.Sale)
	case xdr.ReviewableRequestTypeLimitsUpdate:
		return newLimitsUpdateRequest(request.ID, *request.Details.LimitsUpdate)
	case xdr.ReviewableRequestTypeAmlAlert:
		return newAmlAlertRequest(request.ID, *request.Details.AmlAlert)
	case xdr.ReviewableRequestTypeChangeRole:
		return newUpdateKYCRequest(request.ID, *request.Details.ChangeRole)
	case xdr.ReviewableRequestTypeUpdateSaleDetails:
		return newUpdateSaleDetailsRequest(request.ID, *request.Details.UpdateSaleDetails)
	case xdr.ReviewableRequestTypeAssetCreate:
		return newAssetCreateRequest(request.ID, *request.Details.AssetCreation)
	case xdr.ReviewableRequestTypeCreateAtomicSwapBid:
		return newAtomicSwapBidRequest(request.ID, *request.Details.AtomicSwapBidCreation)
	case xdr.ReviewableRequestTypeAtomicSwap:
		return newAtomicSwapRequest(request.ID, *request.Details.AtomicSwap)
	default:
		panic(errors.From(errors.New("unexpected operation type"), logan.F{
			"type": request.RequestType,
		}))
	}
	return nil
}

func newAmlAlertRequest(id int64, details history2.AmlAlertRequest) *regources.AmlAlertRequest {
	return &regources.AmlAlertRequest{
		Key: regources.NewKeyInt64(id, regources.TypeRequestDetailsAMLAlert),
		Attributes: regources.AmlAlertRequestAttrs{
			Amount: details.Amount,
			Reason: details.Reason,
		},
		Relationships: regources.AmlAlertRequestRelations{
			Balance: NewBalanceKey(details.BalanceID).AsRelation(),
		},
	}
}
func newAssetCreateRequest(id int64, details history2.AssetCreationRequest) *regources.AssetCreateRequest {
	return &regources.AssetCreateRequest{
		Key: regources.NewKeyInt64(id, regources.TypeRequestDetailsAssetCreate),
		Attributes: regources.AssetCreateRequestAttrs{
			Asset:                  details.Asset,
			Policies:               details.Policies,
			PreIssuanceAssetSigner: details.PreIssuedAssetSigner,
			MaxIssuanceAmount:      details.MaxIssuanceAmount,
			InitialPreissuedAmount: details.InitialPreissuedAmount,
			Details:                details.Details,
		},
	}
}
func newAssetUpdateRequest(id int64, details history2.AssetUpdateRequest) *regources.AssetUpdateRequest {
	return &regources.AssetUpdateRequest{
		Key: regources.NewKeyInt64(id, regources.TypeRequestDetailsAssetUpdate),
		Attributes: regources.AssetUpdateRequestAttrs{
			Policies: details.Policies,
			Details:  details.Details,
		},
		Relationships: regources.AssetUpdateRequestRelations{
			Asset: NewAssetKey(details.Asset).AsRelation(),
		},
	}
}
func newAtomicSwapRequest(id int64, details history2.AtomicSwap) *regources.AtomicSwapRequest {
	return &regources.AtomicSwapRequest{
		Key: regources.NewKeyInt64(id, regources.TypeRequestDetailsAtomicSwap),
		Attributes: regources.AtomicSwapRequestAttrs{
			BaseAmount: regources.Amount(details.BaseAmount),
		},
		Relationships: regources.AtomicSwapRequestRelations{
			Bid:        regources.NewKeyInt64(int64(details.BidID), regources.TypeAswapBid).AsRelation(),
			QuoteAsset: newQuoteAssetKey(details.QuoteAsset).AsRelation(),
		},
	}
}
func newAtomicSwapBidRequest(id int64, details history2.AtomicSwapBidCreation) *regources.AtomicSwapBidRequest {
	quoteAssets := &regources.RelationCollection{
		Data: make([]regources.Key, 0, len(details.QuoteAssets)),
	}
	for _, quoteAsset := range details.QuoteAssets {
		quoteAssets.Data = append(quoteAssets.Data, newQuoteAssetKey(quoteAsset.Asset))
	}

	return &regources.AtomicSwapBidRequest{
		Key: regources.NewKeyInt64(id, regources.TypeRequestDetailsAtomicSwapBid),
		Attributes: regources.AtomicSwapBidRequestAttrs{
			BaseAmount: regources.Amount(details.BaseAmount),
			Details:    details.Details,
		},
		Relationships: regources.AtomicSwapBidRequestRelations{
			BaseBalance: NewBalanceKey(details.BaseBalance).AsRelation(),
			QuoteAssets: quoteAssets,
		},
	}
}
func newIssuanceRequest(id int64, details history2.IssuanceRequest) *regources.IssuanceRequest {
	return &regources.IssuanceRequest{
		Key: regources.NewKeyInt64(id, regources.TypeRequestDetailsIssuance),
		Attributes: regources.IssuanceRequestAttrs{
			Amount:  details.Amount,
			Details: details.Details,
		},
		Relationships: regources.IssuanceRequestRelations{
			Asset:    NewAssetKey(details.Asset).AsRelation(),
			Receiver: NewBalanceKey(details.Receiver).AsRelation(),
		},
	}
}
func newLimitsUpdateRequest(id int64, details history2.LimitsUpdateRequest) *regources.LimitsUpdateRequest {
	return &regources.LimitsUpdateRequest{
		Key: regources.NewKeyInt64(id, regources.TypeRequestDetailsLimitsUpdate),
		Attributes: regources.LimitsUpdateRequestAttrs{
			DocumentHash: details.DocumentHash,
			Details:      details.Details,
		},
	}
}
func newPreIssuanceRequest(id int64, details history2.PreIssuanceRequest) *regources.PreIssuanceRequest {
	return &regources.PreIssuanceRequest{
		Key: regources.NewKeyInt64(id, regources.TypeRequestDetailsPreIssuance),
		Attributes: regources.PreIssuanceRequestAttrs{
			Amount:    details.Amount,
			Signature: details.Signature,
			Reference: details.Reference,
		},
		Relationships: regources.PreIssuanceRequestRelations{
			Asset: NewAssetKey(details.Asset).AsRelation(),
		},
	}
}
func newSaleRequest(id int64, details history2.SaleRequest) *regources.SaleRequest {
	quoteAssets := &regources.RelationCollection{
		Data: make([]regources.Key, 0, len(details.QuoteAssets)),
	}
	for _, quoteAsset := range details.QuoteAssets {
		quoteAssets.Data = append(quoteAssets.Data, newQuoteAssetKey(quoteAsset.Asset))
	}

	return &regources.SaleRequest{
		Key: regources.NewKeyInt64(id, regources.TypeRequestDetailsSale),
		Attributes: regources.SaleRequestAttrs{
			BaseAssetForHardCap: details.BaseAssetForHardCap,
			StartTime:           details.StartTime,
			EndTime:             details.EndTime,
			SaleType:            details.SaleType,
			Details:             details.Details,
		},
		Relationships: regources.SaleRequestRelations{
			QuoteAssets:       quoteAssets,
			BaseAsset:         NewAssetKey(details.BaseAsset).AsRelation(),
			DefaultQuoteAsset: newQuoteAssetKey(details.DefaultQuoteAsset).AsRelation(),
		},
	}
}
func newUpdateKYCRequest(id int64, details history2.ChangeRoleRequest) *regources.ChangeRoleRequest {
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
			NewDetails: details.NewDetails,
		},
		Relationships: regources.UpdateSaleDetailsRequestRelations{
			Sale: NewSaleKey(int64(details.SaleID)).AsRelation(),
		},
	}
}
func newWithdrawalRequest(id int64, details history2.WithdrawalRequest) *regources.WithdrawalRequest {
	return &regources.WithdrawalRequest{
		Key: regources.NewKeyInt64(id, regources.TypeRequestDetailsWithdrawal),
		Attributes: regources.WithdrawalRequestAttrs{
			Fee: regources.FeeStr{
				Fixed:             details.FixedFee,
				CalculatedPercent: details.PercentFee,
			},
			Amount:          details.Amount,
			Details:         details.Details,
			ReviewerDetails: details.ReviewerDetails,
		},
		Relationships: regources.WithdrawalRequestRelations{
			Balance: NewBalanceKey(details.BalanceID).AsRelation(),
		},
	}
}
