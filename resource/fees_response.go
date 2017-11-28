package resource

type FeesResponse struct {
	Fees             map[string][]FeeEntry `json:"fees"`
}
