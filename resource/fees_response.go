package resource

type FeesResponse struct {
	StorageFeePeriod int64                 `json:"storage_fee_period"`
	PayoutPeriod     int64                 `json:"payout_period"`
	Fees             map[string][]FeeEntry `json:"fees"`
}
