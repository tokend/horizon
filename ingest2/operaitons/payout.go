package operaitons

import (
	"github.com/go-errors/errors"
	"gitlab.com/tokend/go/amount"
	"gitlab.com/tokend/go/xdr"
	"gitlab.com/tokend/horizon/db2/history2"
)

type payoutHandler struct {
	balanceProvider balanceProvider
	pubKeyProvider  publicKeyProvider
}

func (h *payoutHandler) OperationDetails(op rawOperation, res xdr.OperationResultTr,
) (history2.OperationDetails, error) {
	payoutOp := op.Body.MustPayoutOp()
	payoutRes := res.MustPayoutResult().MustSuccess()

	return history2.OperationDetails{
		Type: xdr.OperationTypePayout,
		Payout: &history2.PayoutDetails{
			SourceAccountID:      op.Source.Address(),
			SourceBalanceID:      payoutOp.SourceBalanceId.AsString(),
			Asset:                payoutOp.Asset,
			MaxPayoutAmount:      amount.StringU(uint64(payoutOp.MaxPayoutAmount)),
			MinAssetHolderAmount: amount.StringU(uint64(payoutOp.MinAssetHolderAmount)),
			MinPayoutAmount:      amount.StringU(uint64(payoutOp.MinPayoutAmount)),
			ExpectedFixedFee:     amount.StringU(uint64(payoutOp.Fee.Fixed)),
			ExpectedPercentFee:   amount.StringU(uint64(payoutOp.Fee.Percent)),
			ActualFixedFee:       amount.StringU(uint64(payoutRes.ActualFee.Fixed)),
			ActualPercentFee:     amount.StringU(uint64(payoutRes.ActualFee.Percent)),
			ActualPayoutAmount:   amount.StringU(uint64(payoutRes.ActualPayoutAmount)),
		},
	}, nil
}

func (h *payoutHandler) ParticipantsEffects(opBody xdr.OperationBody,
	res xdr.OperationResultTr, source history2.ParticipantEffect,
) ([]history2.ParticipantEffect, error) {
	payoutOp := opBody.MustPayoutOp()
	payoutRes := res.MustPayoutResult().MustSuccess()

	balance := h.balanceProvider.GetBalanceByID(payoutOp.SourceBalanceId)

	if balance.AccountID != source.AccountID {
		return nil, errors.New("unexpected state, expected source owns source balance")
	}

	source.BalanceID = &balance.BalanceID
	source.AssetCode = &balance.AssetCode
	source.Effect = history2.Effect{
		Type: history2.EffectTypeCharged,
		Charged: &history2.ChargedEffect{
			Amount: amount.StringU(uint64(payoutRes.ActualPayoutAmount)),
			FeePaid: history2.FeePaid{
				Fixed:             amount.StringU(uint64(payoutRes.ActualFee.Fixed)),
				CalculatedPercent: amount.StringU(uint64(payoutRes.ActualFee.Percent)),
			},
		},
	}

	participants := []history2.ParticipantEffect{source}

	responses := payoutRes.PayoutResponses
	for _, response := range responses {
		balanceID := h.pubKeyProvider.GetBalanceID(response.ReceiverBalanceId)

		participants = append(participants, history2.ParticipantEffect{
			AccountID: h.pubKeyProvider.GetAccountID(response.ReceiverId),
			BalanceID: &balanceID,
			AssetCode: &balance.AssetCode, // source balance has the same asset as receivers
			Effect: history2.Effect{
				Type: history2.EffectTypeFunded,
				Funded: &history2.FundedEffect{
					Amount: amount.StringU(uint64(response.ReceivedAmount)),
				},
			},
		})
	}

	return participants, nil
}
