package operations

import (
	"github.com/go-errors/errors"
	"gitlab.com/tokend/go/xdr"
	"gitlab.com/tokend/horizon/db2/history2"
	"gitlab.com/tokend/horizon/ingest2/internal"
	"gitlab.com/tokend/regources/v2"
)

type payoutHandler struct {
	balanceProvider balanceProvider
	pubKeyProvider  IDProvider
}

// CreatorDetails returns details about payout operation
func (h *payoutHandler) Details(op rawOperation, res xdr.OperationResultTr,
) (history2.OperationDetails, error) {
	payoutOp := op.Body.MustPayoutOp()
	payoutRes := res.MustPayoutResult().MustSuccess()

	return history2.OperationDetails{
		Type: xdr.OperationTypePayout,
		Payout: &history2.PayoutDetails{
			SourceAccountAddress: op.Source.Address(),
			SourceBalanceAddress: payoutOp.SourceBalanceId.AsString(),
			Asset:                string(payoutOp.Asset),
			MaxPayoutAmount:      regources.Amount(payoutOp.MaxPayoutAmount),
			MinAssetHolderAmount: regources.Amount(payoutOp.MinAssetHolderAmount),
			MinPayoutAmount:      regources.Amount(payoutOp.MinPayoutAmount),
			ExpectedFee:          internal.FeeFromXdr(payoutOp.Fee),
			ActualFee:            internal.FeeFromXdr(payoutRes.ActualFee),
			ActualPayoutAmount:   regources.Amount(payoutRes.ActualPayoutAmount),
		},
	}, nil
}

// ParticipantsEffects returns `charged` and `funded` effects
func (h *payoutHandler) ParticipantsEffects(opBody xdr.OperationBody,
	res xdr.OperationResultTr, source history2.ParticipantEffect, _ []xdr.LedgerEntryChange,
) ([]history2.ParticipantEffect, error) {
	payoutOp := opBody.MustPayoutOp()
	payoutRes := res.MustPayoutResult().MustSuccess()

	balance := h.balanceProvider.MustBalance(payoutOp.SourceBalanceId)

	if balance.AccountID != source.AccountID {
		return nil, errors.New("unexpected state, expected source owns source balance")
	}

	source.BalanceID = &balance.ID
	source.AssetCode = &balance.AssetCode
	source.Effect = &history2.Effect{
		Type: history2.EffectTypeCharged,
		Charged: &history2.BalanceChangeEffect{
			Amount: regources.Amount(payoutRes.ActualPayoutAmount),
			Fee:    internal.FeeFromXdr(payoutRes.ActualFee),
		},
	}

	participants := []history2.ParticipantEffect{source}

	responses := payoutRes.PayoutResponses
	for _, response := range responses {
		balanceID := h.pubKeyProvider.MustBalanceID(response.ReceiverBalanceId)

		participants = append(participants, history2.ParticipantEffect{
			AccountID: h.pubKeyProvider.MustAccountID(response.ReceiverId),
			BalanceID: &balanceID,
			AssetCode: &balance.AssetCode, // source balance has the same asset as receivers
			Effect: &history2.Effect{
				Type: history2.EffectTypeFunded,
				Funded: &history2.BalanceChangeEffect{
					Amount: regources.Amount(response.ReceivedAmount),
				},
			},
		})
	}

	return participants, nil
}
