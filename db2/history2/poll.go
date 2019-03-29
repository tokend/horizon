package history2

import (
	"time"

	"gitlab.com/tokend/regources/v2"
)

// Poll - represents instance of voting campaign
type Poll struct {
	ID                       int64               `db:"id"`
	PermissionType           uint32              `db:"permission_type"`
	NumberOfChoices          uint32              `db:"number_of_choices"`
	Data                     regources.PollData  `db:"data"`
	StartTime                time.Time           `db:"start_time"`
	EndTime                  time.Time           `db:"end_time"`
	OwnerID                  string              `db:"owner_id"`
	ResultProviderID         string              `db:"result_provider_id"`
	VoteConfirmationRequired bool                `db:"vote_confirmation_required"`
	Details                  regources.Details   `db:"details"`
	State                    regources.PollState `db:"state"`
}
