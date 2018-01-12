package reviewablerequest

import (
	"gitlab.com/swarmfund/horizon/db2/history"
	"encoding/json"
	"gitlab.com/distributed_lab/logan/v3/errors"
)

type UpdateKYCRequest struct {
	KYCData string `json:"KYC_data"`
}

func (r *UpdateKYCRequest) Populate(histRequest history.UpdateKYCRequest) {
	r.KYCData = histRequest.KYCData
}

func (r *UpdateKYCRequest) PopulateFromRawJsonHistory(rawJson []byte) error {
	var histRequest history.UpdateKYCRequest
	err := json.Unmarshal(rawJson, &histRequest)
	if err != nil {
		return errors.Wrap(err, "failed to unmarshal history.UpdateKYCRequest")
	}

	r.Populate(histRequest)
	return nil
}
