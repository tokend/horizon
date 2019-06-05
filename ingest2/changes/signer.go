package changes

import "gitlab.com/tokend/go/xdr"

type accountStateStorage interface {
	SetRecovery(address string, inProgress bool) error
}

type signerHandler struct {
	accountStateStorage accountStateStorage
}

func newSignerHandler(storage accountStateStorage) *signerHandler {
	return &signerHandler{
		accountStateStorage: storage,
	}
}

//We don't care about other causes
func (p *signerHandler) Removed(lc ledgerChange) error {
	op := lc.Operation
	switch op.Body.Type {
	case xdr.OperationTypeInitiateKycRecovery:
		initKycRecovery := op.Body.MustInitiateKycRecoveryOp()
		return p.accountStateStorage.SetRecovery(initKycRecovery.Account.Address(), true)
	}
	return nil
}

//We don't care about other causes
func (p *signerHandler) Created(lc ledgerChange) error {
	op := lc.Operation
	address := lc.LedgerChange.MustRemoved().MustSigner().AccountId.Address()
	switch op.Body.Type {
	case xdr.OperationTypeReviewRequest:
		reviewRequestOp := op.Body.MustReviewRequestOp()
		switch reviewRequestOp.RequestDetails.RequestType {
		case xdr.ReviewableRequestTypeKycRecovery:
			return p.accountStateStorage.SetRecovery(address, false)
		}
	}

	return nil
}
