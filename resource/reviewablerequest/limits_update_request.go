package reviewablerequest

import (
	"gitlab.com/swarmfund/horizon/db2/history"
)

type LimitsUpdateRequest struct {
	DocumentHash string `json:"document_hash"`
}

func (r *LimitsUpdateRequest) Populate(histRequest history.LimitsUpdateRequest) (error) {
	r.DocumentHash = histRequest.DocumentHash;
	return nil
}