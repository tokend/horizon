package reviewablerequest2

type LimitsUpdateRequest struct {
	DocumentHash string                 `json:"document_hash"`
	Details      map[string]interface{} `json:"details"`
}
