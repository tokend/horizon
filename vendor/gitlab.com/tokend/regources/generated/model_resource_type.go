/*
 * GENERATED. Do not modify. Your changes might be overwritten!
 */

package resources

type ResourceType string

// List of ResourceType
const (
	ACCOUNTS                                                ResourceType = "accounts"
	BALANCES                                                ResourceType = "balances"
	ASSETS                                                  ResourceType = "assets"
	ASWAP_BID                                               ResourceType = "aswap-bid"
	ASSET_PAIRS                                             ResourceType = "asset-pairs"
	BALANCES_STATE                                          ResourceType = "balances-state"
	EXTERNAL_SYSTEM_IDS                                     ResourceType = "external-system-ids"
	FEE_RULES                                               ResourceType = "fee-rules"
	KEY_VALUE_ENTRIES                                       ResourceType = "key-value-entries"
	LIMITS                                                  ResourceType = "limits"
	LEDGER_ENTRY_CHANGES                                    ResourceType = "ledger-entry-changes"
	OFFERS                                                  ResourceType = "offers"
	ORDER_BOOKS                                             ResourceType = "order-books"
	ORDER_BOOK_ENTRIES                                      ResourceType = "order-book-entries"
	ACCOUNT_ROLES                                           ResourceType = "account-roles"
	ACCOUNT_RULES                                           ResourceType = "account-rules"
	SALES                                                   ResourceType = "sales"
	SIGNERS                                                 ResourceType = "signers"
	SIGNER_ROLES                                            ResourceType = "signer-roles"
	SIGNER_RULES                                            ResourceType = "signer-rules"
	PUBLIC_KEY_ENTRIES                                      ResourceType = "public-key-entries"
	PARTICIPANT_EFFECTS                                     ResourceType = "participant-effects"
	OPERATIONS                                              ResourceType = "operations"
	QUOTE_ASSETS                                            ResourceType = "quote-assets"
	SALE_QUOTE_ASSETS                                       ResourceType = "sale-quote-assets"
	FEES                                                    ResourceType = "fees"
	CALCULATED_FEE                                          ResourceType = "calculated-fee"
	POLLS                                                   ResourceType = "polls"
	POLLS_PARTICIPATION                                     ResourceType = "polls-participation"
	VOTES                                                   ResourceType = "votes"
	TRANSACTIONS                                            ResourceType = "transactions"
	EFFECTS_FUNDED                                          ResourceType = "effects-funded"
	EFFECTS_ISSUED                                          ResourceType = "effects-issued"
	EFFECTS_CHARGED                                         ResourceType = "effects-charged"
	EFFECTS_WITHDRAWN                                       ResourceType = "effects-withdrawn"
	EFFECTS_LOCKED                                          ResourceType = "effects-locked"
	EFFECTS_UNLOCKED                                        ResourceType = "effects-unlocked"
	EFFECTS_CHARGED_FROM_LOCKED                             ResourceType = "effects-charged-from-locked"
	EFFECTS_MATCHED                                         ResourceType = "effects-matched"
	OPERATIONS_CREATE_ACCOUNT                               ResourceType = "operations-create-account"
	OPERATIONS_CREATE_ISSUANCE_REQUEST                      ResourceType = "operations-create-issuance-request"
	OPERATIONS_SET_FEES                                     ResourceType = "operations-set-fees"
	OPERATIONS_CREATE_WITHDRAWAL_REQUEST                    ResourceType = "operations-create-withdrawal-request"
	OPERATIONS_MANAGE_BALANCE                               ResourceType = "operations-manage-balance"
	OPERATIONS_MANAGE_ASSET                                 ResourceType = "operations-manage-asset"
	OPERATIONS_CREATE_PREISSUANCE_REQUEST                   ResourceType = "operations-create-preissuance-request"
	OPERATIONS_MANAGE_LIMITS                                ResourceType = "operations-manage-limits"
	OPERATIONS_MANAGE_ASSET_PAIR                            ResourceType = "operations-manage-asset-pair"
	OPERATIONS_MANAGE_OFFER                                 ResourceType = "operations-manage-offer"
	OPERATIONS_MANAGE_INVOICE_REQUEST                       ResourceType = "operations-manage-invoice-request"
	OPERATIONS_REVIEW_REQUEST                               ResourceType = "operations-review-request"
	OPERATIONS_CREATE_SALE_REQUEST                          ResourceType = "operations-create-sale-request"
	OPERATIONS_CHECK_SALE_STATE                             ResourceType = "operations-check-sale-state"
	OPERATIONS_CREATE_AML_ALERT                             ResourceType = "operations-create-aml-alert"
	OPERATIONS_CREATE_CHANGE_ROLE_REQUEST                   ResourceType = "operations-create-change-role-request"
	OPERATIONS_PAYMENT_V2                                   ResourceType = "operations-payment-v2"
	OPERATIONS_MANAGE_EXTERNAL_SYSTEM_ACCOUNT_ID_POOL_ENTRY ResourceType = "operations-manage-external-system-account-id-pool-entry"
	OPERATIONS_BIND_EXTERNAL_SYSTEM_ACCOUNT_ID              ResourceType = "operations-bind-external-system-account-id"
	OPERATIONS_MANAGE_SALE                                  ResourceType = "operations-manage-sale"
	OPERATIONS_MANAGE_KEY_VALUE                             ResourceType = "operations-manage-key-value"
	OPERATIONS_CREATE_MANAGE_LIMITS_REQUEST                 ResourceType = "operations-create-manage-limits-request"
	OPERATIONS_MANAGE_CONTRACT_REQUEST                      ResourceType = "operations-manage-contract-request"
	OPERATIONS_MANAGE_CONTRACT                              ResourceType = "operations-manage-contract"
	OPERATIONS_MANAGE_CREATE_POLL_REQUEST                   ResourceType = "operations-manage-create-poll-request"
	OPERATIONS_MANAGE_POLL                                  ResourceType = "operations-manage-poll"
	OPERATIONS_MANAGE_VOTE                                  ResourceType = "operations-manage-vote"
	OPERATIONS_CANCEL_SALE_REQUEST                          ResourceType = "operations-cancel-sale-request"
	OPERATIONS_PAYOUT                                       ResourceType = "operations-payout"
	OPERATIONS_CREATE_ACCOUNT_ROLE                          ResourceType = "operations-create-account-role"
	OPERATIONS_UPDATE_ACCOUNT_ROLE                          ResourceType = "operations-update-account-role"
	OPERATIONS_REMOVE_ACCOUNT_ROLE                          ResourceType = "operations-remove-account-role"
	OPERATIONS_CREATE_ACCOUNT_RULE                          ResourceType = "operations-create-account-rule"
	OPERATIONS_UPDATE_ACCOUNT_RULE                          ResourceType = "operations-update-account-rule"
	OPERATIONS_REMOVE_ACCOUNT_RULE                          ResourceType = "operations-remove-account-rule"
	OPERATIONS_CREATE_SIGNER_ROLE                           ResourceType = "operations-create-signer-role"
	OPERATIONS_UPDATE_SIGNER_ROLE                           ResourceType = "operations-update-signer-role"
	OPERATIONS_REMOVE_SIGNER_ROLE                           ResourceType = "operations-remove-signer-role"
	OPERATIONS_CREATE_SIGNER_RULE                           ResourceType = "operations-create-signer-rule"
	OPERATIONS_UPDATE_SIGNER_RULE                           ResourceType = "operations-update-signer-rule"
	OPERATIONS_REMOVE_SIGNER_RULE                           ResourceType = "operations-remove-signer-rule"
	OPERATIONS_CREATE_SIGNER                                ResourceType = "operations-create-signer"
	OPERATIONS_UPDATE_SIGNER                                ResourceType = "operations-update-signer"
	OPERATIONS_REMOVE_SIGNER                                ResourceType = "operations-remove-signer"
	OPERATIONS_CREATE_ASWAP_BID_REQUEST                     ResourceType = "operations-create-aswap-bid-request"
	OPERATIONS_CANCEL_ASWAP_BID                             ResourceType = "operations-cancel-aswap-bid"
	OPERATIONS_CREATE_ASWAP_REQUEST                         ResourceType = "operations-create-aswap-request"
	OPERATIONS_STAMP                                        ResourceType = "operations-stamp"
	OPERATIONS_LICENSE                                      ResourceType = "operations-license"
	REQUESTS                                                ResourceType = "requests"
	REQUEST_DETAILS_AML_ALERT                               ResourceType = "request-details-aml-alert"
	REQUEST_DETAILS_ASSET_CREATE                            ResourceType = "request-details-asset-create"
	REQUEST_DETAILS_ASSET_UPDATE                            ResourceType = "request-details-asset-update"
	REQUEST_DETAILS_ATOMIC_SWAP                             ResourceType = "request-details-atomic-swap"
	REQUEST_DETAILS_ASWAP_BID                               ResourceType = "request-details-aswap-bid"
	REQUEST_DETAILS_CREATE_POLL                             ResourceType = "request-details-create-poll"
	REQUEST_DETAILS_ISSUANCE                                ResourceType = "request-details-issuance"
	REQUEST_DETAILS_LIMITS_UPDATE                           ResourceType = "request-details-limits-update"
	REQUEST_DETAILS_PRE_ISSUANCE                            ResourceType = "request-details-pre-issuance"
	REQUEST_DETAILS_SALE                                    ResourceType = "request-details-sale"
	REQUEST_DETAILS_CHANGE_ROLE                             ResourceType = "request-details-change-role"
	REQUEST_DETAILS_UPDATE_SALE_DETAILS                     ResourceType = "request-details-update-sale-details"
	REQUEST_DETAILS_UPDATE_SALE_END_TIME                    ResourceType = "request-details-update-sale-end-time"
	REQUEST_DETAILS_WITHDRAWAL                              ResourceType = "request-details-withdrawal"
)
