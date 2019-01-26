package resources

import (
	"gitlab.com/distributed_lab/logan/v3"
	"gitlab.com/distributed_lab/logan/v3/errors"
	"gitlab.com/tokend/go/xdr"
	"gitlab.com/tokend/horizon/db2/history2"
	"gitlab.com/tokend/regources/v2"
)

func newManageLimits(id int64, details history2.ManageLimitsDetails) *regources.ManageLimits {
	result := regources.ManageLimits{
		Key: regources.NewKeyInt64(id, regources.TypeManageLimits),
		Attributes: regources.ManageLimitsAttributes{
			Action: details.Action,
		},
	}

	switch details.Action {
	case xdr.ManageLimitsActionCreate:
		result.Attributes.Create = newManageLimitsRemove(*details.Creation)
	case xdr.ManageLimitsActionRemove:
		result.Attributes.Remove = &regources.ManageLimitsRemoval{
			LimitsID: details.Removal.LimitsID,
		}
	default:
		panic(errors.From(errors.New("unexpected manage limits action"), logan.F{
			"action": details.Action,
		}))
	}

	return &result
}

func newManageLimitsRemove(details history2.ManageLimitsCreationDetails) *regources.ManageLimitsCreation {
	return &regources.ManageLimitsCreation{
		AccountType:     details.AccountType,
		AccountAddress:  details.AccountAddress,
		StatsOpType:     details.StatsOpType,
		AssetCode:       details.AssetCode,
		IsConvertNeeded: details.IsConvertNeeded,
		DailyOut:        details.DailyOut,
		WeeklyOut:       details.WeeklyOut,
		AnnualOut:       details.AnnualOut,
		MonthlyOut:      details.MonthlyOut,
	}
}
