package changes

import (
	"encoding/hex"
	"time"

	"gitlab.com/distributed_lab/logan/v3"
	"gitlab.com/distributed_lab/logan/v3/errors"
	"gitlab.com/tokend/go/amount"
	"gitlab.com/tokend/go/xdr"
	"gitlab.com/tokend/horizon/db2"
	history "gitlab.com/tokend/horizon/db2/history2"
	"gitlab.com/tokend/horizon/ingest2/internal"
	"gitlab.com/tokend/horizon/utf8"
	"gitlab.com/tokend/regources/v2"
)

var errUnknownRemoveReason = errors.New("request was removed due to unknown reason")

type reviewableRequestStorage interface {
	//Inserts Reviewable request into DB
	Insert(request history.ReviewableRequest) error
	//Updates Reviewable request
	Update(request history.ReviewableRequest) error
	//Approves reviewable request
	Approve(id uint64) error
	//PermanentReject - rejects permanently reviewable request
	PermanentReject(id uint64, rejectReason string) error
	//Cancel - cancels reviewable request
	Cancel(id uint64) error
}

type reviewableRequestHandler struct {
	storage reviewableRequestStorage
}

func newReviewableRequestHandler(storage reviewableRequestStorage) *reviewableRequestHandler {
	return &reviewableRequestHandler{
		storage: storage,
	}
}

//Created - handles creation of new reviewable request
func (c *reviewableRequestHandler) Created(lc ledgerChange) error {
	reviewableRequest := lc.LedgerChange.MustCreated().Data.MustReviewableRequest()
	histReviewableReq, err := c.convertReviewableRequest(&reviewableRequest, lc.LedgerCloseTime)
	if err != nil {
		return errors.Wrap(err, "failed to convert reviewable request", logan.F{
			"request":         reviewableRequest,
			"ledger_sequence": lc.LedgerSeq,
		})
	}

	err = c.storage.Insert(*histReviewableReq)
	if err != nil {
		return errors.Wrap(err, "failed to insert reviewable request", logan.F{
			"request":         histReviewableReq,
			"ledger_sequence": lc.LedgerSeq,
		})
	}

	return nil
}

//Updated - handles update of the request due to approve or reject op
func (c *reviewableRequestHandler) Updated(lc ledgerChange) error {
	reviewableRequest := lc.LedgerChange.MustUpdated().Data.MustReviewableRequest()
	histReviewableRequest, err := c.convertReviewableRequest(&reviewableRequest, lc.LedgerCloseTime)
	if err != nil {
		return errors.Wrap(err, "failed to convert reviewable request", logan.F{
			"request":         reviewableRequest,
			"ledger_sequence": lc.LedgerSeq,
		})
	}

	err = c.storage.Update(*histReviewableRequest)
	if err != nil {
		return errors.Wrap(err, "failed to update reviewable request", logan.F{
			"request":         histReviewableRequest,
			"ledger_sequence": lc.LedgerSeq,
		})
	}

	return nil
}

//Removed - handles removal of request from core due to approval, cancellation or permanent reject
func (c *reviewableRequestHandler) Removed(lc ledgerChange) error {
	// The request is deleted in 3 cases:
	// 1. Due to approve via reviewRequestOp
	// 2. Due to permanentReject via reviewRequestOp
	// 3. Due to cancel via specific operation
	op := lc.Operation.Body
	switch op.Type {
	case xdr.OperationTypeReviewRequest:
		return c.removedOnReview(lc)
	case xdr.OperationTypeManageAsset:
		return c.cancel(lc)
	case xdr.OperationTypeManageSale:
		return c.handleRemoveOnManageSale(lc)
	case xdr.OperationTypeCancelAswapBid:
		return c.cancel(lc)
	// auto review is handled by each operation separately
	case xdr.OperationTypeCreateIssuanceRequest,
		xdr.OperationTypeCheckSaleState,
		xdr.OperationTypeCreateWithdrawalRequest,
		xdr.OperationTypeCreatePreissuanceRequest,
		xdr.OperationTypeManageLimits,
		xdr.OperationTypeManageInvoiceRequest,
		xdr.OperationTypeCreateSaleRequest,
		xdr.OperationTypeCreateAmlAlert,
		xdr.OperationTypeCreateChangeRoleRequest,
		xdr.OperationTypeCreateAswapBidRequest:
		return nil
	default: // safeguard for future updates
		return errors.From(errUnknownRemoveReason, logan.F{
			"op_type": op.Type,
		})
	}
}

