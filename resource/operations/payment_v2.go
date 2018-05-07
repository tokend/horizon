package operations

type PaymentV2 struct {
	Base
	PaymentID                            uint64 `json:"payment_id"`
	From                                 string `json:"from,omitempty"`
	To                                   string `json:"to,omitempty"`
	FromBalance                          string `json:"from_balance,omitempty"`
	ToBalance                            string `json:"to_balance,omitempty"`
	Amount                               string `json:"amount"`
	Asset                                string `json:"asset"`
	SourceActualPaymentFee               string `json:"source_actual_payment_fee"`
	SourceActualPaymentFeeAssetCode      string `json:"source_actual_payment_fee_asset_code"`
	DestinationActualPaymentFee          string `json:"destination_actual_payment_fee"`
	DestinationActualPaymentFeeAssetCode string `json:"destination_actual_payment_fee_asset_code"`
	SourceFixedFee                       string `json:"source_fixed_fee"`
	DestinationFixedFee                  string `json:"destination_fixed_fee"`
	SourcePaysForDest                    bool   `json:"source_pays_for_dest"`
	Subject                              string `json:"subject"`
	Reference                            string `json:"reference"`
	SourceSentUniversal                  string `json:"source_sent_universal"`
}
