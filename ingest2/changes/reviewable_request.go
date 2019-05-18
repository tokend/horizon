package changes

import (
	"encoding/hex"
	"time"

	"gitlab.com/distributed_lab/logan/v3"
	"gitlab.com/distributed_lab/logan/v3/errors"
	"gitlab.com/tokend/go/xdr"
	"gitlab.com/tokend/horizon/db2"
	history "gitlab.com/tokend/horizon/db2/history2"
	"gitlab.com/tokend/horizon/ingest2/internal"
	"gitlab.com/tokend/horizon/utf8"
	regources "gitlab.com/tokend/regources/generated"
)

var errUnknownRemoveReason = errors.New("request was removed due to unknown reason")
var removeOnKYCRecoveryInit = "New KYC recovery was initiated"

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
		return c.handleRemoveOnManageAsset(lc)
	case xdr.OperationTypeManageSale:
		return c.handleRemoveOnManageSale(lc)
	case xdr.OperationTypeCancelAswapBid:
		return c.cancel(lc)
	// auto review is handled by each operation separately
	case xdr.OperationTypeCreateIssuanceRequest:
		return c.handleRemoveOnCreationOp(lc,
			lc.OperationResult.MustCreateIssuanceRequestResult().MustSuccess().Fulfilled)
	case xdr.OperationTypeCreatePreissuanceRequest:
		return c.handleRemoveOnCreationOp(lc,
			lc.OperationResult.MustCreatePreIssuanceRequestResult().MustSuccess().Fulfilled)
	case xdr.OperationTypeCreateManageLimitsRequest:
		return c.handleRemoveOnCreationOp(lc,
			lc.OperationResult.MustCreateManageLimitsRequestResult().MustSuccess().Fulfilled)
	case xdr.OperationTypeCreateSaleRequest:
		return c.handleRemoveOnCreationOp(lc,
			lc.OperationResult.MustCreateSaleCreationRequestResult().MustSuccess().Fulfilled)
	case xdr.OperationTypeCreateAmlAlert:
		return c.handleRemoveOnCreationOp(lc,
			lc.OperationResult.MustCreateAmlAlertRequestResult().MustSuccess().Fulfilled)
	case xdr.OperationTypeCreateChangeRoleRequest:
		return c.handleRemoveOnCreationOp(lc,
			lc.OperationResult.MustCreateChangeRoleRequestResult().MustSuccess().Fulfilled)
	case xdr.OperationTypeCheckSaleState:
		// if check sale state was successful, all the requests created by it were fulfilled
		return c.handleRemoveOnCreationOp(lc, true)
	case xdr.OperationTypeCancelSaleRequest:
		return c.cancel(lc)
	case xdr.OperationTypeManageCreatePollRequest:
		return c.handleRemoveOnManageCreatePollRequest(lc)
	case xdr.OperationTypeInitiateKycRecovery:
		return c.handleInitiateKycRecovery(lc)
	default: // safeguard for future updates
		return errors.From(errUnknownRemoveReason, logan.F{
			"op_type": op.Type.String(),
		})
	}
}

func (c *reviewableRequestHandler) handleInitiateKycRecovery(lc ledgerChange) error {
	id := uint64(lc.LedgerChange.MustRemoved().MustReviewableRequest().RequestId)
	return c.storage.PermanentReject(id, removeOnKYCRecoveryInit)
}

func (c *reviewableRequestHandler) handleRemoveOnCreationOp(lc ledgerChange, fulfilled bool) error {
	id := uint64(lc.LedgerChange.MustRemoved().MustReviewableRequest().RequestId)
	if !fulfilled {
		return errors.From(errors.New("unexpected state: request has been removed on creation op, but fulfilled is false"), logan.F{
			"id": uint64(id),
		})
	}

	return c.storage.Approve(id)
}

