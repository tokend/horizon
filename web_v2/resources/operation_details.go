package resources

import (
	"gitlab.com/distributed_lab/logan/v3"
	"gitlab.com/distributed_lab/logan/v3/errors"
	"gitlab.com/tokend/go/xdr"
	"gitlab.com/tokend/horizon/db2/history2"
	"gitlab.com/tokend/regources/v2"
)

//NewOperationDetails - populates operation details into appropriate resource
func NewOperationDetails(op history2.Operation) regources.Resource {
	switch op.Type {
	case xdr.OperationTypeCreateAccount:
		return &regources.CreateAccountOp{
			Key: regources.NewKeyInt64(op.ID, regources.TypeCreateAccount),
			Relationships: regources.CreateAccountOpRelation{
				Account: NewAccountKey(op.Details.CreateAccount.AccountAddress).AsRelation(),
				Role:    NewAccountRoleKey(op.Details.CreateAccount.AccountRole).AsRelation(),
			},
		}
	case xdr.OperationTypeCreateIssuanceRequest:
		return newCreateIssuanceOpDetails(op.ID, *op.Details.CreateIssuanceRequest)
	case xdr.OperationTypeSetFees:
		return &regources.SetFeeOp{
			Key:        regources.NewKeyInt64(op.ID, regources.TypeSetFees),
			Attributes: regources.SetFeeOpAttrs(*op.Details.SetFee),
		}
	case xdr.OperationTypeCreateWithdrawalRequest:
		return newCreateWithdrawalRequestOp(op.ID, *op.Details.CreateWithdrawRequest)
	case xdr.OperationTypeManageBalance:
		return newManageBalanceOp(op.ID, *op.Details.ManageBalance)
	case xdr.OperationTypeManageAsset:
		return newManageAssetOp(op.ID, *op.Details.ManageAsset)
	case xdr.OperationTypeCreatePreissuanceRequest:
		return newPreIssuanceRequestOp(op.ID, *op.Details.CreatePreIssuanceRequest)
	case xdr.OperationTypeManageAssetPair:
		return newManageAssetPairOp(op.ID, *op.Details.ManageAssetPair)
	case xdr.OperationTypeManageOffer:
		return newManageOfferOp(op.ID, *op.Details.ManageOffer)
	case xdr.OperationTypeManageInvoiceRequest:
		return regources.NewKeyInt64(op.ID, regources.TypeManageInvoiceRequest).GetKeyP()
	case xdr.OperationTypeReviewRequest:
		return newReviewRequestOp(op.ID, *op.Details.ReviewRequest)
	case xdr.OperationTypeCreateSaleRequest:
		return newCreateSaleRequestOp(op.ID, *op.Details.CreateSaleRequest)
	case xdr.OperationTypeCheckSaleState:
		return newCheckSaleStateOp(op.ID, *op.Details.CheckSaleState)
	case xdr.OperationTypeCreateAmlAlert:
		return newCreateAMLAlertRequestOp(op.ID, *op.Details.CreateAMLAlertRequest)
	case xdr.OperationTypeCreateChangeRoleRequest:
		return newChangeRoleRequestOp(op.ID, *op.Details.CreateChangeRoleRequest)
	case xdr.OperationTypePayment:
		return newPaymentOp(op.ID, *op.Details.Payment)
	case xdr.OperationTypeManageExternalSystemAccountIdPoolEntry:
		return newManageExternalSystemPool(op.ID, *op.Details.ManageExternalSystemPool)
	case xdr.OperationTypeBindExternalSystemAccountId:
		return &regources.BindExternalSystemAccountOp{
			Key:        regources.NewKeyInt64(op.ID, regources.TypeBindExternalSystemAccountID),
			Attributes: regources.BindExternalSystemAccountOpAttrs(*op.Details.BindExternalSystemAccount),
		}
	case xdr.OperationTypeManageSale:
		return &regources.ManageSaleOp{
			Key:        regources.NewKeyInt64(op.ID, regources.TypeManageSale),
			Attributes: regources.ManageSaleOpAttrs(*op.Details.ManageSale),
		}
	case xdr.OperationTypeManageKeyValue:
		return &regources.ManageKeyValueOp{
			Key:        regources.NewKeyInt64(op.ID, regources.TypeManageKeyValue),
			Attributes: regources.ManageKeyValueOpAttrs(*op.Details.ManageKeyValue),
		}
	case xdr.OperationTypeCreateManageLimitsRequest:
		return newCreateManageLimitsRequestOp(op.ID, *op.Details.CreateManageLimitsRequest)
	case xdr.OperationTypeManageContractRequest:
		return regources.NewKeyInt64(op.ID, regources.TypeManageContractRequest).GetKeyP()
	case xdr.OperationTypeManageContract:
		return regources.NewKeyInt64(op.ID, regources.TypeManageContract).GetKeyP()
	case xdr.OperationTypeCancelSaleRequest:
		return regources.NewKeyInt64(op.ID, regources.TypeCancelSaleRequest).GetKeyP()
	case xdr.OperationTypePayout:
		return newPayoutOp(op.ID, *op.Details.Payout)
	case xdr.OperationTypeManageAccountRole:
		return newManageAccountRole(op.ID, *op.Details.ManageAccountRole)
	case xdr.OperationTypeManageAccountRule:
		return newManageAccountRule(op.ID, *op.Details.ManageAccountRule)
	case xdr.OperationTypeCreateAswapBidRequest:
		return regources.NewKeyInt64(op.ID, regources.TypeCreateAswapBidRequest).GetKeyP()
	case xdr.OperationTypeCancelAswapBid:
		return regources.NewKeyInt64(op.ID, regources.TypeCancelAswapBid).GetKeyP()
	case xdr.OperationTypeCreateAswapRequest:
		return regources.NewKeyInt64(op.ID, regources.TypeCreateAswapRequest).GetKeyP()
	case xdr.OperationTypeManageSignerRole:
		return newManageSignerRole(op.ID, *op.Details.ManageSignerRole)
	case xdr.OperationTypeManageSignerRule:
		return newManageSignerRule(op.ID, *op.Details.ManageSignerRule)
	case xdr.OperationTypeManageSigner:
		return newManageSigner(op.ID, *op.Details.ManageSigner)
	case xdr.OperationTypeManageLimits:
		return newManageLimitsOp(op.ID, *op.Details.ManageLimits)
	case xdr.OperationTypeStamp:
		return newStampOp(op.ID, *op.Details.Stamp)
	case xdr.OperationTypeLicense:
		return newLicenseOp(op.ID, *op.Details.License)
	default:
		panic(errors.From(errors.New("unexpected operation type"), logan.F{
			"type": op.Type,
		}))
	}
}

