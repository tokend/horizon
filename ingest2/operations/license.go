package operations

import (
	"encoding/hex"

	"gitlab.com/tokend/go/xdr"
	"gitlab.com/tokend/horizon/db2/history2"
)

type licenseOpHandler struct {
	effectsProvider
}

// Details returns details about manage asset operation
func (h *licenseOpHandler) Details(op rawOperation, opRes xdr.OperationResultTr,
) (history2.OperationDetails, error) {
	licenseOp := op.Body.MustLicenseOp()

	opDetails := history2.OperationDetails{
		Type: xdr.OperationTypeLicense,
		License: &history2.LicenseDetails{
			AdminCount:      uint64(licenseOp.AdminCount),
			DueDate:         uint64(licenseOp.DueDate),
			LedgerHash:      hex.EncodeToString([]byte(licenseOp.LedgerHash[:])),
			PrevLicenseHash: hex.EncodeToString([]byte(licenseOp.PrevLicenseHash[:])),
		},
	}

	return opDetails, nil
}
