package operations

import (
	"gitlab.com/distributed_lab/logan/v3/errors"
	"gitlab.com/tokend/go/xdr"
	"gitlab.com/tokend/horizon/db2/history2"
	"gitlab.com/tokend/regources/v2"
)

type manageLimitsOpHandler struct {
	pubKeyProvider IDProvider
}

// Details returns details about manage limits operation
func (h *manageLimitsOpHandler) Details(op rawOperation, opRes xdr.OperationResultTr,
) (regources.OperationDetails, error) {
	manageLimitsOp := op.Body.MustManageLimitsOp()

	opDetails := regources.OperationDetails{
		Type: xdr.OperationTypeManageLimits,
		ManageLimits: &regources.ManageLimitsDetails{
			Action: manageLimitsOp.Details.Action,
		},
	}

	switch opDetails.ManageLimits.Action {
	case xdr.ManageLimitsActionCreate:
		creationDetails := manageLimitsOp.Details.MustLimitsCreateDetails()

		opDetails.ManageLimits.Creation = &regources.ManageLimitsCreationDetails{
			AccountAddress:  creationDetails.AccountId.Address(), // Address() - smart, check for nil inside
			AccountType:     creationDetails.AccountType,
			StatsOpType:     creationDetails.StatsOpType,
			AssetCode:       string(creationDetails.AssetCode),
			IsConvertNeeded: creationDetails.IsConvertNeeded,
			DailyOut:        regources.Amount(creationDetails.DailyOut),
			WeeklyOut:       regources.Amount(creationDetails.WeeklyOut),
			MonthlyOut:      regources.Amount(creationDetails.MonthlyOut),
			AnnualOut:       regources.Amount(creationDetails.AnnualOut),
		}
	case xdr.ManageLimitsActionRemove:
		opDetails.ManageLimits.Removal = &regources.ManageLimitsRemovalDetails{
			LimitsID: int64(manageLimitsOp.Details.MustId()),
		}
	default:
		return regources.OperationDetails{}, errors.New("unexpected manage limits action")
	}

	return opDetails, nil
}

//ParticipantsEffects - returns source of the operation and account for which limits were managed if they are different
func (h *manageLimitsOpHandler) ParticipantsEffects(opBody xdr.OperationBody,
	_ xdr.OperationResultTr, source history2.ParticipantEffect, _ []xdr.LedgerEntryChange,
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

	accountID := h.pubKeyProvider.MustAccountID(*creationDetails.AccountId)

	if source.AccountID == accountID {
		return participants, nil
	}

	assetCode := string(creationDetails.AssetCode)
	participants = append(participants, history2.ParticipantEffect{
		AccountID: accountID,
		AssetCode: &assetCode,
	})

	return participants, nil
}
