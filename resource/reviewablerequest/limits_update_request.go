package reviewablerequest

import (
	"gitlab.com/swarmfund/horizon/db2/history"
	"encoding/json"
	"gitlab.com/distributed_lab/logan/v3/errors"
)

type LimitsUpdateRequest struct {
	DocumentHash string `json:"document_hash"`
}

func (r *LimitsUpdateRequest) Populate(histRequest history.LimitsUpdateRequest) {
	r.DocumentHash = histRequest.DocumentHash;
}

func (r *LimitsUpdateRequest) PopulateFromRawJsonHistory(rawJson []byte) error {
	var histRequest history.LimitsUpdateRequest
	err := json.Unmarshal(rawJson, &histRequest)
	if err != nil {
		return errors.Wrap(err, "failed to unmarshal history.LimitsUpdateRequest")
	}

	r.Populate(histRequest)
	return nil
}