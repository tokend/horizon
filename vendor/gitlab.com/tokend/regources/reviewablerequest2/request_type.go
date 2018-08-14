package reviewablerequest2

// RequestType - provides frontend friendly representation of ReviewableRequestType
type RequestType struct {
	// RequestTypeI  - integer representation of request type
	RequestTypeI int32 `json:"request_type_i"`
	// RequestType  - string representation of request type
	RequestType string `json:"request_type"`
}
