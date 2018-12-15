package operations

import (
	"testing"

	"gitlab.com/tokend/go/xdr"
)

func TestAllReviewableRequestsHandled(t *testing.T) {
	reviewRequestOpHandler := newReviewRequestOpHandler(&mockPublicKeyProvider{}, &mockBalanceProvider{})
	for _, requestT := range xdr.ReviewableRequestTypeAll {
		if requestT == xdr.ReviewableRequestTypeNone {
			continue
		}
		_, ok := reviewRequestOpHandler.allRequestHandlers[requestT]
		if !ok {
			t.Fatalf("All reivable requests must be handled. Request type: %s is not handled", requestT)
		}
	}
}
