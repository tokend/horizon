//Package codes is a helper package to help convert to transaction and operation result codes
//to strings used in horizon.
package codes

import (
	"gitlab.com/swarmfund/go/xdr"
	"github.com/go-errors/errors"
	"gitlab.com/distributed_lab/logan"
)

// ErrUnknownCode is returned when an unexepcted value is provided to `String`
var ErrUnknownCode = errors.New("Unknown result code")

type shortStr interface {
	ShortString() string
}

//String returns the appropriate string representation of the provided result code
func String(rawCode interface{}) (string, error) {
	code, ok := rawCode.(shortStr)
	if !ok {
		return "", ErrUnknownCode
	}

	return code.ShortString(), nil
}

// ForOperationResult returns the strong represtation used by horizon for the
// error code `opr`
func ForOperationResult(opr xdr.OperationResult) (string, error) {
	if opr.Code != xdr.OperationResultCodeOpInner {
		return String(opr.Code)
	}

	ir := opr.MustTr()
	var ic interface{}

	switch ir.Type {
	case xdr.OperationTypeCreateAccount:
		ic = ir.MustCreateAccountResult().Code
	case xdr.OperationTypePayment:
		ic = ir.MustPaymentResult().Code
	case xdr.OperationTypeSetOptions:
		ic = ir.MustSetOptionsResult().Code
	case xdr.OperationTypeSetFees:
		ic = ir.MustSetFeesResult().Code
	case xdr.OperationTypeManageAccount:
		ic = ir.MustManageAccountResult().Code
	case xdr.OperationTypeManageForfeitRequest:
		ic = ir.MustManageForfeitRequestResult().Code
	case xdr.OperationTypeRecover:
		ic = ir.MustRecoverResult().Code
	case xdr.OperationTypeManageBalance:
		ic = ir.MustManageBalanceResult().Code
	case xdr.OperationTypeReviewPaymentRequest:
		ic = ir.MustReviewPaymentRequestResult().Code
	case xdr.OperationTypeManageAsset:
		ic = ir.MustManageAssetResult().Code
	case xdr.OperationTypeSetLimits:
		ic = ir.MustSetLimitsResult().Code
	case xdr.OperationTypeDirectDebit:
		ic = ir.MustDirectDebitResult().Code
	case xdr.OperationTypeManageAssetPair:
		ic = ir.MustManageAssetPairResult().Code
	case xdr.OperationTypeManageOffer:
		ic = ir.MustManageOfferResult().Code
	case xdr.OperationTypeManageInvoice:
		ic = ir.MustManageInvoiceResult().Code
	case xdr.OperationTypeReviewRequest:
		ic = ir.MustReviewRequestResult().Code
	case xdr.OperationTypeCreatePreissuanceRequest:
		ic = ir.MustCreatePreIssuanceRequestResult().Code
	case xdr.OperationTypeCreateIssuanceRequest:
		ic = ir.MustCreateIssuanceRequestResult().Code
	}

	return String(ic)
}

func ForTxResult(txResult xdr.TransactionResult) (txResultCode string, opResultCodes []string, err error) {
	txResultCode, err = String(txResult.Result.Code)
	if err != nil {
		err = logan.Wrap(err, "Failed to convert to string tx result code")
		return
	}

	if txResult.Result.Results == nil {
		return
	}

	opResults := txResult.Result.MustResults()
	opResultCodes = make([]string, len(opResults))
	for i := range opResults {
		opResultCodes[i], err = ForOperationResult(opResults[i])
		if err != nil {
			err = logan.Wrap(err, "Failed to convert to string op result")
			return
		}
	}

	return
}