func (c *reviewableRequestHandler) handleRemoveOnManageAsset(lc ledgerChange) error {
	op := lc.Operation.Body.MustManageAssetOp()
	switch op.Request.Action {
	// must be handled by operation
	case xdr.ManageAssetActionCreateAssetCreationRequest,
		xdr.ManageAssetActionCreateAssetUpdateRequest,
		xdr.ManageAssetActionChangePreissuedAssetSigner,
		xdr.ManageAssetActionUpdateMaxIssuance:
		return nil
	case xdr.ManageAssetActionCancelAssetRequest:
		return c.cancel(lc)
	default:
		return errors.From(errUnknownRemoveReason, logan.F{
			"manage_asset_action": op.Request.Action,
		})
	}
}

func (c *reviewableRequestHandler) handleRemoveOnManageSale(lc ledgerChange) error {
	op := lc.Operation.Body.MustManageSaleOp()
	switch op.Data.Action {
	// must be handled by operation
	case xdr.ManageSaleActionCreateUpdateDetailsRequest:
		return nil
	case xdr.ManageSaleActionCancel:
		return c.cancel(lc)
	default:
		return errors.From(errUnknownRemoveReason, logan.F{
			"manage_sale_action": op.Data.Action,
		})
	}
}

func (c *reviewableRequestHandler) removedOnReview(lc ledgerChange) error {
	key := lc.LedgerChange.MustRemoved().MustReviewableRequest()
	op := lc.Operation.Body.MustReviewRequestOp()
	switch op.Action {
	case xdr.ReviewRequestOpActionApprove:
		err := c.storage.Approve(uint64(key.RequestId))
		if err != nil {
			return errors.Wrap(err, "Failed to delete reviewable request", logan.F{
				"ledger_entry_key": key,
			})
		}
	case xdr.ReviewRequestOpActionPermanentReject:
		err := c.storage.PermanentReject(uint64(key.RequestId), string(op.Reason))
		if err != nil {
			return errors.Wrap(err, "Failed to delete reviewable request", logan.F{
				"ledger_entry_key": key,
			})
		}
	default:
		return errors.From(
			errors.New("unknown action during handle of removed reviewable request on review operation"),
			logan.F{
				"action": op.Action.String(),
			})
	}

	return nil
}

func (c *reviewableRequestHandler) cancel(lc ledgerChange) error {
	key := lc.LedgerChange.MustRemoved().MustReviewableRequest()
	err := c.storage.Cancel(uint64(key.RequestId))
	if err != nil {
		return errors.Wrap(err, "failed to cancel reivewalbe request on remove by operation", logan.F{
			"op_type": lc.Operation.Body.Type,
		})
	}

	return nil
}

func (c *reviewableRequestHandler) convertReviewableRequest(request *xdr.ReviewableRequestEntry,
	ledgerCloseTime time.Time) (*history.ReviewableRequest, error) {

	var reference *string
	if request.Reference != nil {
		reference = new(string)
		*reference = utf8.Scrub(string(*request.Reference))
	}

	details, err := c.getReviewableRequestDetails(&request.Body)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get reviewable request details")
	}

	state := history.ReviewableRequestStatePending
	if string(request.RejectReason) != "" {
		state = history.ReviewableRequestStateRejected
	}

	result := history.ReviewableRequest{
		TotalOrderID: db2.TotalOrderID{
			ID: int64(request.RequestId),
		},
		Requestor:    request.Requestor.Address(),
		Reviewer:     request.Reviewer.Address(),
		Reference:    reference,
		RejectReason: string(request.RejectReason),
		RequestType:  request.Body.Type,
		RequestState: state,
		Hash:         hex.EncodeToString(request.Hash[:]),
		Details:      details,
		CreatedAt:    unixToTime(int64(request.CreatedAt)),
		UpdatedAt:    ledgerCloseTime,
	}

	tasksExt := request.Tasks
	result.AllTasks = uint32(tasksExt.AllTasks)
	result.PendingTasks = uint32(tasksExt.PendingTasks)

	externalDetails := make([]regources.Details, 0, len(tasksExt.ExternalDetails))
	for _, item := range tasksExt.ExternalDetails {
		externalDetails = append(externalDetails, internal.MarshalCustomDetails(item))
	}

	// we use key "data" for compatibility with db2.Details
	// the value for the key "data" is a slice of map[string]interface{}
	result.ExternalDetails = map[string]interface{}{
		"data": externalDetails,
	}

	return &result, nil
}

