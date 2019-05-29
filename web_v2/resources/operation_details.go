package resources

import (
	"gitlab.com/distributed_lab/logan/v3"
	"gitlab.com/distributed_lab/logan/v3/errors"
	"gitlab.com/tokend/go/xdr"
	"gitlab.com/tokend/horizon/db2/history2"
	regources "gitlab.com/tokend/regources/generated"
)

type operationDetailsProvider func(op history2.Operation) regources.Resource

var operationDetailsProviders = map[xdr.OperationType]operationDetailsProvider{
	xdr.OperationTypeCreateAccount:                          newCreateAccountOp,
	xdr.OperationTypeCreateIssuanceRequest:                  newCreateIssuanceOp,
	xdr.OperationTypeSetFees:                                newSetFeesOp,
	xdr.OperationTypeCreateWithdrawalRequest:                newCreateWithdrawalRequestOp,
	xdr.OperationTypeManageBalance:                          newManageBalanceOp,
	xdr.OperationTypeManageAsset:                            newManageAssetOp,
	xdr.OperationTypeCreatePreissuanceRequest:               newPreIssuanceRequestOp,
	xdr.OperationTypeManageLimits:                           newManageLimitsOp,
	xdr.OperationTypeManageAssetPair:                        newManageAssetPairOp,
	xdr.OperationTypeManageOffer:                            newManageOfferOp,
	xdr.OperationTypeManageInvoiceRequest:                   newManageInvoiceRequestOp,
	xdr.OperationTypeReviewRequest:                          newReviewRequestOp,
	xdr.OperationTypeCreateSaleRequest:                      newCreateSaleRequestOp,
	xdr.OperationTypeCheckSaleState:                         newCheckSaleStateOp,
	xdr.OperationTypeCreateAmlAlert:                         newCreateAMLAlertRequestOp,
	xdr.OperationTypeCreateChangeRoleRequest:                newChangeRoleRequestOp,
	xdr.OperationTypePayment:                                newPaymentOp,
	xdr.OperationTypeManageExternalSystemAccountIdPoolEntry: newManageExternalSystemPool,
	xdr.OperationTypeBindExternalSystemAccountId:            newBindExternalSystemAccountIDOp,
	xdr.OperationTypeManageSale:                             newManageSaleOp,
	xdr.OperationTypeManageKeyValue:                         newManageKeyValueOp,
	xdr.OperationTypeCreateManageLimitsRequest:              newManageLimitsOp,
	xdr.OperationTypeManageContractRequest:                  newManageContractRequestOp,
	xdr.OperationTypeManageContract:                         newManageContractOp,
	xdr.OperationTypeCancelSaleRequest:                      newCancelSaleRequestOp,
	xdr.OperationTypePayout:                                 newPayoutOp,
	xdr.OperationTypeManageAccountRole:                      newManageAccountRoleOp,
	xdr.OperationTypeManageAccountRule:                      newManageAccountRuleOp,
	xdr.OperationTypeCreateAswapBidRequest:                  newCreateASwapBidRequestOp,
	xdr.OperationTypeCancelAswapBid:                         newCancelASwapBidOp,
	xdr.OperationTypeCreateAswapRequest:                     newCreateASwapRequestOp,
	xdr.OperationTypeManageSigner:                           newManageSignerOp,
	xdr.OperationTypeManageSignerRole:                       newManageSignerRoleOp,
	xdr.OperationTypeManageSignerRule:                       newManageSignerRuleOp,
	xdr.OperationTypeStamp:                                  newStampOp,
	xdr.OperationTypeLicense:                                newLicenseOp,
	xdr.OperationTypeManageCreatePollRequest:                newManageCreatePollRequestOp,
	xdr.OperationTypeManagePoll:                             newManagePollOp,
	xdr.OperationTypeManageVote:                             newManageVoteOp,
	xdr.OperationTypeManageAccountSpecificRule:              newManageAccountSpecificRule,
	xdr.OperationTypeCancelChangeRoleRequest:                newCancelChangeRoleRequest,
	xdr.OperationTypeInitiateKycRecovery:                    newInitiateKYCRecoveryOp,
	xdr.OperationTypeCreateKycRecoveryRequest:               newCreateKYCRecoveryRequestOp,
}

