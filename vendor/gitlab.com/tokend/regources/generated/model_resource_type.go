/*
 * GENERATED. Do not modify. Your changes might be overwritten!
 */

package regources

type ResourceType string

// List of ResourceType
const (
	ACCOUNTS                                                ResourceType = "accounts"
	ACCOUNT_ROLES                                           ResourceType = "account-roles"
	ACCOUNT_RULES                                           ResourceType = "account-rules"
	ACCOUNT_SPECIFIC_RULES                                  ResourceType = "account-specific-rules"
	ASSETS                                                  ResourceType = "assets"
	ASSET_PAIRS                                             ResourceType = "asset-pairs"
	ATOMIC_SWAP_ASK                                         ResourceType = "atomic-swap-ask"
	BALANCES                                                ResourceType = "balances"
	BALANCES_STATE                                          ResourceType = "balances-state"
	OPERATIONS_BIND_EXTERNAL_SYSTEM_ACCOUNT_ID              ResourceType = "operations-bind-external-system-account-id"
	CALCULATED_FEE                                          ResourceType = "calculated-fee"
	OPERATIONS_CANCEL_ATOMIC_SWAP_ASK                       ResourceType = "operations-cancel-atomic-swap-ask"
	OPERATIONS_CANCEL_CHANGE_ROLE_REQUEST                   ResourceType = "operations-cancel-change-role-request"
	OPERATIONS_CANCEL_SALE_REQUEST                          ResourceType = "operations-cancel-sale-request"
	OPERATIONS_CHECK_SALE_STATE                             ResourceType = "operations-check-sale-state"
	CONVERTED_BALANCE_STATES                                ResourceType = "converted-balance-states"
	CONVERTED_BALANCES_COLLECTIONS                          ResourceType = "converted-balances-collections"
	OPERATIONS_CREATE_ACCOUNT                               ResourceType = "operations-create-account"
	OPERATIONS_CREATE_AML_ALERT                             ResourceType = "operations-create-aml-alert"
	OPERATIONS_CREATE_ATOMIC_SWAP_ASK_REQUEST               ResourceType = "operations-create-atomic-swap-ask-request"
	OPERATIONS_CREATE_ATOMIC_SWAP_BID_REQUEST               ResourceType = "operations-create-atomic-swap-bid-request"
	OPERATIONS_CREATE_CHANGE_ROLE_REQUEST                   ResourceType = "operations-create-change-role-request"
	OPERATIONS_CREATE_ISSUANCE_REQUEST                      ResourceType = "operations-create-issuance-request"
	OPERATIONS_CREATE_KYC_RECOVERY_REQUEST                  ResourceType = "operations-create-kyc-recovery-request"
	OPERATIONS_CREATE_MANAGE_LIMITS_REQUEST                 ResourceType = "operations-create-manage-limits-request"
	OPERATIONS_CREATE_PREISSUANCE_REQUEST                   ResourceType = "operations-create-preissuance-request"
	OPERATIONS_CREATE_SALE_REQUEST                          ResourceType = "operations-create-sale-request"
	OPERATIONS_CREATE_WITHDRAWAL_REQUEST                    ResourceType = "operations-create-withdrawal-request"
	EFFECTS_FUNDED                                          ResourceType = "effects-funded"
	EFFECTS_ISSUED                                          ResourceType = "effects-issued"
	EFFECTS_CHARGED                                         ResourceType = "effects-charged"
	EFFECTS_WITHDRAWN                                       ResourceType = "effects-withdrawn"
	EFFECTS_LOCKED                                          ResourceType = "effects-locked"
	EFFECTS_UNLOCKED                                        ResourceType = "effects-unlocked"
	EFFECTS_CHARGED_FROM_LOCKED                             ResourceType = "effects-charged-from-locked"
	EFFECTS_MATCHED                                         ResourceType = "effects-matched"
	EXTERNAL_SYSTEM_IDS                                     ResourceType = "external-system-ids"
	FEES                                                    ResourceType = "fees"
	HORIZON_STATE                                           ResourceType = "horizon-state"
	OPERATIONS_INITIATE_KYC_RECOVERY                        ResourceType = "operations-initiate-kyc-recovery"
	REQUESTS                                                ResourceType = "requests"
	REQUEST_DETAILS_AML_ALERT                               ResourceType = "request-details-aml-alert"
	REQUEST_DETAILS_ASSET_CREATE                            ResourceType = "request-details-asset-create"
	REQUEST_DETAILS_ASSET_UPDATE                            ResourceType = "request-details-asset-update"
	REQUEST_DETAILS_ATOMIC_SWAP_ASK                         ResourceType = "request-details-atomic-swap-ask"
	REQUEST_DETAILS_ATOMIC_SWAP_BID                         ResourceType = "request-details-atomic-swap-bid"
	REQUEST_DETAILS_ISSUANCE                                ResourceType = "request-details-issuance"
	REQUEST_DETAILS_LIMITS_UPDATE                           ResourceType = "request-details-limits-update"
	REQUEST_DETAILS_PRE_ISSUANCE                            ResourceType = "request-details-pre-issuance"
	REQUEST_DETAILS_SALE                                    ResourceType = "request-details-sale"
	REQUEST_DETAILS_CHANGE_ROLE                             ResourceType = "request-details-change-role"
	REQUEST_DETAILS_UPDATE_SALE_DETAILS                     ResourceType = "request-details-update-sale-details"
	REQUEST_DETAILS_CREATE_POLL                             ResourceType = "request-details-create-poll"
	REQUEST_DETAILS_WITHDRAWAL                              ResourceType = "request-details-withdrawal"
	REQUEST_DETAILS_KYC_RECOVERY                            ResourceType = "request-details-kyc-recovery"
	SALES                                                   ResourceType = "sales"
	KEY_VALUE_ENTRIES                                       ResourceType = "key-value-entries"
	LEDGER_ENTRY_CHANGES                                    ResourceType = "ledger-entry-changes"
	OPERATIONS_LICENSE                                      ResourceType = "operations-license"
	LIMITS                                                  ResourceType = "limits"
	LIMITS_WITH_STATS                                       ResourceType = "limits-with-stats"
	OPERATIONS_CREATE_ACCOUNT_ROLE                          ResourceType = "operations-create-account-role"
	OPERATIONS_UPDATE_ACCOUNT_ROLE                          ResourceType = "operations-update-account-role"
	OPERATIONS_REMOVE_ACCOUNT_ROLE                          ResourceType = "operations-remove-account-role"
	OPERATIONS_CREATE_ACCOUNT_RULE                          ResourceType = "operations-create-account-rule"
	OPERATIONS_UPDATE_ACCOUNT_RULE                          ResourceType = "operations-update-account-rule"
	OPERATIONS_REMOVE_ACCOUNT_RULE                          ResourceType = "operations-remove-account-rule"
	OPERATIONS_MANAGE_ACCOUNT_SPECIFIC_RULE                 ResourceType = "operations-manage-account-specific-rule"
	OPERATIONS_MANAGE_ASSET                                 ResourceType = "operations-manage-asset"
	OPERATIONS_MANAGE_ASSET_PAIR                            ResourceType = "operations-manage-asset-pair"
	OPERATIONS_MANAGE_BALANCE                               ResourceType = "operations-manage-balance"
	OPERATIONS_MANAGE_CONTRACT                              ResourceType = "operations-manage-contract"
	OPERATIONS_MANAGE_CONTRACT_REQUEST                      ResourceType = "operations-manage-contract-request"
	OPERATIONS_MANAGE_CREATE_POLL_REQUEST                   ResourceType = "operations-manage-create-poll-request"
	OPERATIONS_MANAGE_EXTERNAL_SYSTEM_ACCOUNT_ID_POOL_ENTRY ResourceType = "operations-manage-external-system-account-id-pool-entry"
	OPERATIONS_MANAGE_INVOICE                               ResourceType = "operations-manage-invoice"
	OPERATIONS_MANAGE_KEY_VALUE                             ResourceType = "operations-manage-key-value"
	OPERATIONS_MANAGE_LIMITS                                ResourceType = "operations-manage-limits"
	OPERATIONS_MANAGE_OFFER                                 ResourceType = "operations-manage-offer"
	OPERATIONS_MANAGE_POLL                                  ResourceType = "operations-manage-poll"
	OPERATIONS_MANAGE_SALE                                  ResourceType = "operations-manage-sale"
	OPERATIONS_CREATE_SIGNER                                ResourceType = "operations-create-signer"
	OPERATIONS_UPDATE_SIGNER                                ResourceType = "operations-update-signer"
	OPERATIONS_REMOVE_SIGNER                                ResourceType = "operations-remove-signer"
	OPERATIONS_CREATE_SIGNER_ROLE                           ResourceType = "operations-create-signer-role"
	OPERATIONS_UPDATE_SIGNER_ROLE                           ResourceType = "operations-update-signer-role"
	OPERATIONS_REMOVE_SIGNER_ROLE                           ResourceType = "operations-remove-signer-role"
	OPERATIONS_CREATE_SIGNER_RULE                           ResourceType = "operations-create-signer-rule"
	OPERATIONS_UPDATE_SIGNER_RULE                           ResourceType = "operations-update-signer-rule"
	OPERATIONS_REMOVE_SIGNER_RULE                           ResourceType = "operations-remove-signer-rule"
	OPERATIONS_MANAGE_VOTE                                  ResourceType = "operations-manage-vote"
	MATCHES                                                 ResourceType = "matches"
	OFFERS                                                  ResourceType = "offers"
	OPERATIONS                                              ResourceType = "operations"
	ORDER_BOOK_ENTRIES                                      ResourceType = "order-book-entries"
	ORDER_BOOKS                                             ResourceType = "order-books"
	PARTICIPANT_EFFECTS                                     ResourceType = "participant-effects"
	OPERATIONS_PAYMENT                                      ResourceType = "operations-payment"
	OPERATIONS_PAYOUT                                       ResourceType = "operations-payout"
	POLLS                                                   ResourceType = "polls"
	POLL_OUTCOME                                            ResourceType = "poll-outcome"
	PUBLIC_KEY_ENTRIES                                      ResourceType = "public-key-entries"
	QUOTE_ASSETS                                            ResourceType = "quote-assets"
	OPERATIONS_REVIEW_REQUEST                               ResourceType = "operations-review-request"
	SALE_PARTICIPATION                                      ResourceType = "sale-participation"
	SALE_QUOTE_ASSETS                                       ResourceType = "sale-quote-assets"
	SALE_WHITELIST                                          ResourceType = "sale-whitelist"
	OPERATIONS_SET_FEES                                     ResourceType = "operations-set-fees"
	SIGNERS                                                 ResourceType = "signers"
	SIGNER_ROLES                                            ResourceType = "signer-roles"
	SIGNER_RULES                                            ResourceType = "signer-rules"
	OPERATIONS_STAMP                                        ResourceType = "operations-stamp"
	STATISTICS                                              ResourceType = "statistics"
	TRANSACTIONS                                            ResourceType = "transactions"
	VOTES                                                   ResourceType = "votes"
)
