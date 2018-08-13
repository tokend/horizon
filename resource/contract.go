package resource

import (
	"strconv"

	"gitlab.com/swarmfund/horizon/db2/history"
	"gitlab.com/swarmfund/horizon/resource/base"
	"gitlab.com/tokend/go/xdr"
	"gitlab.com/tokend/regources"
)

func PopulateContract(row history.Contract) regources.Contract {
	var details []map[string]interface{}
	for _, item := range row.Details {
		details = append(details, item)
	}
	return regources.Contract{
		ID:            strconv.FormatInt(row.ID, 10),
		PT:            row.PagingToken(),
		Contractor:    row.Contractor,
		Customer:      row.Customer,
		Escrow:        row.Escrow,
		Disputer:      row.Disputer,
		StartTime:     row.StartTime,
		EndTime:       row.EndTime,
		Details:       details,
		DisputeReason: row.DisputeReason,
		State:         base.FlagFromXdrContractState(row.State, xdr.ContractStateAll),
	}
}
