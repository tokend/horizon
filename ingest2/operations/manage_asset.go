package operations

import (
	"gitlab.com/distributed_lab/logan/v3/errors"
	"gitlab.com/tokend/go/xdr"
	"gitlab.com/tokend/horizon/db2/history2"
	"gitlab.com/tokend/horizon/ingest2/internal"
	"gitlab.com/tokend/regources/v2"
)

type manageAssetOpHandler struct {
}

// Details returns details about manage asset operation
func (h *manageAssetOpHandler) Details(op rawOperation, opRes xdr.OperationResultTr,
) (regources.OperationDetails, error) {
	manageAssetOp := op.Body.MustManageAssetOp()

	opDetails := regources.OperationDetails{
		Type: xdr.OperationTypeManageAsset,
		ManageAsset: &regources.ManageAssetDetails{
			RequestID: int64(manageAssetOp.RequestId),
			Action:    manageAssetOp.Request.Action,
		},
	}

	if manageAssetOp.RequestId == 0 {
		opDetails.ManageAsset.RequestID = int64(opRes.MustManageAssetResult().MustSuccess().RequestId)
	}

	switch opDetails.ManageAsset.Action {
	case xdr.ManageAssetActionCreateAssetCreationRequest:
		creationDetails := manageAssetOp.Request.MustCreateAssetCreationRequest().CreateAsset

		policies := xdr.AssetPolicy(creationDetails.Policies)

		opDetails.ManageAsset.AssetCode = string(creationDetails.Code)
		opDetails.ManageAsset.Details = internal.MarshalCustomDetails(creationDetails.Details)
		opDetails.ManageAsset.Policies = &policies
		opDetails.ManageAsset.PreissuedSigner = creationDetails.PreissuedAssetSigner.Address()
		opDetails.ManageAsset.MaxIssuanceAmount = regources.Amount(creationDetails.MaxIssuanceAmount)
	case xdr.ManageAssetActionCreateAssetUpdateRequest:
		updateDetails := manageAssetOp.Request.MustCreateAssetUpdateRequest().UpdateAsset

		policies := xdr.AssetPolicy(updateDetails.Policies)

		opDetails.ManageAsset.AssetCode = string(updateDetails.Code)
		opDetails.ManageAsset.Details = internal.MarshalCustomDetails(updateDetails.Details)
		opDetails.ManageAsset.Policies = &policies
	case xdr.ManageAssetActionCancelAssetRequest:
	case xdr.ManageAssetActionChangePreissuedAssetSigner:
		data := manageAssetOp.Request.MustChangePreissuedSigner()

		opDetails.ManageAsset.AssetCode = string(data.Code)
		opDetails.ManageAsset.PreissuedSigner = data.AccountId.Address()
	case xdr.ManageAssetActionUpdateMaxIssuance:
		data := manageAssetOp.Request.MustUpdateMaxIssuance()

		opDetails.ManageAsset.AssetCode = string(data.AssetCode)
		opDetails.ManageAsset.MaxIssuanceAmount = regources.Amount(data.MaxIssuanceAmount)
	default:
		return regources.OperationDetails{}, errors.New("unexpected manage asset action")
	}

	return opDetails, nil
}

//ParticipantsEffects - returns source of the operation
func (h *manageAssetOpHandler) ParticipantsEffects(opBody xdr.OperationBody,
	opRes xdr.OperationResultTr, source history2.ParticipantEffect, _ []xdr.LedgerEntryChange,
) ([]history2.ParticipantEffect, error) {
	return []history2.ParticipantEffect{source}, nil
}