func (c *reviewableRequestHandler) handleRemoveOnManageCreatePollRequest(lc ledgerChange) error {
	data := lc.Operation.Body.MustManageCreatePollRequestOp().Data
	switch data.Action {
	case xdr.ManageCreatePollRequestActionCreate:
		return c.handleRemoveOnCreationOp(lc,
			lc.OperationResult.MustManageCreatePollRequestResult().
				MustSuccess().Details.MustResponse().Fulfilled)
	case xdr.ManageCreatePollRequestActionCancel:
		return c.cancel(lc)
	default:
		return errors.Wrap(errUnknownRemoveReason,
			"Unexpected manage create poll request action", logan.F{
				"action": data.Action.String(),
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
		fulfilled := lc.OperationResult.MustManageAssetResult().MustSuccess().Fulfilled
		return c.handleRemoveOnCreationOp(lc, fulfilled)
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
		fulfilled := lc.OperationResult.MustManageSaleResult().MustSuccess().Fulfilled
		return c.handleRemoveOnCreationOp(lc, fulfilled)
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
	ledgerCloseTime time.Time,
) (*history.ReviewableRequest, error) {

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

func (c *reviewableRequestHandler) getAssetCreation(request *xdr.AssetCreationRequest) *history.CreateAssetRequest {
	return &history.CreateAssetRequest{
		Asset:                  string(request.Code),
		Type:                   uint64(request.Type),
		Policies:               int32(request.Policies),
		PreIssuedAssetSigner:   request.PreissuedAssetSigner.Address(),
		MaxIssuanceAmount:      regources.Amount(request.MaxIssuanceAmount),
		InitialPreissuedAmount: regources.Amount(request.InitialPreissuedAmount),
		CreatorDetails:         internal.MarshalCustomDetails(request.CreatorDetails),
	}
}

func (c *reviewableRequestHandler) getAssetUpdate(request *xdr.AssetUpdateRequest) *history.UpdateAssetRequest {
	return &history.UpdateAssetRequest{
		Asset:          string(request.Code),
		Policies:       int32(request.Policies),
		CreatorDetails: internal.MarshalCustomDetails(request.CreatorDetails),
	}
}

func (c *reviewableRequestHandler) getPreIssuanceRequest(request *xdr.PreIssuanceRequest,
) (*history.CreatePreIssuanceRequest, error) {

	signature, err := xdr.MarshalBase64(request.Signature)
	if err != nil {
		return nil, errors.Wrap(err, "failed to marshal signature")
	}

	return &history.CreatePreIssuanceRequest{
		Asset:     string(request.Asset),
		Amount:    regources.Amount(request.Amount),
		Signature: signature,
		Reference: string(request.Reference),
	}, nil
}

func (c *reviewableRequestHandler) getIssuanceRequest(request *xdr.IssuanceRequest) *history.CreateIssuanceRequest {
	return &history.CreateIssuanceRequest{
		Asset:          string(request.Asset),
		Amount:         regources.Amount(request.Amount),
		Receiver:       request.Receiver.AsString(),
		CreatorDetails: internal.MarshalCustomDetails(request.CreatorDetails),
	}
}

func (c *reviewableRequestHandler) getWithdrawalRequest(request *xdr.WithdrawalRequest) *history.CreateWithdrawalRequest {
	return &history.CreateWithdrawalRequest{
		BalanceID: request.Balance.AsString(),
		Amount:    regources.Amount(request.Amount),
		Fee: regources.Fee{
			Fixed:             regources.Amount(request.Fee.Fixed),
			CalculatedPercent: regources.Amount(request.Fee.Percent),
		},
		CreatorDetails: internal.MarshalCustomDetails(request.CreatorDetails),
	}
}

func (c *reviewableRequestHandler) getAmlAlertRequest(request *xdr.AmlAlertRequest) *history.CreateAmlAlertRequest {
	return &history.CreateAmlAlertRequest{
		BalanceID:      request.BalanceId.AsString(),
		Amount:         regources.Amount(request.Amount),
		CreatorDetails: internal.MarshalCustomDetails(request.CreatorDetails),
	}
}

func (c *reviewableRequestHandler) getSaleRequest(request *xdr.SaleCreationRequest) *history.CreateSaleRequest {
	quoteAssets := make([]regources.AssetPrice, 0, len(request.QuoteAssets))
	for i := range request.QuoteAssets {
		quoteAssets = append(quoteAssets, regources.AssetPrice{
			Price: regources.Amount(int64(request.QuoteAssets[i].Price)),
			Asset: string(request.QuoteAssets[i].QuoteAsset),
		})
	}

	saleType := request.SaleTypeExt.SaleType
	return &history.CreateSaleRequest{
		BaseAsset:           string(request.BaseAsset),
		DefaultQuoteAsset:   string(request.DefaultQuoteAsset),
		StartTime:           unixToTime(int64(request.StartTime)),
		EndTime:             unixToTime(int64(request.EndTime)),
		SoftCap:             regources.Amount(request.SoftCap),
		HardCap:             regources.Amount(request.HardCap),
		CreatorDetails:      internal.MarshalCustomDetails(request.CreatorDetails),
		QuoteAssets:         quoteAssets,
		SaleType:            saleType,
		BaseAssetForHardCap: regources.Amount(request.RequiredBaseAssetForHardCap),
	}
}

func (c *reviewableRequestHandler) getLimitsUpdateRequest(request *xdr.LimitsUpdateRequest,
) *history.UpdateLimitsRequest {

	return &history.UpdateLimitsRequest{
		CreatorDetails: internal.MarshalCustomDetails(request.CreatorDetails),
	}
}

func (c *reviewableRequestHandler) getChangeRoleRequest(request *xdr.ChangeRoleRequest,
) *history.ChangeRoleRequest {
	return &history.ChangeRoleRequest{
		DestinationAccount: request.DestinationAccount.Address(),
		AccountRoleToSet:   uint64(request.AccountRoleToSet),
		CreatorDetails:     internal.MarshalCustomDetails(request.CreatorDetails),
		SequenceNumber:     uint32(request.SequenceNumber),
	}
}

func (c *reviewableRequestHandler) getUpdateSaleDetailsRequest(
	request *xdr.UpdateSaleDetailsRequest) *history.UpdateSaleDetailsRequest {
	return &history.UpdateSaleDetailsRequest{
		SaleID:         uint64(request.SaleId),
		CreatorDetails: internal.MarshalCustomDetails(request.CreatorDetails),
	}
}

func (c *reviewableRequestHandler) getCreatePollRequest(
	request *xdr.CreatePollRequest) *history.CreatePollRequest {
	return &history.CreatePollRequest{
		NumberOfChoices:          uint32(request.NumberOfChoices),
		VoteConfirmationRequired: request.VoteConfirmationRequired,
		ResultProviderID:         request.ResultProviderId.Address(),
		PermissionType:           uint32(request.PermissionType),
		PollData:                 request.Data,
		StartTime:                unixToTime(int64(request.StartTime)),
		EndTime:                  unixToTime(int64(request.EndTime)),
		CreatorDetails:           internal.MarshalCustomDetails(request.CreatorDetails),
	}
}

func (c *reviewableRequestHandler) getAtomicSwapBidCreationRequest(request *xdr.ASwapBidCreationRequest,
) *history.CreateAtomicSwapBidRequest {
	quoteAssets := make([]regources.AssetPrice, 0, len(request.QuoteAssets))
	for _, quoteAsset := range request.QuoteAssets {
		quoteAssets = append(quoteAssets, regources.AssetPrice{
			Asset: string(quoteAsset.QuoteAsset),
			Price: regources.Amount(quoteAsset.Price),
		})
	}

	return &history.CreateAtomicSwapBidRequest{
		BaseBalance:    request.BaseBalance.AsString(),
		BaseAmount:     regources.Amount(request.Amount),
		CreatorDetails: internal.MarshalCustomDetails(request.CreatorDetails),
		QuoteAssets:    quoteAssets,
	}
}

func (c *reviewableRequestHandler) getAtomicSwapRequest(request *xdr.ASwapRequest,
) *history.CreateAtomicSwapRequest {
	return &history.CreateAtomicSwapRequest{
		BidID:      uint64(request.BidId),
		BaseAmount: regources.Amount(request.BaseAmount),
		QuoteAsset: string(request.QuoteAsset),
	}
}

func (c *reviewableRequestHandler) getKYCRecovery(request *xdr.KycRecoveryRequest,
) *history.KYCRecoveryRequest {
	signersData := make([]history.UpdateSignerDetails, 0, len(request.SignersData))
	for _, signer := range request.SignersData {
		signersData = append(signersData, history.UpdateSignerDetails{
			Details:  internal.MarshalCustomDetails(signer.Details),
			RoleID:   uint64(signer.RoleId),
			Identity: uint32(signer.Identity),
			Weight:   uint32(signer.Weight),
		})
	}

	return &history.KYCRecoveryRequest{
		TargetAccount:  request.TargetAccount.Address(),
		CreatorDetails: internal.MarshalCustomDetails(request.CreatorDetails),
		SequenceNumber: uint32(request.SequenceNumber),
		SignersData:    signersData,
	}
}

func (c *reviewableRequestHandler) getReviewableRequestDetails(
	body *xdr.ReviewableRequestEntryBody,
) (history.ReviewableRequestDetails, error) {

	details := history.ReviewableRequestDetails{
		Type: body.Type,
	}

	var err error
	switch body.Type {
	case xdr.ReviewableRequestTypeCreateAsset:
		details.CreateAsset = c.getAssetCreation(body.AssetCreationRequest)
	case xdr.ReviewableRequestTypeUpdateAsset:
		details.UpdateAsset = c.getAssetUpdate(body.AssetUpdateRequest)
	case xdr.ReviewableRequestTypeCreateIssuance:
		details.CreateIssuance = c.getIssuanceRequest(body.IssuanceRequest)
	case xdr.ReviewableRequestTypeCreatePreIssuance:
		details.CreatePreIssuance, err = c.getPreIssuanceRequest(body.PreIssuanceRequest)
		if err != nil {
			return details, errors.Wrap(err, "failed to get pre issuance request")
		}
	case xdr.ReviewableRequestTypeCreateWithdraw:
		details.CreateWithdraw = c.getWithdrawalRequest(body.WithdrawalRequest)
	case xdr.ReviewableRequestTypeCreateSale:
		details.CreateSale = c.getSaleRequest(body.SaleCreationRequest)
	case xdr.ReviewableRequestTypeUpdateLimits:
		details.UpdateLimits = c.getLimitsUpdateRequest(body.LimitsUpdateRequest)
	case xdr.ReviewableRequestTypeCreateAmlAlert:
		details.CreateAmlAlert = c.getAmlAlertRequest(body.AmlAlertRequest)
	case xdr.ReviewableRequestTypeChangeRole:
		details.ChangeRole = c.getChangeRoleRequest(body.ChangeRoleRequest)
	case xdr.ReviewableRequestTypeUpdateSaleDetails:
		details.UpdateSaleDetails = c.getUpdateSaleDetailsRequest(body.UpdateSaleDetailsRequest)
	case xdr.ReviewableRequestTypeCreateAtomicSwapBid:
		details.CreateAtomicSwapBid = c.getAtomicSwapBidCreationRequest(body.ASwapBidCreationRequest)
	case xdr.ReviewableRequestTypeCreateAtomicSwap:
		details.CreateAtomicSwap = c.getAtomicSwapRequest(body.ASwapRequest)
	case xdr.ReviewableRequestTypeCreatePoll:
		details.CreatePoll = c.getCreatePollRequest(body.CreatePollRequest)
	case xdr.ReviewableRequestTypeKycRecovery:
		details.KYCRecovery = c.getKYCRecovery(body.KycRecoveryRequest)
	default:
		return details, errors.From(errors.New("unexpected reviewable request type"), map[string]interface{}{
			"request_type": body.Type.String(),
		})
	}

	return details, nil
}
