package codes

var messages = map[string]string{
	"":                                                  "",
	"tx_success":                                        "Success",
	"tx_failed":                                         "Transaction failed",
	"tx_too_early":                                      "Too early",
	"tx_too_late":                                       "Too late",
	"tx_missing_operation":                              "Missing Operation",
	"tx_bad_auth":                                       "Bad auth",
	"tx_no_account":                                     "No source account",
	"tx_bad_auth_extra":                                 "Bad auth extra",
	"tx_internal_error":                                 "Internal error",
	"tx_account_blocked":                                "Account blocked",
	"tx_duplication":                                    "Transaction duplication",
	"tx_insufficient_fee":                               "The total fee amount is greater than the max fee amount specified by tx source",
	"tx_source_underfunded":                             "Not enough funds for tx fee",
	"tx_commission_line_full":                           "Charging tx fee would cause uint64 overflow",
	"op_inner":                                          "Op inner",
	"op_bad_auth":                                       "You dont have permission to complete this action",
	"op_no_account":                                     "Target account does not exist",
	"op_account_blocked":                                "Operations from blocked account are not allowed",
	"op_no_counterparty":                                "No counterparty",
	"op_counterparty_blocked":                           "Counterparty account is blocked",
	"op_counterparty_wrong_type":                        "Counterparty has wrong account type",
	"op_bad_auth_extra":                                 "Bad auth extra",
	"op_account_type_mismatched":                        "Wrong account type in operation. Refresh page and try again",
	"op_type_not_allowed":                               "Type of account you trying to create is not allowed",
	"op_name_duplication":                               "Name duplication",
	"op_referrer_not_found":                             "Referrer not found",
	"op_invalid_account_version":                        "Invalid package version",
	"op_invoice_not_found":                              "Invoice not found",
	"op_invoice_wrong_amount":                           "Amount must be a positive number",
	"op_invoice_balance_mismatch":                       "This account id has no such balance",
	"op_invoice_account_mismatch":                       "This account id has no such balance",
	"op_invoice_already_paid":                           "Invoice have already been paid",
	"op_too_many_signers":                               "Signers limit is exceeded",
	"op_threshold_out_of_range":                         "Threshold out of range",
	"op_bad_signer":                                     "Invalid signer",
	"op_trust_malformed":                                "Trust malformed",
	"op_trust_too_many":                                 "Trust many",
	"op_invalid_signer_version":                         "Invalid package version",
	"op_invalid_fee_type":                               "Invalid fee type",
	"op_malformed_range":                                "Invalid range",
	"op_range_overlap":                                  "Range you entered overlapped with another one. Delete or reduce an old one before creating new",
	"op_sub_type_not_exist":                             "Subtype not exist",
	"op_not_allowed":                                    "Not allowed",
	"op_type_mismatch":                                  "Type mismatch",
	"op_invalid_amount":                                 "Invalid amount",
	"op_balance_mismatch":                               "Token asset of balance and token asset of operation are not equal",
	"op_reviewer_not_found":                             "Reviewer not found",
	"op_invalid_details":                                "Invalid details",
	"op_fee_mismatch":                                   "Fees mismatched",
	"op_old_signer_not_found":                           "Old signer not found",
	"op_signer_already_exists":                          "Signer already exist",
	"op_destination_not_found":                          "Destination not found",
	"op_request_not_found":                              "Request not found",
	"op_asset_already_exists":                           "Token already exists",
	"op_invalid_max_issuance_amount":                    "Invalid max issuance amount",
	"op_invalid_code":                                   "Invalid token code",
	"op_invalid_name":                                   "Invalid token name",
	"op_request_already_exists":                         "This request already exists in the system",
	"op_stats_asset_already_exists":                     "It can be only one stats asset in the system",
	"op_line_full":                                      "Payment would cause a destination account to exceed their declared trust limit for the token being sent",
	"op_fee_mismatched":                                 "Fees mismatched",
	"op_balance_account_mismatched":                     "Account id has no such balance",
	"op_balance_assets_mismatched":                      "Token asset of balance and token asset of operation are not equal",
	"op_src_balance_not_found":                          "Source balance not found",
	"op_reference_duplication":                          "You cannot make two issuances with the same reference",
	"op_stats_overflow":                                 "Overflow during statistics calculation",
	"op_limits_exceeded":                                "Limits exceeded",
	"op_not_allowed_by_asset_policy":                    "This action is not allowed by token policy",
	"op_no_trust":                                       "No trust",
	"op_already_exists":                                 "Entry already exists",
	"op_invalid_asset":                                  "Invalid token asset",
	"op_invalid_action":                                 "Invalid action",
	"op_invalid_policies":                               "Invalid policies",
	"op_asset_not_found":                                "Token not found",
	"op_pair_not_traded":                                "Token not tradable",
	"op_underfunded":                                    "Not enough funds. Reduce the amount and try again",
	"op_cross_self":                                     "Current order crosses your existing order",
	"op_offer_overflow":                                 "Failed to create offer",
	"op_asset_pair_not_tradable":                        "Token not tradable",
	"op_physical_price_restriction":                     "Price cannot be lower than physical",
	"op_current_price_restriction":                      "Price cannot be lower than current",
	"op_invalid_percent_fee":                            "Invalid percent fee",
	"op_insufficient_price":                             "Order insufficient price",
	"op_success":                                        "Success",
	"op_malformed":                                      "Operation you are trying to create is malformed in some way",
	"op_balance_not_found":                              "Balance not found",
	"op_invoice_overflow":                               "Failed to create invoice",
	"op_not_found":                                      "Not found",
	"op_too_many_invoices":                              "Too many invoices",
	"op_can_not_delete_in_progress":                     "Cannot delete request while it is progress",
	"op_invalid_external_details":                       "External details are too long",
	"op_asset_is_not_withdrawable":                      "It is not allowed to withdraw specified asset",
	"op_conversion_price_is_not_available":              "Conversion price is not available",
	"op_conversion_overflow":                            "Overflow during conversion",
	"op_converted_amount_mismatched":                    "Specified converted amount does not match calculated",
	"op_balance_lock_overflow":                          "Too much assets are locked in specified balance",
	"op_invalid_universal_amount":                       "Unexpected universal amount value",
	"op_initial_preissued_exceeds_max_issuance":         "Number of tokens available for issuance exceeds max number of tokens to be issued",
	"op_base_asset_or_asset_request_not_found":          "Asset of asset creation request for base asset not found",
	"op_quote_asset_not_found":                          "Quote asset not found",
	"op_start_end_invalid":                              "IO should not end before start",
	"op_invalid_end":                                    "Trying to create IO which already ended",
	"op_invalid_price":                                  "Price can not be 0",
	"op_invalid_cap":                                    "Soft cap should not exceed Hard cap",
	"op_insufficient_max_issuance":                      "Max number of tokens can be issued is not sufficient to fulfill soft cap",
	"op_invalid_asset_pair":                             "One of the assets (base or quote) has invalid code or they are equal",
	"op_request_or_sale_already_exists":                 "IO creation request or IO already exists for specified token",
	"op_not_authorized":                                 "Account not authorized to perform issuance of the asset",
	"op_exceeds_max_issuance_amount":                    "Maximal issuance amount will be exceeded after issuance",
	"op_receiver_full_line":                             "Total funds of receiver will exceed UINT64_MAX after issuance",
	"op_fee_exceeds_amount":                             "Fee is more than amount to issue",
	"op_order_book_does_not_exists":                     "Specified IO does not exists or already closed",
	"op_sale_is_not_started_yet":                        "IO has not been started yet",
	"op_sale_already_ended":                             "IO already ended",
	"op_order_violates_hard_cap":                        "Offer violates hard cap restriction of the IO",
	"op_cant_participate_own_sale":                      "Can not participate in the own IO",
	"op_asset_mismatched":                               "Assets mismatched",
	"op_price_does_not_match":                           "Prices does not match",
	"op_insufficient_preissued":                         "Insufficient amount of tokens available for issuance",
	"op_not_verified_cannot_have_policies":              "Not verified account can not have policies",
	"op_price_is_invalid":                               "Price is invalid in some way",
	"op_update_is_not_allowed":                          "Update is not allowed",
	"op_sale_is_not_active":                             "IO is not active",
	"op_invalid_reason":                                 "Reason must be empty if approving and not empty if rejecting",
	"op_hash_mismatched":                                "Request hash mismatched",
	"op_type_mismatched":                                "Request type mismathed",
	"op_reject_not_allowed":                             "Reject not allowed, use permanent reject",
	"op_asset_does_not_exists":                          "Asset does not exist",
	"op_max_issuance_amount_exceeded":                   "Max issuance amount exceeded",
	"op_insufficient_available_for_issuance_amount":     "Insufficient available for issuance amount",
	"op_full_line":                                      "Can't fund balance - total funds exceed system limit",
	"op_base_asset_does_not_exists":                     "Base asset does not exist",
	"op_hard_cap_will_exceed_max_issuance":              "Hard cap will exceed max issuance",
	"op_insufficient_preissued_for_hard_cap":            "Insufficient amount of tokens available for hard cap",
	"op_external_sys_acc_not_allowed":                   "Op contains external system account ID which should be generated on core level",
	"op_external_sys_id_exists":                         "External system account ID already exists",
	"op_limits_update_request_reference_duplication":    "Such request already exists",
	"op_invalid_pre_confirmation_details":               "Invalid pre confirmation details",
	"op_requires_kyc":                                   "You or your counterpary need to complete KYC to use specified asset",
	"op_requestor_is_blocked":                           "Requestor is blocked",
	"op_version_is_not_supported_yet":                   "Version of this operation is not supported yet",
	"op_balance_already_exists":                         "Balance already exists",
	"op_no_available_id":                                "No available external system account id for binding",
	"op_auto_generated_type_not_allowed":                "Auto generated external system type is not allowed to bind",
	"op_acc_to_update_does_not_exist":                   "Account to update KYC data doesn't exist",
	"op_request_exist":                                  "Request already exists",
	"op_same_acc_type_to_set":                           "Account type and kyc level are the same",
	"op_request_does_not_exist":                         "Request does not exist",
	"op_permanent_reject_not_allowed":                   "Permanent reject not allowed, use reject",
	"op_pending_request_update_not_allowed":             "User not allowed to update UpdateKYCRequest if it isn't rejected",
	"op_not_allowed_to_update_request":                  "Master not allowed to update UpdateKYCRequest",
	"op_invalid_update_kyc_request_data":                "Invalid UpdateKYCRequest data",
	"op_invalid_kyc_data":                               "Invalid KYC data",
	"op_non_zero_tasks_to_remove_not_allowed":           "Non-zero value of tasksToRemove field is not allowed in reject KYC request",
	"op_invalid_fee_version":                            "Version of fee entry is greater than ledger version",
	"op_invalid_fee_asset":                              "Asset code of fee asset is invalid",
	"op_fee_asset_not_allowed":                          "Fee asset not allowed",
	"op_cross_asset_fee_not_allowed":                    "Fee asset on payment fee type can differ from asset iff payment fee subtype is OUTGOING",
	"op_fee_asset_not_found":                            "Fee asset not found",
	"op_asset_pair_not_found":                           "Cannot create cross asset fee entry without existing asset pair",
	"op_invalid_asset_pair_price":                       "Asset pair price is <= 0",
	"op_destination_account_not_found":                  "Destination account not found",
	"op_destination_balance_not_found":                  "Destination balance not found",
	"op_invalid_destination_fee":                        "Destination fee is invalid",
	"op_invalid_destination_fee_asset":                  "Destination fee asset must be the same as source balance asset",
	"op_fee_asset_mismatched":                           "Fee asset from operation not the same as fee asset from database",
	"op_insufficient_fee_amount":                        "Insufficient fee amount",
	"op_balance_to_charge_fee_from_not_found":           "Balance to charge fee from not found",
	"op_payment_amount_is_less_than_dest_fee":           "Payment amount is less than destination fee",
	"op_sale_not_found":                                 "Sale not found",
	"op_invalid_new_details":                            "New sale details is invalid JSON",
	"op_update_details_request_already_exists":          "Update sale details request already exists",
	"op_update_details_request_not_found":               "Update sale details request to amend not found",
	"op_kyc_rule_not_found":                             "KYC rule not found",
	"op_invalid_type":                                   "Invalid key value type at KYC rule",
	"op_source_underfunded":                             "Source account underfunded",
	"op_source_balance_lock_overflow":                   "Overflow while locking amount from source balance",
	"op_requires_verification":                          "You or your counterpary need to be verified to use specified asset",
	"op_invalid_sale_state":                             "Invalid sale state",
	"op_promotion_update_request_invalid_asset_pair":    "One of the assets (base or quote) has invalid code or they are equal",
	"op_promotion_update_request_invalid_price":         "Price can not be 0",
	"op_promotion_update_request_start_end_invalid":     "IO should not end before start",
	"op_promotion_update_request_invalid_cap":           "Soft cap should not exceed Hard cap",
	"op_promotion_update_request_invalid_details":       "Details is invalid JSON",
	"op_promotion_update_request_already_exists":        "PromotionUpdateRequest already exists",
	"op_promotion_update_request_not_found":             "PromotionUpdateRequest not found",
	"op_invalid_sale_new_end_time":                      "New sale end time is before start time or current ledger close time",
	"op_invalid_new_end_time":                           "New end time is before start time or current ledger close time",
	"op_update_end_time_request_already_exists":         "UpdateSaleEndTimeRequest already exists",
	"op_update_end_time_request_not_found":              "UpdateSaleEndTimeRequest already not found",
	"op_payment_v1_no_longer_supported":                 "Use payment v2 to perform any payment",
	"op_not_allowed_to_remove":                          "Only request creator can remove request",
	"op_contract_not_found":                             "There is no opened contract with such id",
	"op_only_contractor_can_attach_invoice_to_contract": "Not allowed to attach invoice to contract",
	"op_sender_account_mismatched":                      "Not allowed to use not customer account in contract invoice",
	"op_invoice_is_approved":                            "Not allowed to remove approved invoice request",
	"op_amount_mismatched":                              "Amount in invoice request and details in review invoice request must be equal",
	"op_destination_balance_mismatched":                 "Balance from invoice request and details in review invoice request must be equal",
	"op_not_allowed_account_destination":                "Not allowed to send to account in review invoice request",
	"op_required_source_pay_for_destination":            "Source must pay fee for destination",
	"op_source_balance_mismatched":                      "Source balance must be equal source balance from response create invoice request",
	"op_invoice_receiver_balance_lock_amount_overflow":  "Receiver balance has to much lock amount",
	"op_invoice_already_approved":                       "Not allowed to approve invoice request second time",
	"op_payment_v2_malformed":                           "Payment v2 malformed in some way in approve invoice request",
	"op_too_many_contracts":                             "Reached contract max count for contractor limit",
	"op_details_too_long":                               "Details to long",
	"op_too_many_contract_details":                      "Reached max contract details count limit",
	"op_dispute_reason_too_long":                        "Dispute reason to long",
	"op_already_confirmed":                              "Not allowed to confirm contract second time",
	"op_invoice_not_approved":                           "Not allowed to confirm contract when all invoices not approved",
	"op_dispute_already_started":                        "Not allowed start dispute second time",
	"op_resolve_dispute_now_allowed":                    "Only escrow can resolve dispute",
	"op_confirm_not_allowed":                            "Only customer and contractor can confirm contract",
	"op_customer_balance_overflow":                      "Customer balances amounts exceed max amount",
	"op_system_tasks_not_allowed":                       "Source is trying to set one of the core flags",
	"op_issuance_tasks_not_found":                       "Issuance tasks have not been provided by the source and don't exist in KeyValue table",
	"op_cannot_create_for_acc_id_and_acc_type":          "Limits cannot be created for account ID and account type simultaneously",
	"op_invalid_limits":                                 "Invalid limits",
	"op_contract_details_too_long":                      "Customer details has exceeded max contract details length",
}

func getMessage(rawCode string) string {
	return messages[rawCode]
}
