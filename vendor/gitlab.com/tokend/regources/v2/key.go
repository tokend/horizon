package regources

import "strconv"

type ResourceType string

const (
	TypeAccounts           ResourceType = "accounts"
	TypeBalances                        = "balances"
	TypeAssets                          = "assets"
	TypeAswapBid                        = "aswap-bid"
	TypeAssetPairs                      = "asset-pairs"
	TypeBalancesState                   = "balances-state"
	TypeExternalSystemID                = "external-system-ids"
	TypeFeeRules                        = "fee-rules"
	TypeKeyValueEntries                 = "key-value-entries"
	TypeLimits                          = "limits"
	TypeOffers                          = "offers"
	TypeOrderBookEntries                = "order-book-entries"
	TypeAccountRoles                    = "account-roles"
	TypeAccountRules                    = "account-rules"
	TypeSales                           = "sales"
	TypeSigners                         = "signers"
	TypeSignerRoles                     = "signer-roles"
	TypeSignerRules                     = "signer-rules"
	TypeParticipantEffects              = "participant-effects"
	TypeOperations                      = "operations"
	TypeQuoteAssets                     = "quote-assets"
	TypeSaleQuoteAssets                 = "sale-quote-assets"
	TypeFees                            = "fees"
	TypeCalculatedFee                   = "calculated-fee"
	TypeTxs                             = "transactions"
	// TypeEffectsFunded - balance received funds from other balance
	TypeEffectsFunded = "effects-funded"
	// TypeEffectsIssued - funds have been issued to the balance
	TypeEffectsIssued = "effects-issued"
	// TypeEffectsCharged - balance has been charged
	TypeEffectsCharged = "effects-charged"
	// TypeEffectsWithdrawn - balance has been charged and corresponding amount of tokens has been destroyed
	TypeEffectsWithdrawn = "effects-withdrawn"
	// TypeEffectsLocked - funds has been locked on the balance
	TypeEffectsLocked = "effects-locked"
	// TypeEffectsUnlocked - funds has been unlocked on the balance
	TypeEffectsUnlocked = "effects-unlocked"
	// TypeEffectsChargedFromLocked - funds has been charged from locked amount on balance
	TypeEffectsChargedFromLocked = "effects-charged-from-locked"
	// TypeEffectsMatched - balance has been charged or received funds due to match of the offers
	TypeEffectsMatched = "effects-matched"
	// TypeCreateAccount - details of createAccountOp
	TypeCreateAccount                          = "operations-create-account"
	TypePayment                                = "operations-payment"
	TypeCreateIssuanceRequest                  = "operations-create-issuance-request"
	TypeSetFees                                = "operations-set-fees"
	TypeCreateWithdrawalRequest                = "operations-create-withdrawal-request"
	TypeManageBalance                          = "operations-manage-balance"
	TypeManageAsset                            = "operations-manage-asset"
	TypeCreatePreissuanceRequest               = "operations-create-preissuance-request"
	TypeManageLimits                           = "operations-manage-limits"
	TypeDirectDebit                            = "operations-direct-debit"
	TypeManageAssetPair                        = "operations-manage-asset-pair"
	TypeManageOffer                            = "operations-manage-offer"
	TypeManageInvoiceRequest                   = "operations-manage-invoice-request"
	TypeReviewRequest                          = "operations-review-request"
	TypeCreateSaleRequest                      = "operations-create-sale-request"
	TypeCheckSaleState                         = "operations-check-sale-state"
	TypeCreateAmlAlert                         = "operations-create-aml-alert"
	TypeCreateChangeRoleRequest                = "operations-create-change-role-request"
	TypePaymentV2                              = "operations-payment-v2"
	TypeManageExternalSystemAccountIDPoolEntry = "operations-manage-external-system-account-id-pool-entry"
	TypeBindExternalSystemAccountID            = "operations-bind-external-system-account-id"
	TypeManageSale                             = "operations-manage-sale"
	TypeManageKeyValue                         = "operations-manage-key-value"
	TypeCreateManageLimitsRequest              = "operations-create-manage-limits-request"
	TypeManageContractRequest                  = "operations-manage-contract-request"
	TypeManageContract                         = "operations-manage-contract"
	TypeCancelSaleRequest                      = "operations-cancel-sale-request"
	TypePayout                                 = "operations-payout"
	TypeManageAccountRole                      = "operations-manage-account-role"
	TypeManageAccountRule                      = "operations-manage-account-rule"
	TypeCreateAswapBidRequest                  = "operations-create-aswap-bid-request"
	TypeCancelAswapBid                         = "operations-cancel-aswap-bid"
	TypeCreateAswapRequest                     = "operations-create-aswap-request"

	TypeRequests                        = "requests"
	TypeRequestDetailsAMLAlert          = "request-details-aml-alert"
	TypeRequestDetailsAssetCreate       = "request-details-asset-create"
	TypeRequestDetailsAssetUpdate       = "request-details-asset-update"
	TypeRequestDetailsAtomicSwap        = "request-details-atomic-swap"
	TypeRequestDetailsAtomicSwapBid     = "request-details-aswap-bid"
	TypeRequestDetailsIssuance          = "request-details-issuance"
	TypeRequestDetailsLimitsUpdate      = "request-details-limits-update"
	TypeRequestDetailsPreIssuance       = "request-details-pre-issuance"
	TypeRequestDetailsSale              = "request-details-sale"
	TypeRequestDetailsChangeRole        = "request-details-change-role"
	TypeRequestDetailsUpdateSaleDetails = "request-details-update-sale-details"
	TypeRequestDetailsUpdateSaleEndTime = "request-details-update-sale-end-time"
	TypeRequestDetailsWithdrawal        = "request-details-withdrawal"
)

// Key - identifier of the Resource
type Key struct {
	ID   string       `json:"id"`
	Type ResourceType `json:"type"`
}

func NewKeyInt64(id int64, resourceType ResourceType) Key {
	return Key{
		ID:   strconv.FormatInt(id, 10),
		Type: resourceType,
	}
}

//GetKey - returns key of the Resource
func (r *Key) GetKey() Key {
	return *r
}

//GetKeyP - returns key pointer
func (r Key) GetKeyP() *Key {
	return &r
}

// AsRelation - converts key to relation
func (r Key) AsRelation() *Relation {
	return &Relation{
		Data: r.GetKeyP(),
	}
}