// newManageLimitsOp - creates new instance of ManageLimitsOp
func newManageLimitsOp(id int64, details history2.ManageLimitsDetails) *regources.ManageLimitsOp {
	result := regources.ManageLimitsOp{
		Key: regources.NewKeyInt64(id, regources.TypeManageLimits),
		Attributes: regources.ManageLimitsOpAttributes{
			Action: details.Action,
		},
	}

	switch details.Action {
	case xdr.ManageLimitsActionCreate:
		result.Attributes.Create = newManageLimitsCreationOp(*details.Creation)
	case xdr.ManageLimitsActionRemove:
		result.Attributes.Remove = &regources.ManageLimitsRemovalOp{
			LimitsID: details.Removal.LimitsID,
		}
	default:
		panic(errors.From(errors.New("unexpected manage limits action"), logan.F{
			"action": details.Action,
		}))
	}

	return &result
}

// newManageLimitsCreationOp - creates new instance of ManageLimitsCreationOp
func newManageLimitsCreationOp(details history2.ManageLimitsCreationDetails) *regources.ManageLimitsCreationOp {
	return &regources.ManageLimitsCreationOp{
		AccountRole:     details.AccountRole,
		AccountAddress:  details.AccountAddress,
		StatsOpType:     details.StatsOpType,
		AssetCode:       details.AssetCode,
		IsConvertNeeded: details.IsConvertNeeded,
		DailyOut:        details.DailyOut,
		WeeklyOut:       details.WeeklyOut,
		AnnualOut:       details.AnnualOut,
		MonthlyOut:      details.MonthlyOut,
	}
}

// newReviewRequestOp - creates new instance of ReviewRequestOp
func newReviewRequestOp(id int64, details history2.ReviewRequestDetails) *regources.ReviewRequestOp {
	return &regources.ReviewRequestOp{
		Key: regources.NewKeyInt64(id, regources.TypeReviewRequest),
		Attributes: regources.ReviewRequestOpAttrs{
			Action:          details.Action,
			Reason:          details.Reason,
			RequestHash:     details.RequestHash,
			RequestID:       details.RequestID,
			IsFulfilled:     details.IsFulfilled,
			AddedTasks:      details.AddedTasks,
			RemovedTasks:    details.RemovedTasks,
			ExternalDetails: details.ExternalDetails,
		},
	}
}

