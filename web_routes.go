package horizon

import (
	"github.com/go-chi/chi"
	"gitlab.com/tokend/horizon/web_v2/handlers"
)

func registerV3Routes(r chi.Router) {
	r.Get("/accounts/{id}", handlers.GetAccount)
	r.Get("/accounts/{id}/signers", handlers.GetAccountSigners)
	r.Get("/accounts/{id}/calculated_fees", handlers.GetCalculatedFees)
	r.Get("/assets/{code}", handlers.GetAsset)
	r.Get("/assets", handlers.GetAssetList)
	r.Get("/balances", handlers.GetBalanceList)
	r.Get("/fees", handlers.GetFeeList)
	r.Get("/history", handlers.GetHistory)
	r.Get("/asset_pairs/{id}", handlers.GetAssetPair)
	r.Get("/asset_pairs", handlers.GetAssetPairList)
	r.Get("/offers/{id}", handlers.GetOffer)
	r.Get(`/offers`, handlers.GetOfferList)
	r.Get("/public_key_entries/{id}", handlers.GetPublicKeyEntry)
	r.Get("/transactions", handlers.GetTransactions)

	// reviewable requests
	r.Get("/requests", handlers.GetRequests)
	r.Get("/requests/{id}", handlers.GetRequests)
	r.Get("/create_asset_requests", handlers.GetCreateAssetRequests)
	r.Get("/create_asset_requests/{id}", handlers.GetCreateAssetRequests)
	r.Get("/create_sale_requests", handlers.GetCreateSaleRequests)
	r.Get("/create_sale_requests/{id}", handlers.GetCreateSaleRequests)
	r.Get("/update_asset_requests", handlers.GetUpdateAssetRequests)
	r.Get("/update_asset_requests/{id}", handlers.GetUpdateAssetRequests)
	r.Get("/create_pre_issuance_requests", handlers.GetCreatePreIssuanceRequests)
	r.Get("/create_pre_issuance_requests/{id}", handlers.GetCreatePreIssuanceRequests)
	r.Get("/create_issuance_requests", handlers.GetCreateIssuanceRequests)
	r.Get("/create_issuance_requests/{id}", handlers.GetCreateIssuanceRequests)
	r.Get("/create_withdraw_requests", handlers.GetCreateWithdrawRequests)
	r.Get("/create_withdraw_requests/{id}", handlers.GetCreateWithdrawRequests)
	r.Get("/update_limits_requests", handlers.GetUpdateLimitsRequests)
	r.Get("/update_limits_requests/{id}", handlers.GetUpdateLimitsRequests)
	r.Get("/create_aml_alert_requests", handlers.GetCreateAmlAlertRequests)
	r.Get("/create_aml_alert_requests/{id}", handlers.GetCreateAmlAlertRequests)
	r.Get("/change_role_requests", handlers.GetChangeRoleRequests)
	r.Get("/change_role_requests/{id}", handlers.GetChangeRoleRequests)
	r.Get("/update_sale_details_requests", handlers.GetUpdateSaleDetailsRequests)
	r.Get("/update_sale_details_requests/{id}", handlers.GetUpdateSaleDetailsRequests)
	r.Get("/create_atomic_swap_bid_requests", handlers.GetCreateAtomicSwapBidRequests)
	r.Get("/create_atomic_swap_bid_requests/{id}", handlers.GetCreateAtomicSwapBidRequests)
	r.Get("/create_atomic_swap_requests", handlers.GetCreateAtomicSwapRequests)
	r.Get("/create_atomic_swap_requests{id}", handlers.GetCreateAtomicSwapRequests)
	r.Get("/create_poll_requests", handlers.GetCreatePollRequests)
	r.Get("/create_poll_requests/{id}", handlers.GetCreatePollRequests)

	r.Get("/key_values", handlers.GetKeyValueList)
	r.Get("/key_values/{key}", handlers.GetKeyValue)

	r.Get("/polls/{id}", handlers.GetPoll)
	r.Get("/polls", handlers.GetPollList)
	r.Get("/polls/{id}/relationships/votes", handlers.GetVoteList)
	r.Get("/polls/{id}/relationships/votes/{voter}", handlers.GetVote)

	r.Get("/sales", handlers.GetSaleList)
	r.Get("/sales/{id}", handlers.GetSale)

	r.Get("/order_book/{id}", handlers.DeprecatedGetOrderBook)

	r.Get("/account_roles/{id}", handlers.GetAccountRole)
	r.Get("/account_roles", handlers.GetAccountRoleList)
	r.Get("/account_rules/{id}", handlers.GetAccountRule)
	r.Get("/account_rules", handlers.GetAccountRuleList)

	r.Get("/signer_roles/{id}", handlers.GetSignerRole)
	r.Get("/signer_roles", handlers.GetSignerRoleList)
	r.Get("/signer_rules/{id}", handlers.GetSignerRule)
	r.Get("/signer_rules", handlers.GetSignerRuleList)
}

func registerV4Routes(r chi.Router) {
	r.Get("/order_book/{base}:{quote}:{order_book_id}", handlers.GetOrderBook)
}
