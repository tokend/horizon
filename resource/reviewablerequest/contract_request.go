package reviewablerequest

import (
	"gitlab.com/swarmfund/horizon/db2/history"
	"gitlab.com/tokend/regources/reviewablerequest2"
)

func PopulateContractRequest(histRequest history.ContractRequest) (
	*reviewablerequest2.ContractRequest, error,
) {
	return &reviewablerequest2.ContractRequest{
		Escrow:    histRequest.Escrow,
		Details:   histRequest.Details,
		StartTime: histRequest.StartTime,
		EndTime:   histRequest.EndTime,
	}, nil
}
