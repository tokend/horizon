package reviewablerequest

import (
	"strconv"

	"gitlab.com/swarmfund/horizon/db2/history"
	"gitlab.com/tokend/go/amount"
	"gitlab.com/tokend/regources/reviewablerequest2"
)

func PopulateInvoiceRequest(histRequest history.InvoiceRequest) (
	*reviewablerequest2.InvoiceRequest, error,
) {
	var contractID string
	if histRequest.ContractID != nil {
		contractID = strconv.FormatInt(*histRequest.ContractID, 10)
	}
	return &reviewablerequest2.InvoiceRequest{
		Amount:     amount.StringU(histRequest.Amount),
		Asset:      histRequest.Asset,
		ContractID: contractID,
		Details:    histRequest.Details,
	}, nil
}
