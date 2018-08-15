package reviewablerequest

import (
	"strconv"

	"gitlab.com/swarmfund/horizon/db2/history"
	"gitlab.com/tokend/go/amount"
	"gitlab.com/tokend/regources"
)

func PopulateInvoiceRequest(histRequest history.InvoiceRequest) (
	*regources.InvoiceRequest, error,
) {
	var contractID string
	if histRequest.ContractID != nil {
		contractID = strconv.FormatInt(*histRequest.ContractID, 10)
	}
	return &regources.InvoiceRequest{
		Amount:     amount.StringU(histRequest.Amount),
		Asset:      histRequest.Asset,
		ContractID: contractID,
		Details:    histRequest.Details,
	}, nil
}