// newManageExternalSystemPool - creates new instance of ManageExternalSystemPoolOp
func newManageExternalSystemPool(id int64, details history2.ManageExternalSystemPoolDetails) *regources.ManageExternalSystemPoolOp {
	result := &regources.ManageExternalSystemPoolOp{
		Key: regources.NewKeyInt64(id, regources.TypeManageExternalSystemAccountIDPoolEntry),
	}

	switch details.Action {
	case xdr.ManageExternalSystemAccountIdPoolEntryActionCreate:
		result.Attributes.Create = new(regources.CreateExternalSystemPoolOp)
		*result.Attributes.Create = regources.CreateExternalSystemPoolOp(*details.Create)
	case xdr.ManageExternalSystemAccountIdPoolEntryActionRemove:
		result.Attributes.Remove = new(regources.RemoveExternalSystemPoolOp)
		*result.Attributes.Remove = regources.RemoveExternalSystemPoolOp(*details.Remove)
	default:
		panic(errors.From(errors.New("unexpected action for manage ex sys id pool"), logan.F{
			"action": details.Action,
		}))
	}

	return result
}

// newChangeRoleRequest - creates new instance of CreateKYCRequestOp
func newChangeRoleRequestOp(id int64, details history2.CreateChangeRoleRequestDetails,
) *regources.CreateChangeRoleRequest {
	return &regources.CreateChangeRoleRequest{
		Key: regources.NewKeyInt64(id, regources.TypeCreateChangeRoleRequest),
		Attributes: regources.CreateChangeRoleRequestAttrs{
			CreatorDetails: details.CreatorDetails,
			AllTasks:       details.AllTasks,
		},
		Relationships: regources.CreateChangeRoleRequestOpRelations{
			AccountToUpdateRole: NewAccountKey(details.DestinationAccount).AsRelation(),
			Request:             NewRequestKey(details.RequestDetails.RequestID).AsRelation(),
			RoleToSet:           NewAccountRoleKey(details.AccountRoleToSet).AsRelation(),
		},
	}
}

// newCreateIssuanceOpDetails - creates new instance of CreateIssuanceRequestOp
func newCreateIssuanceOpDetails(id int64, details history2.CreateIssuanceRequestDetails) *regources.CreateIssuanceRequestOp {
	return &regources.CreateIssuanceRequestOp{
		Key: regources.NewKeyInt64(id, regources.TypeCreateIssuanceRequest),
		Attributes: regources.CreateIssuanceRequestOpAttrs{
			Fee:            details.Fee,
			Reference:      details.Reference,
			Amount:         details.Amount,
			CreatorDetails: details.CreatorDetails,
			AllTasks:       details.AllTasks,
		},
		Relationships: regources.CreateIssuanceRequestOpRelations{
			Asset:           NewAssetKey(details.Asset).AsRelation(),
			ReceiverAccount: NewAccountKey(details.ReceiverAccountAddress).AsRelation(),
			ReceiverBalance: NewBalanceKey(details.ReceiverBalanceAddress).AsRelation(),
			Request:         NewRequestKey(details.RequestDetails.RequestID).AsRelation(),
		},
	}
}

// newCreateWithdrawalRequestOp int6 creates new instance of
func newCreateWithdrawalRequestOp(id int64,
	details history2.CreateWithdrawRequestDetails) *regources.CreateWithdrawRequestOp {
	return &regources.CreateWithdrawRequestOp{
		Key: regources.NewKeyInt64(id, regources.TypeCreateWithdrawalRequest),
		Attributes: regources.CreateWithdrawRequestOpAttrs{
			Amount:         details.Amount,
			Fee:            details.Fee,
			CreatorDetails: details.CreatorDetails,
		},
		Relationships: regources.CreateWithdrawRequestOpRelations{
			Balance: NewBalanceKey(details.BalanceAddress).AsRelation(),
		},
	}
}

// newManageBalanceOp - creates new instance of ManageBalanceOp
func newManageBalanceOp(id int64, details history2.ManageBalanceDetails) *regources.ManageBalanceOp {
	return &regources.ManageBalanceOp{
		Key: regources.NewKeyInt64(id, regources.TypeManageBalance),
		Attributes: regources.ManageBalanceOpAttrs{
			Action:         details.Action,
			BalanceAddress: details.BalanceAddress,
		},
		Relationships: regources.ManageBalanceOpRelations{
			DestinationAccount: NewAccountKey(details.DestinationAccount).AsRelation(),
			Asset:              NewAssetKey(details.Asset).AsRelation(),
		},
	}
}