//NewOperationDetails - populates operation details into appropriate resource
func NewOperationDetails(op history2.Operation) regources.Resource {
	if _, ok := operationDetailsProviders[op.Type]; !ok {
		panic(errors.From(errors.New("unexpected operation type"), logan.F{
			"type": op.Type,
		}))
	}

	return operationDetailsProviders[op.Type](op)
}

func newManageCreatePollRequestOp(op history2.Operation) regources.Resource {
	details := op.Details.ManageCreatePollRequest
	manageCreateRequestPollOp := regources.ManageCreatePollRequestOp{
		Key: regources.NewKeyInt64(op.ID, regources.OPERATIONS_MANAGE_CREATE_POLL_REQUEST),
		Attributes: &regources.ManageCreatePollRequestOpAttributes{
			Action: details.Action,
		},
	}

	switch details.Action {
	case xdr.ManageCreatePollRequestActionCreate:
		manageCreateRequestPollOp.Attributes.Create = &regources.CreatePollRequestOp{
			AllTasks:        details.CreateDetails.AllTasks,
			CreatorDetails:  details.CreateDetails.CreatorDetails,
			EndTime:         details.CreateDetails.EndTime,
			NumberOfChoices: uint64(details.CreateDetails.NumberOfChoices),
			PermissionType:  uint64(details.CreateDetails.PermissionType),
			PollData: regources.PollData{
				Type: details.CreateDetails.PollData.Type,
			},
			ResultProviderId:         details.CreateDetails.ResultProviderID,
			StartTime:                details.CreateDetails.StartTime,
			VoteConfirmationRequired: details.CreateDetails.VoteConfirmationRequired,
		}
		manageCreateRequestPollOp.Relationships = &regources.ManageCreatePollRequestOpRelationships{
			Request:        NewRequestKey(details.CreateDetails.RequestDetails.RequestID).AsRelation(),
			ResultProvider: NewAccountKey(details.CreateDetails.ResultProviderID).AsRelation(),
		}
	case xdr.ManageCreatePollRequestActionCancel:
		manageCreateRequestPollOp.Relationships = &regources.ManageCreatePollRequestOpRelationships{
			Request: NewRequestKey(details.CancelDetails.RequestID).AsRelation(),
		}
	default:
		panic(errors.From(errors.New("unexpected poll request action"), logan.F{
			"action": details.Action,
		}))
	}

	return &manageCreateRequestPollOp
}

func newManagePollOp(op history2.Operation) regources.Resource {
	details := op.Details.ManagePoll
	managePollOp := regources.ManagePollOp{
		Key: regources.NewKeyInt64(op.ID, regources.OPERATIONS_MANAGE_POLL),
		Attributes: regources.ManagePollOpAttributes{
			Action: details.Action,
			PollId: details.PollID,
		},
		Relationships: regources.ManagePollOpRelationships{
			Poll: NewPollKey(details.PollID).AsRelation(),
		},
	}

	switch details.Action {
	case xdr.ManagePollActionClose:
		managePollOp.Attributes.Close = &regources.ClosePollOp{
			Details:    regources.Details(details.ClosePoll.Details),
			PollId:     details.PollID,
			PollResult: details.ClosePoll.PollResult,
		}
	case xdr.ManagePollActionUpdateEndTime:
		managePollOp.Attributes.UpdateEndTime = &regources.UpdatePollEndTimeOp{
			NewEndTime: details.UpdatePollEndTime.EndTime,
		}
	case xdr.ManagePollActionCancel:
	default:
		panic(errors.From(errors.New("unexpected manage poll action"), logan.F{
			"action": details.Action,
		}))
	}

	return &managePollOp
}
func newCreateAccountOp(op history2.Operation) regources.Resource {
	return &regources.CreateAccountOp{
		Key: regources.NewKeyInt64(op.ID, regources.OPERATIONS_CREATE_ACCOUNT),
		Relationships: regources.CreateAccountOpRelationships{
			Account: NewAccountKey(op.Details.CreateAccount.AccountAddress).AsRelation(),
			Role:    NewAccountRoleKey(op.Details.CreateAccount.AccountRole).AsRelation(),
		},
	}
}

