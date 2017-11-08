package resource

import (
	"gitlab.com/tokend/horizon/db2/core"
)

func (this *AccountThresholds) Populate(row core.Account) {
	this.LowThreshold = row.Thresholds[1]
	this.MedThreshold = row.Thresholds[2]
	this.HighThreshold = row.Thresholds[3]
}
