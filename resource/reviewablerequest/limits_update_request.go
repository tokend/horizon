package reviewablerequest

import (
	"gitlab.com/swarmfund/horizon/db2/history"
)

type LimitsUpdateRequest struct {
	DocumentHash string `json:"document_hash"`
	Details 	 string `json:"details"`
}

func (r *LimitsUpdateRequest) Populate(histRequest history.LimitsUpdateRequest) (error) {
	r.Details = histRequest.Details
	r.DocumentHash = histRequest.DocumentHash
	return nil
}