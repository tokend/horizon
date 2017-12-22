package codes

import (
	"testing"

	"gitlab.com/swarmfund/go/xdr"
	. "github.com/smartystreets/goconvey/convey"
)

func checkCode(t *testing.T, code shortStr) {
	opCode := opCodeToString(code)
	So(opCode, ShouldNotBeBlank)
	message := getMessage(opCode)
	if message == "" {
		t.Errorf("Expected message not to be blanck for %s", opCode)
	}
}

func TestCodes(t *testing.T) {
	Convey("TransactionResultCode", t, func() {
		for _, code := range xdr.TransactionResultCodeAll {
			message := getMessage(code.ShortString())
			if message == "" {
				t.Errorf("Expected message not to be blanck for %s", code.ShortString())
			}
		}
	})
	Convey("OperationResultCode", t, func() {
		for _, code := range xdr.OperationResultCodeAll {
			message := getMessage(code.ShortString())
			if message == "" {
				t.Errorf("Expected message not to be blanck for %s", code.ShortString())
			}
		}
	})
	Convey("CreateAccountResultCode", t, func() {
		for _, code := range xdr.CreateAccountResultCodeAll {
			checkCode(t, code)
		}
	})
	Convey("PaymentResultCode", t, func() {
		for _, code := range xdr.PaymentResultCodeAll {
			checkCode(t, code)
		}
	})
	Convey("SetOptionsResultCode", t, func() {
		for _, code := range xdr.SetOptionsResultCodeAll {
			checkCode(t, code)
		}
	})
	Convey("SetFeesResultCode", t, func() {
		for _, code := range xdr.SetFeesResultCodeAll {
			checkCode(t, code)
		}
	})
	Convey("ManageAccountResultCode", t, func() {
		for _, code := range xdr.ManageAccountResultCodeAll {
			checkCode(t, code)
		}
	})
	Convey("CreateWithdrawalRequestResultCode", t, func() {
		for _, code := range xdr.CreateWithdrawalRequestResultCodeAll {
			checkCode(t, code)
		}
	})
	Convey("RecoverResultCode", t, func() {
		for _, code := range xdr.RecoverResultCodeAll {
			checkCode(t, code)
		}
	})
	Convey("ManageBalanceResultCode", t, func() {
		for _, code := range xdr.ManageBalanceResultCodeAll {
			checkCode(t, code)
		}
	})
	Convey("ReviewPaymentRequestResultCode", t, func() {
		for _, code := range xdr.ReviewPaymentRequestResultCodeAll {
			checkCode(t, code)
		}
	})
	Convey("ManageAssetResultCode", t, func() {
		for _, code := range xdr.ManageAssetResultCodeAll {
			checkCode(t, code)
		}
	})

	Convey("SetLimits", t, func() {
		for _, code := range xdr.SetLimitsResultCodeAll {
			checkCode(t, code)
		}
	})

	Convey("DirectDebit", t, func() {
		for _, code := range xdr.DirectDebitResultCodeAll {
			checkCode(t, code)
		}
	})

	Convey("ManageAssetPair", t, func() {
		for _, code := range xdr.ManageAssetPairResultCodeAll {
			checkCode(t, code)
		}
	})

	Convey("ManageOffer", t, func() {
		for _, code := range xdr.ManageOfferResultCodeAll {
			checkCode(t, code)
		}
	})

	Convey("ManageInvoice", t, func() {
		for _, code := range xdr.ManageInvoiceResultCodeAll {
			checkCode(t, code)
		}
	})

	Convey("Withdrawal", t, func() {
		for _, code := range xdr.CreateWithdrawalRequestResultCodeAll {
			checkCode(t, code)
		}
	})

	Convey("Create sale", t, func() {
		for _, code := range xdr.CreateSaleCreationRequestResultCodeAll {
			checkCode(t, code)
		}
	})
}
