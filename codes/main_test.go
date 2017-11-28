package codes

import (
	"testing"

	"gitlab.com/tokend/go/xdr"
	. "github.com/smartystreets/goconvey/convey"
)

func TestCodes(t *testing.T) {
	Convey("TransactionResultCode", t, func() {
		for _, code := range xdr.TransactionResultCodeAll {
			result, err := String(code)
			So(err, ShouldBeNil)
			So(result, ShouldNotBeBlank)
		}
	})
	Convey("OperationResultCode", t, func() {
		for _, code := range xdr.OperationResultCodeAll {
			result, err := String(code)
			So(err, ShouldBeNil)
			So(result, ShouldNotBeBlank)
		}
	})
	Convey("CreateAccountResultCode", t, func() {
		for _, code := range xdr.CreateAccountResultCodeAll {
			result, err := String(code)
			So(err, ShouldBeNil)
			So(result, ShouldNotBeBlank)
		}
	})
	Convey("PaymentResultCode", t, func() {
		for _, code := range xdr.PaymentResultCodeAll {
			result, err := String(code)
			So(err, ShouldBeNil)
			So(result, ShouldNotBeBlank)
		}
	})
	Convey("SetOptionsResultCode", t, func() {
		for _, code := range xdr.SetOptionsResultCodeAll {
			result, err := String(code)
			So(err, ShouldBeNil)
			So(result, ShouldNotBeBlank)
		}
	})
	Convey("ManageCoinsEmissionRequestResultCode", t, func() {
		for _, code := range xdr.ManageCoinsEmissionRequestResultCodeAll {
			result, err := String(code)
			So(err, ShouldBeNil)
			So(result, ShouldNotBeBlank)
		}
	})
	Convey("ReviewCoinsEmissionRequestResultCode", t, func() {
		for _, code := range xdr.ReviewCoinsEmissionRequestResultCodeAll {
			result, err := String(code)
			So(err, ShouldBeNil)
			So(result, ShouldNotBeBlank)
		}
	})
	Convey("SetFeesResultCode", t, func() {
		for _, code := range xdr.SetFeesResultCodeAll {
			result, err := String(code)
			So(err, ShouldBeNil)
			So(result, ShouldNotBeBlank)
		}
	})
	Convey("ManageAccountResultCode", t, func() {
		for _, code := range xdr.ManageAccountResultCodeAll {
			result, err := String(code)
			So(err, ShouldBeNil)
			So(result, ShouldNotBeBlank)
		}
	})
	Convey("ManageForfeitRequestResultCode", t, func() {
		for _, code := range xdr.ManageForfeitRequestResultCodeAll {
			result, err := String(code)
			So(err, ShouldBeNil)
			So(result, ShouldNotBeBlank)
		}
	})
	Convey("RecoverResultCode", t, func() {
		for _, code := range xdr.RecoverResultCodeAll {
			result, err := String(code)
			So(err, ShouldBeNil)
			So(result, ShouldNotBeBlank)
		}
	})
	Convey("ManageBalanceResultCode", t, func() {
		for _, code := range xdr.ManageBalanceResultCodeAll {
			result, err := String(code)
			So(err, ShouldBeNil)
			So(result, ShouldNotBeBlank)
		}
	})
	Convey("ReviewPaymentRequestResultCode", t, func() {
		for _, code := range xdr.ReviewPaymentRequestResultCodeAll {
			result, err := String(code)
			So(err, ShouldBeNil)
			So(result, ShouldNotBeBlank)
		}
	})
	Convey("ManageAssetResultCode", t, func() {
		for _, code := range xdr.ManageAssetResultCodeAll {
			result, err := String(code)
			So(err, ShouldBeNil)
			So(result, ShouldNotBeBlank)
		}
	})
	Convey("UploadPreemissionsResultCode", t, func() {
		for _, code := range xdr.UploadPreemissionsResultCodeAll {
			result, err := String(code)
			So(err, ShouldBeNil)
			So(result, ShouldNotBeBlank)
		}
	})

	Convey("SetLimits", t, func() {
		for _, code := range xdr.SetLimitsResultCodeAll {
			result, err := String(code)
			So(err, ShouldBeNil)
			So(result, ShouldNotBeBlank)
		}
	})

	Convey("DirectDebit", t, func() {
		for _, code := range xdr.DirectDebitResultCodeAll {
			result, err := String(code)
			So(err, ShouldBeNil)
			So(result, ShouldNotBeBlank)
		}
	})

	Convey("ManageAssetPair", t, func() {
		for _, code := range xdr.ManageAssetPairResultCodeAll {
			result, err := String(code)
			So(err, ShouldBeNil)
			So(result, ShouldNotBeBlank)
		}
	})

	Convey("ManageOffer", t, func() {
		for _, code := range xdr.ManageOfferResultCodeAll {
			result, err := String(code)
			So(err, ShouldBeNil)
			So(result, ShouldNotBeBlank)
		}
	})

	Convey("ManageInvoice", t, func() {
		for _, code := range xdr.ManageInvoiceResultCodeAll {
			result, err := String(code)
			So(err, ShouldBeNil)
			So(result, ShouldNotBeBlank)
		}
	})
}
