package reviewablerequest

import (
	"gitlab.com/tokend/go/xdr"
	"gitlab.com/tokend/regources"
)

// Populate - populates requestType from xdr.ReviewableRequestType
func PopulateRequestType(rawType xdr.ReviewableRequestType) regources.RequestType {
	return regources.RequestType{
		RequestTypeI: int32(rawType),
		RequestType:  rawType.ShortString(),
	}
}