// newCreateIssuanceOp - creates new instance of CreateIssuanceRequestOp
func newCreateIssuanceOp(op history2.Operation) regources.Resource {
	details := op.Details.CreateIssuanceRequest
	return &regources.CreateIssuanceRequestOp{
		Key: regources.NewKeyInt64(op.ID, regources.OPERATIONS_CREATE_ISSUANCE_REQUEST),
		Attributes: regources.CreateIssuanceRequestOpAttributes{
			Fee:            details.Fee,
			Reference:      details.Reference,
			Amount:         details.Amount,
			CreatorDetails: details.CreatorDetails,
			AllTasks:       details.AllTasks,
		},
		Relationships: regources.CreateIssuanceRequestOpRelationships{
			Asset:           NewAssetKey(details.Asset).AsRelation(),
			ReceiverAccount: NewAccountKey(details.ReceiverAccountAddress).AsRelation(),
			ReceiverBalance: NewBalanceKey(details.ReceiverBalanceAddress).AsRelation(),
			Request:         NewRequestKey(details.RequestDetails.RequestID).AsRelation(),
		},
	}
}

func newSetFeesOp(op history2.Operation) regources.Resource {
	return &regources.SetFeeOp{
		Key:        regources.NewKeyInt64(op.ID, regources.OPERATIONS_SET_FEES),
		Attributes: regources.SetFeeOpAttributes(*op.Details.SetFee),
	}
}

// newCreateWithdrawalRequestOp int6 creates new instance of
func newCreateWithdrawalRequestOp(op history2.Operation) regources.Resource {
	details := op.Details.CreateWithdrawRequest
	return &regources.CreateWithdrawRequestOp{
		Key: regources.NewKeyInt64(op.ID, regources.OPERATIONS_CREATE_WITHDRAWAL_REQUEST),
		Attributes: regources.CreateWithdrawRequestOpAttributes{
			Amount:         details.Amount,
			Fee:            details.Fee,
			CreatorDetails: details.CreatorDetails,
		},
		Relationships: regources.CreateWithdrawRequestOpRelationships{
			Balance: NewBalanceKey(details.BalanceAddress).AsRelation(),
		},
	}
}

// newManageBalanceOp - creates new instance of ManageBalanceOp
func newManageBalanceOp(op history2.Operation) regources.Resource {
	details := op.Details.ManageBalance
	return &regources.ManageBalanceOp{
		Key: regources.NewKeyInt64(op.ID, regources.OPERATIONS_MANAGE_BALANCE),
		Attributes: regources.ManageBalanceOpAttributes{
			Action:         details.Action,
			BalanceAddress: details.BalanceAddress,
		},
		Relationships: regources.ManageBalanceOpRelationships{
			DestinationAccount: NewAccountKey(details.DestinationAccount).AsRelation(),
			Asset:              NewAssetKey(details.Asset).AsRelation(),
		},
	}
}

// newManageAssetOp - creates new instance of ManageAsset
func newManageAssetOp(op history2.Operation) regources.Resource {
	details := op.Details.ManageAsset
	return &regources.ManageAssetOp{
		Key: regources.NewKeyInt64(op.ID, regources.OPERATIONS_MANAGE_ASSET),
		Attributes: regources.ManageAssetOpAttributes{
			AssetCode:         details.AssetCode,
			Action:            details.Action,
			Policies:          details.Policies,
			CreatorDetails:    details.CreatorDetails,
			PreIssuanceSigner: details.PreIssuanceSigner,
			MaxIssuanceAmount: details.MaxIssuanceAmount,
		},
		Relationships: regources.ManageAssetOpRelationships{
			Request: NewRequestKey(details.RequestID).AsRelation(),
		},
	}
}

