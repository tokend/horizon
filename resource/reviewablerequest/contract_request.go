package reviewablerequest

import (
	"time"

	"gitlab.com/swarmfund/horizon/db2/history"
)

type ContractRequest struct {
	Escrow    string                 `json:"escrow"`
	Details   map[string]interface{} `json:"details"`
	StartTime time.Time              `json:"start_time"`
	EndTime   time.Time              `json:"end_time"`
}

func (r *ContractRequest) Populate(histRequest history.ContractRequest) error {
	r.Escrow = histRequest.Escrow
	r.Details = histRequest.Details
	r.StartTime = histRequest.StartTime
	r.EndTime = histRequest.EndTime
	return nil
}