// newManageAssetOp - creates new instance of ManageAsset
func newManageAssetOp(id int64, details history2.ManageAssetDetails) *regources.ManageAssetOp {
	return &regources.ManageAssetOp{
		Key: regources.NewKeyInt64(id, regources.TypeManageAsset),
		Attributes: regources.ManageAssetOpAttrs{
			AssetCode:         details.AssetCode,
			Action:            details.Action,
			Policies:          details.Policies,
			CreatorDetails:    details.CreatorDetails,
			PreissuedSigner:   details.PreissuedSigner,
			MaxIssuanceAmount: details.MaxIssuanceAmount,
		},
		Relationships: regources.ManageAssetOpRelations{
			Request: NewRequestKey(details.RequestID).AsRelation(),
		},
	}
}

// newPreIssuanceRequestOp - creates new instance of CreatePreIssuanceRequestOp
func newPreIssuanceRequestOp(id int64, details history2.CreatePreIssuanceRequestDetails) *regources.CreatePreIssuanceRequestOp {
	return &regources.CreatePreIssuanceRequestOp{
		Key: regources.NewKeyInt64(id, regources.TypeCreatePreissuanceRequest),
		Attributes: regources.CreatePreIssuanceRequestOpAttrs{
			Amount: details.Amount,
		},
		Relationships: regources.CreatePreIssuanceRequestOpRelations{
			Asset:   NewAssetKey(details.AssetCode).AsRelation(),
			Request: NewRequestKey(details.RequestID).AsRelation(),
		},
	}
}

// newManageAssetPairOp - creates new instance of ManageAssetPairOp
func newManageAssetPairOp(id int64, details history2.ManageAssetPairDetails) *regources.ManageAssetPairOp {
	return &regources.ManageAssetPairOp{
		Key: regources.NewKeyInt64(id, regources.TypeManageAssetPair),
		Attributes: regources.ManageAssetPairOpAttrs{
			PhysicalPrice:           details.PhysicalPrice,
			PhysicalPriceCorrection: details.PhysicalPriceCorrection,
			MaxPriceStep:            details.MaxPriceStep,
			Policies:                details.Policies,
		},
		Relationships: regources.ManageAssetPairOpRelations{
			BaseAsset:  NewAssetKey(details.BaseAsset).AsRelation(),
			QuoteAsset: NewAssetKey(details.QuoteAsset).AsRelation(),
		},
	}
}

// newManageOfferOp - creates new instance of ManageOfferOp
func newManageOfferOp(id int64, details history2.ManageOfferDetails) *regources.ManageOfferOp {
	return &regources.ManageOfferOp{
		Key: regources.NewKeyInt64(id, regources.TypeManageOffer),
		Attributes: regources.ManageOfferOpAttrs{
			OfferID:     details.OfferID,
			OrderBookID: details.OrderBookID,
			BaseAmount:  details.Amount,
			Price:       details.Price,
			IsBuy:       details.IsBuy,
			Fee:         details.Fee,
			IsDeleted:   details.IsDeleted,
		},
		Relationships: regources.ManageOfferOpRelations{
			BaseAsset:  NewAssetKey(details.BaseAsset).AsRelation(),
			QuoteAsset: NewAssetKey(details.QuoteAsset).AsRelation(),
		},
	}
}

// newCreateSaleRequestOp - creates new instance of CreateSaleRequestOp
func newCreateSaleRequestOp(id int64, details history2.CreateSaleRequestDetails) *regources.CreateSaleRequestOp {
	quoteAssets := &regources.RelationCollection{
		Data: make([]regources.Key, 0, len(details.QuoteAssets)),
	}

	for _, quoteAsset := range details.QuoteAssets {
		quoteAssets.Data = append(quoteAssets.Data, newQuoteAssetKey(quoteAsset.Asset))
	}

	return &regources.CreateSaleRequestOp{
		Key: regources.NewKeyInt64(id, regources.TypeCreateSaleRequest),
		Attributes: regources.CreateSaleRequestOpAttrs{
			StartTime:      details.StartTime,
			EndTime:        details.EndTime,
			SoftCap:        details.SoftCap,
			HardCap:        details.HardCap,
			CreatorDetails: details.CreatorDetails,
		},
		Relationships: regources.CreateSaleRequestOpRelations{
			QuoteAssets:       quoteAssets,
			Request:           NewRequestKey(details.RequestID).AsRelation(),
			BaseAsset:         NewAssetKey(details.BaseAsset).AsRelation(),
			DefaultQuoteAsset: NewAssetKey(details.DefaultQuoteAsset).AsRelation(),
		},
	}
}