// newPreIssuanceRequestOp - creates new instance of CreatePreIssuanceRequestOp
func newPreIssuanceRequestOp(op history2.Operation) regources.Resource {
	details := op.Details.CreatePreIssuanceRequest
	return &regources.CreatePreIssuanceRequestOp{
		Key: regources.NewKeyInt64(op.ID, regources.OPERATIONS_CREATE_PREISSUANCE_REQUEST),
		Attributes: regources.CreatePreIssuanceRequestOpAttributes{
			Amount: details.Amount,
		},
		Relationships: regources.CreatePreIssuanceRequestOpRelationships{
			Asset:   NewAssetKey(details.AssetCode).AsRelation(),
			Request: NewRequestKey(details.RequestID).AsRelation(),
		},
	}
}

// newManageLimitsOp - creates new instance of ManageLimitsOp
func newManageLimitsOp(op history2.Operation) regources.Resource {
	details := op.Details.ManageLimits
	result := regources.ManageLimitsOp{
		Key: regources.NewKeyInt64(op.ID, regources.OPERATIONS_MANAGE_LIMITS),
		Attributes: regources.ManageLimitsOpAttributes{
			Action: details.Action,
		},
	}

	switch details.Action {
	case xdr.ManageLimitsActionCreate:
		result.Attributes.Create = newManageLimitsCreationOp(*details.Creation)
	case xdr.ManageLimitsActionRemove:
		result.Attributes.Remove = &regources.ManageLimitsRemovalOp{
			LimitsId: details.Removal.LimitsID,
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
func newReviewRequestOp(op history2.Operation) regources.Resource {
	details := op.Details.ReviewRequest
	return &regources.ReviewRequestOp{
		Key: regources.NewKeyInt64(op.ID, regources.OPERATIONS_REVIEW_REQUEST),
		Attributes: regources.ReviewRequestOpAttributes{
			Action:          details.Action,
			Reason:          details.Reason,
			RequestHash:     details.RequestHash,
			RequestId:       details.RequestID,
			IsFulfilled:     details.IsFulfilled,
			AddedTasks:      details.AddedTasks,
			RemovedTasks:    details.RemovedTasks,
			ExternalDetails: details.ExternalDetails,
		},
	}
}

// newManageExternalSystemPool - creates new instance of ManageExternalSystemPoolOp
func newManageExternalSystemPool(op history2.Operation) regources.Resource {
	details := op.Details.ManageExternalSystemPool
	result := &regources.ManageExternalSystemAccountIdPoolEntryOp{
		Key: regources.NewKeyInt64(op.ID, regources.OPERATIONS_MANAGE_EXTERNAL_SYSTEM_ACCOUNT_ID_POOL_ENTRY),
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
func newChangeRoleRequestOp(op history2.Operation) regources.Resource {
	details := op.Details.CreateChangeRoleRequest
	return &regources.CreateChangeRoleRequestOp{
		Key: regources.NewKeyInt64(op.ID, regources.OPERATIONS_CREATE_CHANGE_ROLE_REQUEST),
		Attributes: regources.CreateChangeRoleRequestOpAttributes{
			CreatorDetails: details.CreatorDetails,
			AllTasks:       details.AllTasks,
		},
		Relationships: regources.CreateChangeRoleRequestOpRelationships{
			AccountToUpdateRole: NewAccountKey(details.DestinationAccount).AsRelation(),
			Request:             NewRequestKey(details.RequestDetails.RequestID).AsRelation(),
			RoleToSet:           NewAccountRoleKey(details.AccountRoleToSet).AsRelation(),
		},
	}
}

// newCreateIssuanceOpDetails - creates new instance of CreateIssuanceRequestOp
func newCreateIssuanceOpDetails(op history2.Operation) regources.Resource {
	details := op.Details.CreateIssuanceRequest
	return &regources.CreateIssuanceRequestOp{
		Key: regources.NewKeyInt64(op.ID, regources.OPERATIONS_CREATE_ISSUANCE_REQUEST),
		Attributes: regources.CreateIssuanceRequestOpAttributes{
			Fee:            details.Fee,
			Reference:      details.Reference,
			Amount:         details.Amount,
			CreatorDetails: details.CreatorDetails,
			AllTasks:       details.AllTasks,
		},
		Relationships: regources.CreateIssuanceRequestOpRelationships{
			Asset:           NewAssetKey(details.Asset).AsRelation(),
			ReceiverAccount: NewAccountKey(details.ReceiverAccountAddress).AsRelation(),
			ReceiverBalance: NewBalanceKey(details.ReceiverBalanceAddress).AsRelation(),
			Request:         NewRequestKey(details.RequestDetails.RequestID).AsRelation(),
		},
	}
}

// newManageAssetPairOp - creates new instance of ManageAssetPairOp
func newManageAssetPairOp(op history2.Operation) regources.Resource {
	details := op.Details.ManageAssetPair
	return &regources.ManageAssetPairOp{
		Key: regources.NewKeyInt64(op.ID, regources.OPERATIONS_MANAGE_ASSET_PAIR),
		Attributes: regources.ManageAssetPairOpAttributes{
			PhysicalPrice:           details.PhysicalPrice,
			PhysicalPriceCorrection: details.PhysicalPriceCorrection,
			MaxPriceStep:            details.MaxPriceStep,
			Policies:                details.Policies,
		},
		Relationships: regources.ManageAssetPairOpRelationships{
			BaseAsset:  NewAssetKey(details.BaseAsset).AsRelation(),
			QuoteAsset: NewAssetKey(details.QuoteAsset).AsRelation(),
		},
	}
}

// newManageOfferOp - creates new instance of ManageOfferOp
func newManageOfferOp(op history2.Operation) regources.Resource {
	details := op.Details.ManageOffer
	return &regources.ManageOfferOp{
		Key: regources.NewKeyInt64(op.ID, regources.OPERATIONS_MANAGE_OFFER),
		Attributes: regources.ManageOfferOpAttributes{
			OfferId:     details.OfferID,
			OrderBookId: details.OrderBookID,
			BaseAmount:  details.Amount,
			Price:       details.Price,
			IsBuy:       details.IsBuy,
			Fee:         details.Fee,
			IsDeleted:   details.IsDeleted,
		},
		Relationships: regources.ManageOfferOpRelationships{
			BaseAsset:  NewAssetKey(details.BaseAsset).AsRelation(),
			QuoteAsset: NewAssetKey(details.QuoteAsset).AsRelation(),
		},
	}
}

// newCreateSaleRequestOp - creates new instance of CreateSaleRequestOp
func newCreateSaleRequestOp(op history2.Operation) regources.Resource {
	details := op.Details.CreateSaleRequest
	quoteAssets := &regources.RelationCollection{
		Data: make([]regources.Key, 0, len(details.QuoteAssets)),
	}

	for _, quoteAsset := range details.QuoteAssets {
		quoteAssets.Data = append(quoteAssets.Data, newQuoteAssetKey(quoteAsset.Asset))
	}

	return &regources.CreateSaleRequestOp{
		Key: regources.NewKeyInt64(op.ID, regources.OPERATIONS_CREATE_SALE_REQUEST),
		Attributes: regources.CreateSaleRequestOpAttributes{
			StartTime:      details.StartTime,
			EndTime:        details.EndTime,
			SoftCap:        details.SoftCap,
			HardCap:        details.HardCap,
			CreatorDetails: details.CreatorDetails,
		},
		Relationships: regources.CreateSaleRequestOpRelationships{
			QuoteAssets:       quoteAssets,
			Request:           NewRequestKey(details.RequestID).AsRelation(),
			BaseAsset:         NewAssetKey(details.BaseAsset).AsRelation(),
			DefaultQuoteAsset: NewAssetKey(details.DefaultQuoteAsset).AsRelation(),
		},
	}
}

// newCreateAMLAlertRequestOp - creates new instance of CreateAMLAlertRequestOp
func newCreateAMLAlertRequestOp(op history2.Operation) regources.Resource {
	details := op.Details.CreateAMLAlertRequest
	return &regources.CreateAmlAlertRequestOp{
		Key: regources.NewKeyInt64(op.ID, regources.OPERATIONS_CREATE_AML_ALERT),
		Attributes: regources.CreateAmlAlertRequestOpAttributes{
			Amount:         details.Amount,
			CreatorDetails: details.CreatorDetails,
		},
		Relationships: regources.CreateAmlAlertRequestOpRelationships{
			Balance: NewBalanceKey(details.BalanceAddress).AsRelation(),
		},
	}
}

// newCheckSaleStateOp - creates new instance of CheckSaleStateOp
func newCheckSaleStateOp(op history2.Operation) regources.Resource {
	details := op.Details.CheckSaleState
	return &regources.CheckSaleStateOp{
		Key: regources.NewKeyInt64(op.ID, regources.OPERATIONS_CHECK_SALE_STATE),
		Attributes: regources.CheckSaleStateOpAttributes{
			Effect: details.Effect,
		},
		Relationships: regources.CheckSaleStateOpRelationships{
			Sale: NewSaleKey(details.SaleID).AsRelation(),
		},
	}

}

// newPaymentOp - creates new instance of PaymentOp
func newPaymentOp(op history2.Operation) regources.Resource {
	details := op.Details.Payment
	return &regources.PaymentOp{
		Key: regources.NewKeyInt64(op.ID, regources.OPERATIONS_PAYMENT),
		Attributes: regources.PaymentOpAttributes{
			Amount:                  details.Amount,
			SourceFee:               details.SourceFee,
			DestinationFee:          details.DestinationFee,
			SourcePayForDestination: details.SourcePayForDestination,
			Subject:                 details.Subject,
			Reference:               details.Reference,
		},
		Relationships: regources.PaymentOpRelationships{
			AccountFrom: NewAccountKey(details.AccountFrom).AsRelation(),
			AccountTo:   NewAccountKey(details.AccountTo).AsRelation(),
			BalanceFrom: NewBalanceKey(details.BalanceFrom).AsRelation(),
			BalanceTo:   NewBalanceKey(details.BalanceTo).AsRelation(),
			Asset:       NewAssetKey(details.Asset).AsRelation(),
		},
	}

}

// newCreateManageLimitsRequestOp - creates new instance of CreateManageLimitsRequestOp
func newCreateManageLimitsRequestOp(op history2.Operation) regources.Resource {
	details := op.Details.CreateManageLimitsRequest
	return &regources.CreateManageLimitsRequestOp{
		Key: regources.NewKeyInt64(op.ID, regources.OPERATIONS_CREATE_MANAGE_LIMITS_REQUEST),
		Attributes: regources.CreateManageLimitsRequestOpAttributes{
			CreatorDetails: details.CreatorDetails,
		},
		Relationships: regources.CreateManageLimitsRequestOpRelationships{
			Request: NewRequestKey(details.RequestID).AsRelation(),
		},
	}
}

// newPayoutOp - creates new instance of PayoutOp
func newPayoutOp(op history2.Operation) regources.Resource {
	details := op.Details.Payout
	return &regources.PayoutOp{
		Key: regources.NewKeyInt64(op.ID, regources.OPERATIONS_PAYOUT),
		Attributes: regources.PayoutOpAttributes{
			MaxPayoutAmount:      details.MaxPayoutAmount,
			MinAssetHolderAmount: details.MinAssetHolderAmount,
			MinPayoutAmount:      details.MinPayoutAmount,
			ExpectedFee:          details.ExpectedFee,
			ActualFee:            details.ActualFee,
			ActualPayoutAmount:   details.ActualPayoutAmount,
		},
		Relationships: regources.PayoutOpRelationships{
			SourceAccount: NewAccountKey(details.SourceAccountAddress).AsRelation(),
			SourceBalance: NewBalanceKey(details.SourceBalanceAddress).AsRelation(),
			Asset:         NewAssetKey(details.Asset).AsRelation(),
		},
	}
}

func newLicenseOp(op history2.Operation) regources.Resource {
	details := op.Details.License
	return &regources.LicenseOp{
		Key: regources.NewKeyInt64(op.ID, regources.OPERATIONS_LICENSE),
		Attributes: regources.LicenseOpAttributes{
			PrevLicenseHash: details.PrevLicenseHash,
			LedgerHash:      details.LedgerHash,
			DueDate:         details.DueDate,
			AdminCount:      details.AdminCount,
			Signatures:      details.Signatures,
		},
	}
}

func newStampOp(op history2.Operation) regources.Resource {
	details := op.Details.Stamp
	return &regources.StampOp{
		Key: regources.NewKeyInt64(op.ID, regources.OPERATIONS_STAMP),
		Attributes: regources.StampOpAttributes{
			LedgerHash:  details.LedgerHash,
			LicenseHash: details.LicenseHash,
		},
	}
}

func newManageAccountSpecificRule(op history2.Operation) regources.Resource {
	details := op.Details.ManageAccountSpecificRule
	manageAccSpecificRuleOp := &regources.ManageAccountSpecificRuleOp{
		Key: regources.NewKeyInt64(op.ID, regources.OPERATIONS_MANAGE_ACCOUNT_SPECIFIC_RULE),
		Attributes: &regources.ManageAccountSpecificRuleOpAttributes{
			Action: details.Action,
		},
		Relationships: &regources.ManageAccountSpecificRuleOpRelationships{
			Rule: regources.NewKeyInt64(int64(details.RuleID), regources.ACCOUNT_SPECIFIC_RULES).AsRelation(),
		},
	}
	switch details.Action {
	case xdr.ManageAccountSpecificRuleActionCreate:
		manageAccSpecificRuleOp.Attributes.CreateData = &regources.CreateAccountSpecificRuleData{
			Forbids:   details.CreateDetails.Forbids,
			LedgerKey: details.CreateDetails.LedgerKey,
		}
		if details.CreateDetails.AccountID != "" {
			manageAccSpecificRuleOp.Attributes.CreateData.AccountId = &details.CreateDetails.AccountID
		}
	case xdr.ManageAccountSpecificRuleActionRemove:
	}
	return manageAccSpecificRuleOp
}

func newCancelChangeRoleRequest(op history2.Operation) regources.Resource {
	key := regources.NewKeyInt64(op.ID, regources.OPERATIONS_CANCEL_CHANGE_ROLE_REQUEST)
	return &key
}

func newInitiateKYCRecoveryOp(op history2.Operation) regources.Resource {
	details := op.Details.InitiateKYCRecovery
	return &regources.InitiateKycRecoveryOp{
		Key: regources.NewKeyInt64(op.ID, regources.OPERATIONS_INITIATE_KYC_RECOVERY),
		Attributes: &regources.InitiateKycRecoveryOpAttributes{
			Signer: details.Signer,
		},
		Relationships: &regources.InitiateKycRecoveryOpRelationships{
			Account: NewAccountKey(details.Account).AsRelation(),
		},
	}
}

func newCreateKYCRecoveryRequestOp(op history2.Operation) regources.Resource {
	details := op.Details.CreateKYCRecoveryRequest
	return &regources.CreateKycRecoveryRequestOp{
		Key: regources.NewKeyInt64(op.ID, regources.OPERATIONS_CREATE_KYC_RECOVERY_REQUEST),
		Attributes: &regources.CreateKycRecoveryRequestOpAttributes{
			AllTasks:       details.AllTasks,
			SignersData:    details.SignersData,
			CreatorDetails: details.CreatorDetails,
		},
		Relationships: &regources.CreateKycRecoveryRequestOpRelationships{
			TargetAccount: NewAccountKey(details.TargetAccount).AsRelation(),
			Request:       NewRequestKey(int64(details.RequestDetails.RequestID)).AsRelation(),
		},
	}
}

func newManageVoteOp(op history2.Operation) regources.Resource {
	details := op.Details.ManageVote
	manageVoteOp := regources.ManageVoteOp{
		Key: regources.NewKeyInt64(op.ID, regources.OPERATIONS_MANAGE_VOTE),
		Attributes: regources.ManageVoteOpAttributes{
			Action: details.Action,
		},
		Relationships: regources.ManageVoteOpRelationships{
			Poll: NewPollKey(details.PollID).AsRelation(),
		},
	}

	switch details.Action {
	case xdr.ManageVoteActionCreate:
		choice := uint64(details.VoteData.Single.Choice)
		manageVoteOp.Attributes.Create = &regources.CreateVoteOp{
			PollId: details.PollID,
		}
		if details.VoteData != nil {
			manageVoteOp.Attributes.Create.VoteData = regources.VoteData{
				PollType:     details.VoteData.PollType,
				SingleChoice: &choice,
			}
		}
	case xdr.ManageVoteActionRemove:
		manageVoteOp.Attributes.Remove = &regources.RemoveVoteOp{
			PollId: details.PollID,
		}
	default:
		panic(errors.From(errors.New("unexpected manage vote action"), logan.F{
			"action": details.Action,
		}))
	}

	return &manageVoteOp
}

func newBindExternalSystemAccountIDOp(op history2.Operation) regources.Resource {
	return &regources.BindExternalSystemAccountIdOp{
		Key:        regources.NewKeyInt64(op.ID, regources.OPERATIONS_BIND_EXTERNAL_SYSTEM_ACCOUNT_ID),
		Attributes: regources.BindExternalSystemAccountIdOpAttributes(*op.Details.BindExternalSystemAccount),
	}
}

func newManageInvoiceRequestOp(op history2.Operation) regources.Resource {
	return regources.NewKeyInt64(op.ID, regources.OPERATIONS_MANAGE_INVOICE).GetKeyP()
}

func newManageSaleOp(op history2.Operation) regources.Resource {
	return &regources.ManageSaleOp{
		Key:        regources.NewKeyInt64(op.ID, regources.OPERATIONS_MANAGE_SALE),
		Attributes: regources.ManageSaleOpAttributes(*op.Details.ManageSale),
	}
}

func newManageKeyValueOp(op history2.Operation) regources.Resource {
	return &regources.ManageKeyValueOp{
		Key:        regources.NewKeyInt64(op.ID, regources.OPERATIONS_MANAGE_KEY_VALUE),
		Attributes: regources.ManageKeyValueOpAttributes(*op.Details.ManageKeyValue),
	}
}

func newManageContractRequestOp(op history2.Operation) regources.Resource {
	return regources.NewKeyInt64(op.ID, regources.OPERATIONS_MANAGE_CONTRACT_REQUEST).GetKeyP()
}

func newManageContractOp(op history2.Operation) regources.Resource {
	return regources.NewKeyInt64(op.ID, regources.OPERATIONS_MANAGE_CONTRACT).GetKeyP()
}

func newCreateASwapBidRequestOp(op history2.Operation) regources.Resource {
	return regources.NewKeyInt64(op.ID, regources.OPERATIONS_CREATE_ASWAP_BID_REQUEST).GetKeyP()
}

func newCreateASwapRequestOp(op history2.Operation) regources.Resource {
	return regources.NewKeyInt64(op.ID, regources.OPERATIONS_CREATE_ASWAP_REQUEST).GetKeyP()
}

func newCancelASwapBidOp(op history2.Operation) regources.Resource {
	return regources.NewKeyInt64(op.ID, regources.OPERATIONS_CANCEL_ASWAP_BID).GetKeyP()
}

func newCancelSaleRequestOp(op history2.Operation) regources.Resource {
	return regources.NewKeyInt64(op.ID, regources.OPERATIONS_CANCEL_SALE_REQUEST).GetKeyP()
}
