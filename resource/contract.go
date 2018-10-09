package resource

import (
	"strconv"

	"gitlab.com/tokend/horizon/db2/history"
	"gitlab.com/tokend/horizon/resource/base"
	"gitlab.com/tokend/go/xdr"
	"gitlab.com/tokend/regources"
)

func PopulateContract(row history.Contract) regources.Contract {
	return regources.Contract{
		ID:              strconv.FormatInt(row.ID, 10),
		PT:              row.PagingToken(),
		Contractor:      row.Contractor,
		Customer:        row.Customer,
		Escrow:          row.Escrow,
		StartTime:       row.StartTime,
		EndTime:         row.EndTime,
		InitialDetails:  row.InitialDetails,
		CustomerDetails: row.CustomerDetails,
		State:           base.FlagFromXdrContractState(row.State, xdr.ContractStateAll),
	}
}

func PopulateContractDetails(row history.ContractDetails) regources.DetailsWithPayload {
	return regources.DetailsWithPayload{
		Details:   row.Details,
		Author:    row.Author,
		CreatedAt: row.CreatedAt,
	}
}

func PopulateContractDispute(row history.ContractDispute) *regources.DetailsWithPayload {
	return &regources.DetailsWithPayload{
		Details:   row.Reason,
		Author:    row.Author,
		CreatedAt: row.CreatedAt,
	}
}
