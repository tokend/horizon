package reviewablerequest

import (
	"gitlab.com/tokend/go/xdr"
	"gitlab.com/tokend/regources/reviewablerequest2"
)

// Populate - populates requestType from xdr.ReviewableRequestType
func PopulateRequestType(rawType xdr.ReviewableRequestType) reviewablerequest2.RequestType {
	return reviewablerequest2.RequestType{
		RequestTypeI: int32(rawType),
		RequestType:  rawType.ShortString(),
	}
}
