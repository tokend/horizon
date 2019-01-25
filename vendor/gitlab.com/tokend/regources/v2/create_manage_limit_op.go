package regources

//CreateManageLimitsRequest - details of corresponding op
type CreateManageLimitsRequest struct {
	Key
	Attributes CreateManageLimitsRequestAttrs `json:"attributes"`
}

//CreateManageLimitsRequestAttrs - details of corresponding op
type CreateManageLimitsRequestAttrs struct {
	Data      Details `json:"data"`
	RequestID int64   `json:"request_id"`
}
