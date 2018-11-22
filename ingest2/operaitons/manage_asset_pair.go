package operaitons

import (
	"gitlab.com/tokend/go/amount"
	"gitlab.com/tokend/go/xdr"
	"gitlab.com/tokend/horizon/db2/history2"
)

type manageAssetPairOpHadler struct {
}

func (h *manageAssetPairOpHadler) OperationDetails(opBody xdr.OperationBody,
	opRes xdr.OperationResultTr,
) (history2.OperationDetails, error) {
	manageAssetPairOp := opBody.MustManageAssetPairOp()

	return history2.OperationDetails{
		Type: xdr.OperationTypeManageAssetPair,
		ManageAssetPair: &history2.ManageAssetPairDetails{
			BaseAsset:               manageAssetPairOp.Base,
			QuoteAsset:              manageAssetPairOp.Quote,
			PhysicalPrice:           amount.String(int64(manageAssetPairOp.PhysicalPrice)),
			PhysicalPriceCorrection: amount.String(int64(manageAssetPairOp.PhysicalPriceCorrection)),
			MaxPriceStep:            amount.String(int64(manageAssetPairOp.MaxPriceStep)),
			PoliciesI:               int32(manageAssetPairOp.Policies),
		},
	}, nil
}

func (h *manageAssetPairOpHadler) ParticipantsEffects(opBody xdr.OperationBody,
	opRes xdr.OperationResultTr, source history2.ParticipantEffect,
) ([]history2.ParticipantEffect, error) {
	return []history2.ParticipantEffect{source}, nil
}
