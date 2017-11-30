package reviewablerequest

import "gitlab.com/swarmfund/go/xdr"

// RequestType - provides frontend friendly representation of ReviewableRequestType
type RequestType struct {
	// RequestTypeI  - integer representation of request type
	RequestTypeI int32 `json:"request_type_i"`
	// RequestType  - string representation of request type
	RequestType string `json:"request_type"`
}

// Populate - populates requestType from xdr.ReviewableRequestType
func (r *RequestType) Populate(rawType xdr.ReviewableRequestType) {
	r.RequestTypeI = int32(rawType)
	r.RequestType = rawType.ShortString()
}