func (c *reviewableRequestHandler) getAssetCreation(request *xdr.AssetCreationRequest) *history.AssetCreationRequest {
	return &history.AssetCreationRequest{
		Asset:                  string(request.Code),
		Policies:               int32(request.Policies),
		PreIssuedAssetSigner:   request.PreissuedAssetSigner.Address(),
		MaxIssuanceAmount:      amount.StringU(uint64(request.MaxIssuanceAmount)),
		InitialPreissuedAmount: amount.StringU(uint64(request.InitialPreissuedAmount)),
		Details:                internal.MarshalCustomDetails(request.CreatorDetails),
	}
}

func (c *reviewableRequestHandler) getAssetUpdate(request *xdr.AssetUpdateRequest) *history.AssetUpdateRequest {
	return &history.AssetUpdateRequest{
		Asset:    string(request.Code),
		Policies: int32(request.Policies),
		Details:  internal.MarshalCustomDetails(request.CreatorDetails),
	}
}

func (c *reviewableRequestHandler) getPreIssuanceRequest(request *xdr.PreIssuanceRequest) (*history.PreIssuanceRequest,
	error) {

	signature, err := xdr.MarshalBase64(request.Signature)
	if err != nil {
		return nil, errors.Wrap(err, "failed to marshal signature")
	}

	return &history.PreIssuanceRequest{
		Asset:     string(request.Asset),
		Amount:    amount.StringU(uint64(request.Amount)),
		Signature: signature,
		Reference: string(request.Reference),
	}, nil
}

func (c *reviewableRequestHandler) getIssuanceRequest(request *xdr.IssuanceRequest) *history.IssuanceRequest {
	return &history.IssuanceRequest{
		Asset:    string(request.Asset),
		Amount:   amount.StringU(uint64(request.Amount)),
		Receiver: request.Receiver.AsString(),
		Details:  internal.MarshalCustomDetails(request.ExternalDetails),
	}
}

func (c *reviewableRequestHandler) getWithdrawalRequest(request *xdr.WithdrawalRequest) *history.WithdrawalRequest {
	return &history.WithdrawalRequest{
		BalanceID:  request.Balance.AsString(),
		Amount:     amount.StringU(uint64(request.Amount)),
		FixedFee:   amount.StringU(uint64(request.Fee.Fixed)),
		PercentFee: amount.StringU(uint64(request.Fee.Percent)),
		Details:    internal.MarshalCustomDetails(request.CreatorDetails),
	}
}

func (c *reviewableRequestHandler) getAmlAlertRequest(request *xdr.AmlAlertRequest) *history.AmlAlertRequest {
	return &history.AmlAlertRequest{
		BalanceID: request.BalanceId.AsString(),
		Amount:    amount.StringU(uint64(request.Amount)),
		Reason:    string(request.CreatorDetails),
	}
}

func (c *reviewableRequestHandler) getSaleRequest(request *xdr.SaleCreationRequest) *history.SaleRequest {
	quoteAssets := make([]regources.AssetPrice, 0, len(request.QuoteAssets))
	for i := range request.QuoteAssets {
		quoteAssets = append(quoteAssets, regources.AssetPrice{
			Price: regources.Amount(int64(request.QuoteAssets[i].Price)),
			Asset: string(request.QuoteAssets[i].QuoteAsset),
		})
	}

	saleType := request.SaleTypeExt.SaleType
	baseAssetForHardCap := uint64(request.RequiredBaseAssetForHardCap)

	return &history.SaleRequest{
		BaseAsset:           string(request.BaseAsset),
		DefaultQuoteAsset:   string(request.DefaultQuoteAsset),
		StartTime:           unixToTime(int64(request.StartTime)),
		EndTime:             unixToTime(int64(request.EndTime)),
		SoftCap:             amount.StringU(uint64(request.SoftCap)),
		HardCap:             amount.StringU(uint64(request.HardCap)),
		Details:             internal.MarshalCustomDetails(request.CreatorDetails),
		QuoteAssets:         quoteAssets,
		SaleType:            saleType,
		BaseAssetForHardCap: amount.StringU(baseAssetForHardCap),
	}
}

func (c *reviewableRequestHandler) getLimitsUpdateRequest(
	request *xdr.LimitsUpdateRequest) *history.LimitsUpdateRequest {

	return &history.LimitsUpdateRequest{
		Details:      internal.MarshalCustomDetails(request.CreatorDetails),
		DocumentHash: hex.EncodeToString(request.DeprecatedDocumentHash[:]),
	}
}

