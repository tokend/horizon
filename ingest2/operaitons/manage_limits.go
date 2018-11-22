package operaitons

import (
	"gitlab.com/distributed_lab/logan/v3/errors"
	"gitlab.com/tokend/go/amount"
	"gitlab.com/tokend/go/xdr"
	"gitlab.com/tokend/horizon/db2/history2"
)

type manageLimitsOpHandler struct {
	pubKeyProvider publicKeyProvider
}

func (h *manageLimitsOpHandler) OperationDetails(opBody xdr.OperationBody,
	opRes xdr.OperationResultTr,
) (history2.OperationDetails, error) {
	manageLimitsOp := opBody.MustManageLimitsOp()

	opDetails := history2.OperationDetails{
		Type: xdr.OperationTypeManageLimits,
		ManageLimits: &history2.ManageLimitsDetails{
			Action: manageLimitsOp.Details.Action,
		},
	}

	switch opDetails.ManageLimits.Action {
	case xdr.ManageLimitsActionCreate:
		creationDetails := manageLimitsOp.Details.MustLimitsCreateDetails()

		var accountID *int64
		if creationDetails.AccountId != nil {
			accountIDInt := h.pubKeyProvider.GetAccountID(*creationDetails.AccountId)
			accountID = &accountIDInt
		}

		opDetails.ManageLimits.Creation = &history2.ManageLimitsCreationDetails{
			AccountID:       accountID,
			AccountType:     creationDetails.AccountType,
			StatsOpType:     creationDetails.StatsOpType,
			AssetCode:       creationDetails.AssetCode,
			IsConvertNeeded: creationDetails.IsConvertNeeded,
			DailyOut:        amount.StringU(uint64(creationDetails.DailyOut)),
			WeeklyOut:       amount.StringU(uint64(creationDetails.WeeklyOut)),
			MonthlyOut:      amount.StringU(uint64(creationDetails.MonthlyOut)),
			AnnualOut:       amount.StringU(uint64(creationDetails.AnnualOut)),
		}
	case xdr.ManageLimitsActionRemove:
		opDetails.ManageLimits.Removal = &history2.ManageLimitsRemovalDetails{
			ID: int64(manageLimitsOp.Details.MustId()),
		}
	default:
		return history2.OperationDetails{}, errors.New("unexpected manage limits action")
	}

	return opDetails, nil
}

func (h *manageLimitsOpHandler) ParticipantsEffects(opBody xdr.OperationBody,
	opRes xdr.OperationResultTr, source history2.ParticipantEffect,
) ([]history2.ParticipantEffect, error) {
	participants := []history2.ParticipantEffect{source}

	manageLimitsOp := opBody.MustManageLimitsOp()

	if manageLimitsOp.Details.Action != xdr.ManageLimitsActionCreate {
		return participants, nil
	}

	creationDetails := manageLimitsOp.Details.MustLimitsCreateDetails()

	if creationDetails.AccountId == nil {
		return participants, nil
	}

	accountID := h.pubKeyProvider.GetAccountID(*creationDetails.AccountId)

	if source.AccountID == accountID {
		return participants, nil
	}

	participants = append(participants, history2.ParticipantEffect{
		AccountID: accountID,
		AssetCode: &creationDetails.AssetCode,
	})

	return participants, nil
}
