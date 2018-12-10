package operaitons

import (
	"gitlab.com/distributed_lab/logan/v3/errors"
	"gitlab.com/tokend/go/amount"
	"gitlab.com/tokend/go/xdr"
	"gitlab.com/tokend/horizon/db2/history2"
)

type manageAssetOpHandler struct {
}

// OperationDetails returns details about manage asset operation
func (h *manageAssetOpHandler) OperationDetails(op rawOperation, opRes xdr.OperationResultTr,
) (history2.OperationDetails, error) {
	manageAssetOp := op.Body.MustManageAssetOp()

	opDetails := history2.OperationDetails{
		Type: xdr.OperationTypeManageAsset,
		ManageAsset: &history2.ManageAssetDetails{
			RequestID: int64(manageAssetOp.RequestId),
			Action:    manageAssetOp.Request.Action,
		},
	}

	if manageAssetOp.RequestId == 0 {
		opDetails.ManageAsset.RequestID = int64(opRes.MustManageAssetResult().MustSuccess().RequestId)
	}

	switch opDetails.ManageAsset.Action {
	case xdr.ManageAssetActionCreateAssetCreationRequest:
		creationDetails := manageAssetOp.Request.MustCreateAsset()

		policies := int32(creationDetails.Policies)

		opDetails.ManageAsset.AssetCode = creationDetails.Code
		opDetails.ManageAsset.Details = customDetailsUnmarshal([]byte(creationDetails.Details))
		opDetails.ManageAsset.Policies = &policies
		opDetails.ManageAsset.PreissuedSigner = creationDetails.PreissuedAssetSigner.Address()
		opDetails.ManageAsset.MaxIssuanceAmount = amount.StringU(uint64(creationDetails.MaxIssuanceAmount))
	case xdr.ManageAssetActionCreateAssetUpdateRequest:
		updateDetails := manageAssetOp.Request.MustUpdateAsset()

		policies := int32(updateDetails.Policies)

		opDetails.ManageAsset.AssetCode = updateDetails.Code
		opDetails.ManageAsset.Details = customDetailsUnmarshal([]byte(updateDetails.Details))
		opDetails.ManageAsset.Policies = &policies
	case xdr.ManageAssetActionCancelAssetRequest:
	case xdr.ManageAssetActionChangePreissuedAssetSigner:
		data := manageAssetOp.Request.MustChangePreissuedSigner()

		opDetails.ManageAsset.AssetCode = data.Code
		opDetails.ManageAsset.PreissuedSigner = data.AccountId.Address()
	case xdr.ManageAssetActionUpdateMaxIssuance:
		data := manageAssetOp.Request.MustUpdateMaxIssuance()

		opDetails.ManageAsset.AssetCode = data.AssetCode
		opDetails.ManageAsset.MaxIssuanceAmount = amount.StringU(uint64(data.MaxIssuanceAmount))
	default:
		return history2.OperationDetails{}, errors.New("unexpected manage asset action")
	}

	return opDetails, nil
}

func (h *manageAssetOpHandler) ParticipantsEffects(opBody xdr.OperationBody,
	opRes xdr.OperationResultTr, source history2.ParticipantEffect,
) ([]history2.ParticipantEffect, error) {
	return []history2.ParticipantEffect{source}, nil
}