func (c *reviewableRequestHandler) getChangeRoleRequest(request *xdr.ChangeRoleRequest) *history.ChangeRoleRequest {
	return &history.ChangeRoleRequest{
		DestinationAccount: request.DestinationAccount.Address(),
		AccountRoleToSet:   uint64(request.AccountRoleToSet),
		KYCData:            internal.MarshalCustomDetails(request.KycData),
		SequenceNumber:     uint32(request.SequenceNumber),
	}
}

func (c *reviewableRequestHandler) getUpdateSaleDetailsRequest(
	request *xdr.UpdateSaleDetailsRequest) *history.UpdateSaleDetailsRequest {
	return &history.UpdateSaleDetailsRequest{
		SaleID:     uint64(request.SaleId),
		NewDetails: internal.MarshalCustomDetails(request.CreatorDetails),
	}
}

func (c *reviewableRequestHandler) getAtomicSwapBidCreationRequest(request *xdr.ASwapBidCreationRequest,
) *history.AtomicSwapBidCreation {
	quoteAssets := make([]regources.AssetPrice, 0, len(request.QuoteAssets))
	for _, quoteAsset := range request.QuoteAssets {
		quoteAssets = append(quoteAssets, regources.AssetPrice{
			Asset: string(quoteAsset.QuoteAsset),
			Price: regources.Amount(quoteAsset.Price),
		})
	}

	return &history.AtomicSwapBidCreation{
		BaseBalance: request.BaseBalance.AsString(),
		BaseAmount:  uint64(request.Amount),
		Details:     internal.MarshalCustomDetails(request.CreatorDetails),
		QuoteAssets: quoteAssets,
	}
}

func (c *reviewableRequestHandler) getAtomicSwapRequest(request *xdr.ASwapRequest,
) *history.AtomicSwap {
	return &history.AtomicSwap{
		BidID:      uint64(request.BidId),
		BaseAmount: uint64(request.BaseAmount),
		QuoteAsset: string(request.QuoteAsset),
	}
}

func (c *reviewableRequestHandler) getReviewableRequestDetails(
	body *xdr.ReviewableRequestEntryBody) (history.ReviewableRequestDetails, error) {

	details := history.ReviewableRequestDetails{
		Type: body.Type,
	}

	var err error
	switch body.Type {
	case xdr.ReviewableRequestTypeCreateAsset:
		details.AssetCreation = c.getAssetCreation(body.AssetCreationRequest)
	case xdr.ReviewableRequestTypeUpdateAsset:
		details.AssetUpdate = c.getAssetUpdate(body.AssetUpdateRequest)
	case xdr.ReviewableRequestTypeCreateIssuance:
		details.IssuanceCreate = c.getIssuanceRequest(body.IssuanceRequest)
	case xdr.ReviewableRequestTypeCreatePreIssuance:
		details.PreIssuanceCreate, err = c.getPreIssuanceRequest(body.PreIssuanceRequest)
		if err != nil {
			return details, errors.Wrap(err, "failed to get pre issuance request")
		}
	case xdr.ReviewableRequestTypeCreateWithdraw:
		details.Withdraw = c.getWithdrawalRequest(body.WithdrawalRequest)
	case xdr.ReviewableRequestTypeCreateSale:
		details.Sale = c.getSaleRequest(body.SaleCreationRequest)
	case xdr.ReviewableRequestTypeUpdateLimits:
		details.LimitsUpdate = c.getLimitsUpdateRequest(body.LimitsUpdateRequest)
	case xdr.ReviewableRequestTypeCreateAmlAlert:
		details.AmlAlert = c.getAmlAlertRequest(body.AmlAlertRequest)
	case xdr.ReviewableRequestTypeChangeRole:
		details.ChangeRole = c.getChangeRoleRequest(body.ChangeRoleRequest)
	case xdr.ReviewableRequestTypeUpdateSaleDetails:
		details.UpdateSaleDetails = c.getUpdateSaleDetailsRequest(body.UpdateSaleDetailsRequest)
	case xdr.ReviewableRequestTypeCreateAtomicSwapBid:
		details.AtomicSwapBidCreation = c.getAtomicSwapBidCreationRequest(body.ASwapBidCreationRequest)
	case xdr.ReviewableRequestTypeCreateAtomicSwap:
		details.AtomicSwap = c.getAtomicSwapRequest(body.ASwapRequest)
	default:
		return details, errors.From(errors.New("unexpected reviewable request type"), map[string]interface{}{
			"request_type": body.Type.String(),
		})
	}

	return details, nil
}
