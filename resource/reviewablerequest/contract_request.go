package reviewablerequest

import (
	"gitlab.com/swarmfund/horizon/db2/history"
	"gitlab.com/tokend/regources"
)

func PopulateContractRequest(histRequest history.ContractRequest) (
	*regources.ContractRequest, error,
) {
	return &regources.ContractRequest{
		Escrow:    histRequest.Escrow,
		Details:   histRequest.Details,
		StartTime: histRequest.StartTime,
		EndTime:   histRequest.EndTime,
	}, nil
}
