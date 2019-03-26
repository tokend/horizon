package core2

import (
	"time"

	"gitlab.com/tokend/go/xdr"
	"gitlab.com/tokend/regources/v2"
)

type Poll struct {
	ID                       int64             `db:"id"`
	PermissionType           uint64            `db:"permission_type"`
	NumberOfChoices          uint64            `db:"number_of_choices"`
	Type                     int32             `db:"type"`
	Data                     xdr.PollData      `db:"data"`
	StartTime                time.Time         `db:"start_time"`
	EndTime                  time.Time         `db:"end_time"`
	OwnerID                  string            `db:"owner_id"`
	ResultProviderID         string            `db:"result_provider_id"`
	VoteConfirmationRequired bool              `db:"vote_confirmation_required"`
	Details                  regources.Details `db:"details"`
	LastModified             int32             `db:"lastmodified"`
}
