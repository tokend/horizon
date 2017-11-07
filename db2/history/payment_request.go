package history

import (
	"encoding/json"
	"time"

	"bullioncoin.githost.io/development/horizon/db2"
	"github.com/go-errors/errors"
	"github.com/guregu/null"
)

type PaymentRequest struct {
	db2.TotalOrderID
	PaymentID     uint64      `db:"payment_id"`
	PaymentState  *uint32     `db:"state"`
	Exchange      string      `db:"exchange"`
	Accepted      *bool       `db:"accepted"`
	DetailsString null.String `db:"details"`
	CreatedAt     time.Time   `db:"created_at"`
	UpdatedAt     time.Time   `db:"updated_at"`
	RequestType   int32       `db:"request_type"`
}

func (r *PaymentRequest) UnmarshalDetails(dest interface{}) error {
	if !r.DetailsString.Valid {
		return nil
	}

	err := json.Unmarshal([]byte(r.DetailsString.String), &dest)
	if err != nil {
		err = errors.Wrap(err, 1)
	}

	return err
}
