package operaitons

import (
	"gitlab.com/tokend/go/amount"
	"gitlab.com/tokend/go/xdr"
	"gitlab.com/tokend/horizon/db2/history2"
)

type cancelAtomicSwapBidOpHandler struct {
	pubKeyProvider        publicKeyProvider
	ledgerChangesProvider ledgerChangesProvider
}

func (h *cancelAtomicSwapBidOpHandler) OperationDetails(op rawOperation,
	opRes xdr.OperationResultTr,
) (history2.OperationDetails, error) {
	return history2.OperationDetails{
		Type: xdr.OperationTypeCancelAswapBid,
		CancelAtomicSwapBid: &history2.CancelAtomicSwapBidDetails{
			BidID: int64(op.Body.MustCancelASwapBidOp().BidId),
		},
	}, nil
}

func (h *cancelAtomicSwapBidOpHandler) ParticipantsEffects(opBody xdr.OperationBody,
	opRes xdr.OperationResultTr, source history2.ParticipantEffect,
) ([]history2.ParticipantEffect, error) {
	successRes := opRes.MustCancelASwapBidResult().MustSuccess()

	if successRes.UnlockedAmount == 0 {
		return []history2.ParticipantEffect{source}, nil
	}

	balanceID := h.pubKeyProvider.GetBalanceID(successRes.BaseBalance)

	source.BalanceID = &balanceID
	source.AssetCode = &successRes.BaseAsset
	source.Effect = history2.Effect{
		Type: history2.EffectTypeUnlocked,
		AtomicSwap: &history2.AtomicSwapEffect{
			Amount: amount.StringU(uint64(successRes.UnlockedAmount)),
		},
	}

	return []history2.ParticipantEffect{source}, nil
}