// newCreateAMLAlertRequestOp - creates new instance of CreateAMLAlertRequestOp
func newCreateAMLAlertRequestOp(id int64, details history2.CreateAMLAlertRequestDetails) *regources.CreateAMLAlertRequestOp {
	return &regources.CreateAMLAlertRequestOp{
		Key: regources.NewKeyInt64(id, regources.TypeCreateAmlAlert),
		Attributes: regources.CreateAMLAlertRequestOpAttrs{
			Amount:         details.Amount,
			CreatorDetails: details.CreatorDetails,
		},
		Relationships: regources.CreateAMLAlertRequestOpRelations{
			Balance: NewBalanceKey(details.BalanceAddress).AsRelation(),
		},
	}
}

// newCheckSaleStateOp - creates new instance of CheckSaleStateOp
func newCheckSaleStateOp(id int64, details history2.CheckSaleStateDetails) *regources.CheckSaleStateOp {
	return &regources.CheckSaleStateOp{
		Key: regources.NewKeyInt64(id, regources.TypeCheckSaleState),
		Attributes: regources.CheckSaleStateOpAttrs{
			Effect: details.Effect,
		},
		Relationships: regources.CheckSaleStateOpRelations{
			Sale: NewSaleKey(details.SaleID).AsRelation(),
		},
	}

}

// newPaymentOp - creates new instance of PaymentOp
func newPaymentOp(id int64, details history2.PaymentDetails) *regources.PaymentOp {
	return &regources.PaymentOp{
		Key: regources.NewKeyInt64(id, regources.TypePaymentV2),
		Attributes: regources.PaymentOpAttrs{
			Amount:                  details.Amount,
			SourceFee:               details.SourceFee,
			DestinationFee:          details.DestinationFee,
			SourcePayForDestination: details.SourcePayForDestination,
			Subject:                 details.Subject,
			Reference:               details.Reference,
		},
		Relationships: regources.PaymentOpRelations{
			AccountFrom: NewAccountKey(details.AccountFrom).AsRelation(),
			AccountTo:   NewAccountKey(details.AccountTo).AsRelation(),
			BalanceFrom: NewBalanceKey(details.BalanceFrom).AsRelation(),
			BalanceTo:   NewBalanceKey(details.BalanceTo).AsRelation(),
			Asset:       NewAssetKey(details.Asset).AsRelation(),
		},
	}

}

// newCreateManageLimitsRequestOp - creates new instance of CreateManageLimitsRequestOp
func newCreateManageLimitsRequestOp(id int64, details history2.CreateManageLimitsRequestDetails) *regources.CreateManageLimitsRequestOp {
	return &regources.CreateManageLimitsRequestOp{
		Key: regources.NewKeyInt64(id, regources.TypeCreateManageLimitsRequest),
		Attributes: regources.CreateManageLimitsRequestOpAttrs{
			CreatorDetails: details.CreatorDetails,
		},
		Relationships: regources.CreateManageLimitsRequestOpRelations{
			Request: NewRequestKey(details.RequestID).AsRelation(),
		},
	}
}

// newPayoutOp - creates new instance of PayoutOp
func newPayoutOp(id int64, details history2.PayoutDetails) *regources.PayoutOp {
	return &regources.PayoutOp{
		Key: regources.NewKeyInt64(id, regources.TypePayout),
		Attributes: regources.PayoutOpAttrs{
			MaxPayoutAmount:      details.MaxPayoutAmount,
			MinAssetHolderAmount: details.MinAssetHolderAmount,
			MinPayoutAmount:      details.MinPayoutAmount,
			ExpectedFee:          details.ExpectedFee,
			ActualFee:            details.ActualFee,
			ActualPayoutAmount:   details.ActualPayoutAmount,
		},
		Relationships: regources.PayoutOpRelations{
			SourceAccount: NewAccountKey(details.SourceAccountAddress).AsRelation(),
			SourceBalance: NewBalanceKey(details.SourceBalanceAddress).AsRelation(),
			Asset:         NewAssetKey(details.Asset).AsRelation(),
		},
	}
}

func newLicenseOp(id int64, details history2.LicenseDetails) *regources.LicenseOp {
	return &regources.LicenseOp{
		Key: regources.NewKeyInt64(id, regources.TypeLicense),
		Attributes: regources.LicenseOpAttrs{
			PrevLicenseHash: details.PrevLicenseHash,
			LedgerHash:      details.LedgerHash,
			DueDate:         details.DueDate,
			AdminCount:      details.AdminCount,
			Signatures:      details.Signatures,
		},
	}
}

func newStampOp(id int64, details history2.StampDetails) *regources.StampOp {
	return &regources.StampOp{
		Key: regources.NewKeyInt64(id, regources.TypeStamp),
		Attributes: regources.StampOpAttributes{
			LedgerHash:  details.LedgerHash,
			LicenseHash: details.LicenseHash,
		},
	}
}
