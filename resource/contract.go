package resource

import (
	"gitlab.com/swarmfund/horizon/db2/history"
	"gitlab.com/tokend/regources"
		"strconv"
)

func PopulateContract(row history.Contract) regources.Contract {
	return regources.Contract{
		ID: strconv.FormatInt(row.ID, 10),
		PT: row.PagingToken(),
		Contractor: row.Contractor,
		Customer: row.Customer,
		Escrow: row.Escrow,
		Disputer: row.Disputer,
		StartTime: row.StartTime,
		EndTime: row.EndTime,
		Details: row.Details,
		DisputeReason: row.DisputeReason,
		State: strconv.FormatInt(int64(row.State), 10),
	}
}
