package ingest

import (
	"fmt"

	"encoding/hex"

	"gitlab.com/swarmfund/go/amount"
	"gitlab.com/swarmfund/go/xdr"
	"gitlab.com/swarmfund/horizon/db2/history"
	"gitlab.com/swarmfund/horizon/utf8"
)

// operationDetails returns the details regarding the current operation, suitable
// for ingestion into a history_operation row
func (is *Session) operationDetails() history.OperationDetails {
	c := is.Cursor
	source := c.OperationSourceAccount()

	details := history.OperationDetails{}
	details.Type = c.OperationType()

	switch c.OperationType() {
	case xdr.OperationTypeCreateAccount:
		op := c.Operation().Body.MustCreateAccountOp()

		details.CreateAccount = &history.CreateAccountDetails{
			SourceAccount: source.Address(),
			Destination:   op.Destination.Address(),
			AccountType:   int32(op.AccountType),
		}

	case xdr.OperationTypePayment:
		op := c.Operation().Body.MustPaymentOp()
		opResult := c.OperationResult().MustPaymentResult()

		details.Payment = &history.PaymentDetails{
			SourceAccount:       source.Address(),
			Destination:         opResult.PaymentResponse.Destination.Address(),
			SourceBalance:       op.SourceBalanceId.AsString(),
			DestinationBalance:  op.DestinationBalanceId.AsString(),
			Amount:              amount.String(int64(op.Amount)),
			SourceFee:           amount.String(int64(op.FeeData.SourceFee.PaymentFee)),
			DestinationFee:      amount.String(int64(op.FeeData.DestinationFee.PaymentFee)),
			SourceFixedFee:      amount.String(int64(op.FeeData.SourceFee.FixedFee)),
			DestinationFixedFee: amount.String(int64(op.FeeData.DestinationFee.FixedFee)),
			SourcePaysForDest:   op.FeeData.SourcePaysForDest,
			Subject:             string(op.Subject),
			Reference:           string(op.Reference),
			Asset:               string(opResult.PaymentResponse.Asset),
		}

	case xdr.OperationTypeSetOptions:
		op := c.Operation().Body.MustSetOptionsOp()

		details.SetOptions = &history.SetOptionsDetails{
			MasterWeight:   (*uint32)(op.MasterWeight),
			LowThreshold:   (*uint32)(op.LowThreshold),
			MedThreshold:   (*uint32)(op.MedThreshold),
			HighThreshold:  (*uint32)(op.HighThreshold),
			SignerID:       op.Signer.PubKey.Address(),
			SignerWeight:   (uint32)(op.Signer.Weight),
			SignerType:     (uint32)(op.Signer.SignerType),
			SignerIdentify: (uint32)(op.Signer.Identity),
		}

	case xdr.OperationTypeSetFees:
		op := c.Operation().Body.MustSetFeesOp()
		accountID := ""
		if op.Fee != nil {
			if op.Fee.AccountId != nil {
				accountID = op.Fee.AccountId.Address()
			}
		}
		accountType := op.Fee.AccountType

		details.SetFees = &history.SetFeesDetails{
			AssetCode:   string(op.Fee.Asset),
			FixedFee:    amount.String(int64(op.Fee.PercentFee)),
			PercentFee:  amount.String(int64(op.Fee.PercentFee)),
			FeeType:     int64(op.Fee.FeeType),
			AccountId:   accountID,
			AccountType: (*int32)(accountType),
			Subtype:     int64(op.Fee.Subtype),
			LowerBound:  int64(op.Fee.LowerBound),
			UpperBound:  int64(op.Fee.UpperBound),
		}

	case xdr.OperationTypeManageAccount:
		op := c.Operation().Body.MustManageAccountOp()

		details.ManageAccount = &history.ManageAccountDetails{
			Account:              op.Account.Address(),
			BlockReasonsToAdd:    op.BlockReasonsToAdd,
			BlockReasonsToRemove: op.BlockReasonsToRemove,
		}

	case xdr.OperationTypeCreateWithdrawalRequest:

		op := c.Operation().Body.MustCreateWithdrawalRequestOp()
		request := op.Request
		details.CreateWithdrawalRequest = &history.CreateWithdrawalRequestDetails{
			Amount:          amount.StringU(uint64(request.Amount)),
			Balance:         request.Balance.AsString(),
			FeeFixed:        amount.StringU(uint64(request.Fee.Fixed)),
			FeePercent:      amount.StringU(uint64(request.Fee.Percent)),
			ExternalDetails: string(request.ExternalDetails),
			DestAsset:       string(request.Details.AutoConversion.DestAsset),
			DestAmount:      amount.StringU(uint64(request.Details.AutoConversion.ExpectedAmount)),
		}

	case xdr.OperationTypeRecover:
		op := c.Operation().Body.MustRecoverOp()

		details.Recover = &history.RecoverDetails{ //TODO: create method for PublicKey
			Account: op.Account.Address(),
			OldSigner:/*op.OldSigner*/ "",
			NewSigner:/* op.NewSigner*/ "",
		}

	case xdr.OperationTypeManageBalance:
		op := c.Operation().Body.MustManageBalanceOp()

		details.ManageBalance = &history.ManageBalanceDetails{
			Destination: op.Destination.Address(), //ask
			Action:      int32(op.Action),
		}

	case xdr.OperationTypeReviewPaymentRequest:
		op := c.Operation().Body.MustReviewPaymentRequestOp()

		details.ReviewPaymentRequest = &history.ReviewPaymentRequestDetails{
			PaymentId:    int64(op.PaymentId),
			Accept:       op.Accept,
			RejectReason: string(*op.RejectReason),
		}

	case xdr.OperationTypeSetLimits:
		op := c.Operation().Body.MustSetLimitsOp()

		details.SetLimits = &history.SetLimitsDetails{
			AccountType: (*int32)(op.AccountType),
			Account:     op.Account.Address(), //TODO: ask
		}

	case xdr.OperationTypeDirectDebit:
		op := c.Operation().Body.MustDirectDebitOp().PaymentOp
		opResult := c.OperationResult().MustDirectDebitResult().MustSuccess()

		details.DirectDebit = &history.DirectDebitDetails{
			SourceId:              source.Address(),
			DestinationId:         opResult.PaymentResponse.Destination.Address(),
			SourceBalance:         op.SourceBalanceId.AsString(),
			DestinationBalance:    op.DestinationBalanceId.AsString(),
			Amount:                amount.String(int64(op.Amount)),
			SourcePaymentFee:      amount.String(int64(op.FeeData.SourceFee.PaymentFee)),
			DestinationPaymentFee: amount.String(int64(op.FeeData.DestinationFee.PaymentFee)),
			SourceFixedFee:        amount.String(int64(op.FeeData.SourceFee.FixedFee)),
			DestinationFixedFee:   amount.String(int64(op.FeeData.DestinationFee.FixedFee)),
			SourcePaysForDest:     op.FeeData.SourcePaysForDest,
			Subject:               string(op.Subject),
			Reference:             string(op.Reference),
			Asset:                 string(opResult.PaymentResponse.Asset),
		}

	case xdr.OperationTypeManageAssetPair:
		op := c.Operation().Body.MustManageAssetPairOp()

		details.ManageAssetPair = &history.ManageAssetPairDetails{
			BaseAsset:               string(op.Base),
			QuoteAsset:              string(op.Quote),
			PhysicalPrice:           amount.String(int64(op.PhysicalPrice)),
			PhysicalPriceCorrection: amount.String(int64(op.PhysicalPriceCorrection)),
			MaxPriceStep:            amount.String(int64(op.MaxPriceStep)),
			Policies:                int32(op.Policies),
		}

	case xdr.OperationTypeManageOffer:
		op := c.Operation().Body.ManageOfferOp

		details.ManageOffer = &history.ManageOfferDetails{
			IsBuy:     op.IsBuy,
			Amount:    amount.String(int64(op.Amount)),
			Price:     amount.String(int64(op.Price)),
			Fee:       amount.String(int64(op.Fee)),
			OfferId:   uint64(op.OfferId),
			IsDeleted: int64(op.OfferId) != 0,
		}

	case xdr.OperationTypeManageInvoice:
		op := c.Operation().Body.MustManageInvoiceOp()
		opResult := c.OperationResult().MustManageInvoiceResult()

		details.ManageInvoice = &history.ManageInvoiceDetails{
			Amount:          amount.String(int64(op.Amount)),
			ReceiverBalance: op.ReceiverBalance.AsString(),
			SenderId:        op.Sender.Address(),
			InvoiceId:       uint64(opResult.Success.InvoiceId),
			Asset:           string(opResult.Success.Asset),
		}

	case xdr.OperationTypeReviewRequest:
		op := c.Operation().Body.MustReviewRequestOp()

		details.ReviewRequest = &history.ReviewRequestDetails{
			Action:      int32(op.Action),
			Reason:      string(op.Reason),
			RequestHash: hex.EncodeToString(op.RequestHash[:]),
			RequestId:   uint64(op.RequestId),
			RequestType: int32(op.RequestDetails.RequestType),
		}

	case xdr.OperationTypeManageAsset:
		op := c.Operation().Body.MustManageAssetOp()

		details.ManageAsset = &history.ManageAssetDetails{
			RequestId: uint64(op.RequestId),
			Action:    int32(op.Request.Action),
		}

	case xdr.OperationTypeCreatePreissuanceRequest:
		// no details needed
	case xdr.OperationTypeCreateIssuanceRequest:
		op := c.Operation().Body.MustCreateIssuanceRequestOp()
<<<<<<< Updated upstream
		details["reference"] = utf8.Scrub(string(op.Reference))
		details["amount"] = amount.StringU(uint64(op.Request.Amount))
		details["asset"] = string(op.Request.Asset)
		details["balance_id"] = op.Request.Receiver.AsString()
======

		details.CreateIssuanceRequest = &history.CreateIssuanceRequestDetails{
			FeeFixed:   amount.StringU(uint64(op.Request.Fee.Fixed)),
			FeePercent: amount.StringU(uint64(op.Request.Fee.Percent)),
			Reference:  utf8.Scrub(string(op.Reference)),
			Amount:     amount.StringU(uint64(op.Request.Amount)),
			Asset:      string(op.Request.Asset),
			BalanceId:  op.Request.Receiver.AsString(),
		}

>>>>>>> Stashed changes
	case xdr.OperationTypeCreateSaleRequest:
		// no details needed
	default:
		panic(fmt.Errorf("Unknown operation type: %s", c.OperationType()))
	}
	return details
}
